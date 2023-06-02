package pkg

import (
  "fmt"
  "context"
  "net/http"
  "os"
  "crypto/tls"
  "crypto/x509"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
)

type Handler struct {}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  ctx := context.Background()
  reg := prometheus.NewRegistry()
  var(
    oidcUp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "oidc_up",
        Help: "Active check that verify the user could connect to the OIDC module",
    }, []string{"module", "url", "username", "status", "reason"})
  )

  reg.MustRegister(oidcUp)

  handler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

  // Run
  mod := r.URL.Query().Get("module")
  if mod == "" {
    w.WriteHeader(http.StatusBadRequest)
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte("Configuration error: param 'module' should be specified"))
    return
  }

  cfg, ok := config.Modules[mod]
  if !ok {
    w.WriteHeader(http.StatusNotFound)
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte(fmt.Sprintf("Configuration error: module=%s not found in configuration", mod)))
    return
  }

  labels := prometheus.Labels{
    "module": mod,
    "url": cfg.Url,
    "username": cfg.Username,
  }

  // CA certificate custom http client
  if config.CaCert != "" {
    cas, err := os.ReadFile(config.CaCert)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      w.Header().Set("Content-Type", "text/plan")
      w.Write([]byte(fmt.Sprintf("Error reading certificate: %s", err)))
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(cas)
    client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: caCertPool}}}
    ctx = oidc.ClientContext(ctx, client)
  }

	provider, err := oidc.NewProvider(ctx, cfg.Url)
	if err != nil {
    labels["status"] = "oidc-failed"
    labels["reason"] = err.Error()
    oidcUp.With(labels).Set(0)
    handler.ServeHTTP(w, r)
    return
	}

  auth := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{"profile", "roles"},
	}

  token, err := auth.PasswordCredentialsToken(ctx, cfg.Username, cfg.Password)
  if err != nil {
    labels["status"] = "oauth2-failed"
    labels["reason"] = err.Error()
    oidcUp.With(labels).Set(0)
    handler.ServeHTTP(w, r)
    return
  }

  if !token.Valid() {
    labels["status"] = "invalid-token"
    labels["reason"] = err.Error()
    oidcUp.With(labels).Set(0)
    handler.ServeHTTP(w, r)
    return
  }

  // OK case
  labels["status"] = "ok"
  labels["reason"] = ""
  oidcUp.With(labels).Set(1)
  handler.ServeHTTP(w, r)
}
