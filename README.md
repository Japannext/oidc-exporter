# OpenIDConnect Exporter

Test an OIDC Idp backend (keycloak, dex, etc) ability to authenticate a static user.
It is very similar to the blackbox-exporter in behavior (it works with a probe).

## Motivation

Authentication backend such as Keycloak are hard to passively monitor based on their
error messages for the following reasons:
* It is often unclear if a user entered its password badly or if there is a misconfiguration
  in the OpenID Connect backend.
* When a downtime happens, end-users will be the first to report, sometimes a long time after the
  downtime actually started.

As such, an active monitoring for OpenID Connect backends is often required in order to measure
the unavailability of such service more precisely. `oidc-exporter` is this active check.

## Installation

For all installation methods, you need to create a static user/service account in the Idp that will be tested.


### On Kubernetes

```bash
helm install oidc-exporter oci://ghcr.io/japannext/helm-charts --version 1.0.0 --values values.yaml
```

`values.yaml` example:
```yaml
---
modules:
  myorg:
    url: https://keycloak.example.com/realms/myorg
    clientId: oidc-exporter
    username: oidc-exporter
    passwordSecretName: oidc-exporter-password
    passwordSecretKey: password

# Only if you have Japannext's keycloak-operator installed
keycloak:
  enabled: true
  realm: myorg

# Example of CA injection from a trust-manager generated configmap
cacertConfigMap: ca-bundle
```

The helm recipe will create:
* A Deployment of oidc-exporter
* A prometheus-operator's Probe
* A Japannext keycloak-operator's KeycloakClient

### Using docker

```bash
docker run ghcr.io/japannext/oidc-exporter:1.0.0 \
  -v /etc/oidc-exporter.yaml:/etc/oidc-exporter.yaml
```

`oidc-exporter.yaml` example:
```yaml
---
modules:
  myorg:
    url: https://keycloak.example.com/realms/myorg
    clientID: oidc-exporter
    clientSecret: xxx-xxx-xxx-xxx
    username: oidc-exporter
    password: xxx-xxx-xxx-xxx

# Optional
cacert: /path/to/ca.pem
```

For testing:
```
curl localhost:9123/metrics?module=myorg
```

## Alerting

Recommended Prometheus alert:

```yaml
- alert: KeycloakAuthDown  # / DexAuthDown / AuthentikAuthDown
  expr: oidc_up == 0
  for: 5m
  labels:
    module: "{{ $labels.module }}"
    url: "{{ $labels.url }}"
    severity: critical
  annotations:
    summary: OpenIDConnect authentication failed (module={{ $labels.module }})
    description: |
      Test account {{ $labels.username }} failed to connect to OIDC backend
      at {{ $labels.url }} for more than 5 minutes.
      Status = {{ $labels.status }}
      Reason = {{ $labels.reason }}
```
