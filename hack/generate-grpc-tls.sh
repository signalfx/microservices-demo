#!/bin/sh

# Generates a local CA and short-lived service certificates for the demo's
# encrypted internal gRPC endpoints. The resulting manifest contains private
# keys and is intentionally ignored by git.

set -eu
umask 077

output="${1:-kubernetes-manifests/grpc-tls.generated.yaml}"
workdir="$(mktemp -d)"
trap 'rm -rf "$workdir"' EXIT

if ! command -v openssl >/dev/null 2>&1; then
  echo "openssl is required (or run hack/generate-grpc-tls.ps1 to use Docker)" >&2
  exit 1
fi

openssl ecparam -name prime256v1 -genkey -noout -out "$workdir/ca.key"
openssl req -x509 -new -sha256 -key "$workdir/ca.key" -days 365 \
  -subj "/CN=microservices-demo gRPC CA" \
  -addext "basicConstraints=critical,CA:TRUE" \
  -addext "keyUsage=critical,keyCertSign,cRLSign" \
  -out "$workdir/ca.crt"

for service in checkoutservice productcatalogservice shippingservice; do
  openssl ecparam -name prime256v1 -genkey -noout -out "$workdir/$service.key"
  openssl req -new -sha256 -key "$workdir/$service.key" \
    -subj "/CN=$service" -out "$workdir/$service.csr"
  cat > "$workdir/$service.ext" <<EOF
basicConstraints=critical,CA:FALSE
keyUsage=critical,digitalSignature
extendedKeyUsage=serverAuth
subjectAltName=DNS:$service,DNS:$service.default,DNS:$service.default.svc,DNS:$service.default.svc.cluster.local
EOF
  openssl x509 -req -sha256 -in "$workdir/$service.csr" \
    -CA "$workdir/ca.crt" -CAkey "$workdir/ca.key" -CAcreateserial \
    -days 30 -extfile "$workdir/$service.ext" -out "$workdir/$service.crt"
done

mkdir -p "$(dirname "$output")"
{
  cat <<'EOF'
apiVersion: v1
kind: ConfigMap
metadata:
  name: grpc-tls-ca
data:
  ca.crt: |
EOF
  sed 's/^/    /' "$workdir/ca.crt"

  for service in checkoutservice productcatalogservice shippingservice; do
    cat <<EOF
---
apiVersion: v1
kind: Secret
metadata:
  name: $service-grpc-tls
type: kubernetes.io/tls
stringData:
  tls.crt: |
EOF
    sed 's/^/    /' "$workdir/$service.crt"
    cat <<'EOF'
  tls.key: |
EOF
    sed 's/^/    /' "$workdir/$service.key"
  done
} > "$output"

echo "Generated $output (contains private keys; do not commit it)."
