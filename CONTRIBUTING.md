# Requirements

Use [gvm](https://github.com/moovweb/gvm) to match the project's golang version,
or use a version of golang from your system that matches the one indicated in `go.mod`.

Install [taskfile](https://taskfile.dev/installation/) if you wish to run local tasks
from our `Taskfile.yaml`.

# Building

Verify the go binary builds:
```bash
go build .
```

Package into a docker image, and upload to a local repo:
```bash
# Set your local repo in the environment (not commited)
echo LOCAL_REGISRY=nexus.example.com:8080/myrepo > .env.local

# Run the task to build and upload the docker image
task develop
```

## Building behind corporate proxy

Docker/Podman support passing the proxy environment variable to the image
being built.
```bash
https_proxy=http://proxy.example.com:8080
http_proxy=http://proxy.example.com:8080
no_proxy=example.com
```

## Building behind TLS termination proxy

To use a custom certificate authority during the docker build,
simply drop your custom CA in pem format in the `.ca-bundle/`
directory.

```bash
# On Ubuntu
cp /usr/local/share/ca-certificates/* .ca-bundle/
```
```bash
# On RHEL
cp /etc/pki/ca-trust/source/anchors/* .ca-bundle/
```

It will be added to the docker intermediate build image to fetch
dependencies, but not to the final image.
