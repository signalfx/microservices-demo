$ErrorActionPreference = "Stop"
$repoRoot = Split-Path -Parent $PSScriptRoot
$docker = Get-Command docker -ErrorAction SilentlyContinue
if ($null -eq $docker) {
    $dockerPath = Join-Path $env:ProgramFiles "Docker/Docker/resources/bin/docker.exe"
    if (-not (Test-Path -LiteralPath $dockerPath)) {
        throw "docker was not found on PATH or in the default Docker Desktop location"
    }
} else {
    $dockerPath = $docker.Source
}

& $dockerPath run --rm `
    --mount "type=bind,source=$repoRoot,target=/workspace" `
    --workdir /workspace `
    alpine:3.24 `
    sh -c "apk add --no-cache openssl >/dev/null && sh ./hack/generate-grpc-tls.sh"

if ($LASTEXITCODE -ne 0) {
    throw "gRPC TLS certificate generation failed with exit code $LASTEXITCODE"
}
