FROM golang:1.20 as build

WORKDIR /code

# Local CA use
COPY .ca-bundle /usr/local/share/ca-certificates/
RUN chmod -R 644 /usr/local/share/ca-certificates/ && update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download

# Build
COPY pkg pkg
COPY main.go main.go
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o /build/oidc-exporter .

FROM scratch
USER 1000
COPY --from=build --chown=1000 --chmod=755 /build/oidc-exporter /
CMD ["/oidc-exporter"]
