FROM golang:1.22 as builder

WORKDIR /code

# Copy the Go Modules manifests.
COPY go.mod go.mod
COPY go.sum go.sum

# Local CA use
COPY .ca-bundle /usr/local/share/ca-certificates/
RUN chmod -R 644 /usr/local/share/ca-certificates/ && update-ca-certificates

# Dependency download
RUN go mod download

# Copy the go source.
COPY main.go main.go
COPY pkg/ pkg/

# Build
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o /build/cosi-powerscale .

# Second stage: building final environment for running the driver.
FROM scratch
LABEL org.opencontainers.image.source https://github.com/japannext/cosi-powerscale

COPY --from=builder /build/cosi-powerscale /app/cosi-powerscale

USER 1001
WORKDIR /var/lib/cosi

# Set the entrypoint.
ENTRYPOINT ["/app/cosi-powerscale"]
CMD []
