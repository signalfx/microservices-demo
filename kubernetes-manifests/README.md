# ./kubernetes-manifests

:warning: Kubernetes manifests provided in this directory are not directly
deployable to a cluster. They are meant to be used with `skaffold` command to
insert the correct `image:` tags.

Use the manifests in [/release](/release) directory which are configured with
pre-built public images.

## Internal gRPC TLS

Checkout, product catalog, and shipping require TLS certificates. Generate a
local CA and one certificate per service before running Skaffold:

```powershell
./hack/generate-grpc-tls.ps1
skaffold run
```

On Linux or macOS, run `hack/generate-grpc-tls.sh` with OpenSSL installed. The
generated `kubernetes-manifests/grpc-tls.generated.yaml` contains private keys,
is ignored by git, and should be replaced by your normal certificate issuer and
Secret-management process outside local demo environments.
