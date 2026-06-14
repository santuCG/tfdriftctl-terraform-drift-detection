$ErrorActionPreference = "Stop"
Write-Host "Generating self-signed certificate for localhost..."
$goRoot = go env GOROOT
go run "$goRoot/src/crypto/tls/generate_cert.go" --host localhost
Write-Host "Certificate generated successfully. cert.pem and key.pem are now in the root directory."