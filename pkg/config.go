package pkg

import (
  "fmt"
  "os"

  log "github.com/sirupsen/logrus"
  yaml "gopkg.in/yaml.v3"
)

type Config struct {
  Modules map[string]*ModuleConfig `yaml:"modules"`
  CaCert string `yaml:"cacert"`
}

type ModuleConfig struct {
  Url string `yaml:"url"`
  ClientID string `yaml:"client_id"`
  ClientSecret string `yaml:"client_secret"`
  Username string `yaml:"username"`
  Password string `yaml:"password"`
}

var config *Config

func initConfig() {
  cf := os.Getenv("OIDC_EXPORTER_CONFIG_FILE")
  if cf == "" {
    cf = "/etc/oidc-exporter.yaml"
  }

  doc, err := os.ReadFile(cf)
  if err != nil {
    log.Fatal(err)
  }

  err = yaml.Unmarshal([]byte(doc), &config)
  if err != nil {
    log.Fatal(err)
  }

  // Detect environment variables
  for name, mod := range config.Modules {
    csKey := fmt.Sprintf("OIDC_EXPORTER_%s_CLIENT_SECRET", name)
    clientSecret := os.Getenv(csKey)
    if mod.ClientSecret == "" && clientSecret != "" {
      log.Infof("[%s] Set 'client_secret' through environment variable (%s)", name, csKey)
      mod.ClientSecret = clientSecret
    }
    if mod.ClientSecret == "" {
      log.Fatalf("[%s] No 'client_secret' set. Set it in the config or through %s", name, csKey)
    }

    pKey := fmt.Sprintf("OIDC_EXPORTER_%s_PASSWORD", name)
    password := os.Getenv(pKey)
    if mod.Password == "" && password != "" {
      log.Infof("[%s] Set 'password' through environment variable (%s)", name, pKey)
      mod.Password = password
    }
    if mod.Password == "" {
      log.Fatalf("[%s] No 'password' set. Set it in the config or through %s", name, pKey)
    }
  }
}
