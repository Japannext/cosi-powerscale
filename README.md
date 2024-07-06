# COSI driver for Dell Powerscale (Isilon/OneFS)

Dell Powerscale is ahardware solution specialized in shared filesystem (NFS/CIFS).
However, since recent versions, they do support basic S3-like bucket functionalities.

This is a COSI driver (Container Object Storage Interface) in order to automate the creation of buckets in Powerscale.

# Installation

```bash
helm install nas1 oci://ghcr.io/japannext/helm-charts/cosi-powerscale --version 1.1.0 --values values.yaml
```

`values.yaml` example:
```yaml
---
config:
  name: nas1
  apiEndpoint: https://isilon1.example.comL8080
  # A basicauth secret
  apiSecret: nas1-api-credentials
  basePath: /ifs/kubernetes/production
  zone: examplezone001
  s3Endpoint: https://data.nas1.example.com:9021
  region: ""
  # A ConfigMap with a "ca.crt" key
  tlsCacertConfigMap: ca-bundle
```

Don't forget to populate the secret for isilon API:
```bash
kubectl create secret generic nas1-api-credentials --from-field=username=root --from-field=password=password123
```

Install the example bucket to test if it's working:
```bash
kubectl apply -f examples/
```

Verify it creates a Bucket object and a Secret:
```bash
kubectl get bucket,bucketclaim,bucketaccess
```

Check the secret for accessing the bucket:
```bash
kubectl get secret example-bucket-secret -o jsonpath="{.data.BucketInfo}" | base64 -d
```
