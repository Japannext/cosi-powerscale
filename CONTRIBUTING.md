# Environment

Use [gvm](https://github.com/moovweb/gvm) to match the project's golang version,
or use a version of golang from your system that matches the one indicated in `go.mod`.

```bash
gvm use go1.22
```

# Building

Verify the driver builds:
```bash
./task build
```

# Deploying to a dev environment

You can specify a helmfile in `.helmfile.yaml`, and it will be executed at the end of the build process.
Here is a relevant example:
```yaml
# .helmfile.yaml
---
environments:
  default:
    kubeContext: my_dev_env
---
releases:
- name: nas1
  namespace: kube-storage
  createNamespace: true
  values:
  - provisioner:
      repository: nexus.example.com:8080/kubernetes/cosi-powerscale
      tag: develop
      pullPolicy: Always
    config:
      name: nas1
      apiEndpoint: https://isilon1.example.com:8080
      apiSecret: nas1-api-credentials
      basePath: /ifs/kubernetes/development
      s3Endpoint: https://data.nas1.example.com:9021
      zone: examplezone001
      region: ""
      tlsCacertConfigMap: ca-bundle
```

It can then be deployed with:
```bash
./task sync
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

On Ubuntu:
```bash
cp /usr/local/share/ca-certificates/* .ca-bundle/
```
On RHEL:
```bash
# On RHEL
cp /etc/pki/ca-trust/source/anchors/* .ca-bundle/
```

It will be added to the docker intermediate build image to fetch
dependencies, but not to the final image.
