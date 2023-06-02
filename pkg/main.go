package pkg

import (
  "net/http"

  log "github.com/sirupsen/logrus"
)

func Run() {
  initConfig()

  handler := &Handler{}
  http.Handle("/metrics", handler)
  http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Ready"))
  })
  http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Healthy"))
  })

  if err := http.ListenAndServe(":9123", nil); err != nil {
    log.Fatal(err)
  }
}
