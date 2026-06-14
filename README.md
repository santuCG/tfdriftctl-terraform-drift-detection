# tfdriftctl - Terraform Drift Detection

`tfdriftctl` is a tool that continuously compares your Terraform state files against live cloud infrastructure (AWS) to detect configuration drift—without needing to run `terraform plan` or `terraform apply`. 

It comes with a REST API, JWT Authentication, and TLS encryption out-of-the-box.

## How It Works (Architecture)

```text
Terraform State -> State Reader -> Extractor -> Expected Model
                                                        |
                                                        v
Cloud APIs      -> Cloud Fetcher -> Extractor -> Actual Model  -> Drift Engine -> Report
```

`tfdriftctl` reads your "Expected Model" directly from your Terraform state file and pulls your "Actual Model" live from AWS Cloud APIs. The Drift Engine then compares the two models attribute-by-attribute to find any unmanaged discrepancies, alerting you to security and configuration risks before your next `terraform apply`!
---

## Prerequisites
- **Go** (1.20+)
- **AWS CLI** configured with your credentials (`aws configure`)

---

## 🚀 Quick Start Guide

Follow these steps in order to run `tfdriftctl` on your local machine.

### 1. Generate TLS Certificates
The API server runs securely over HTTPS. Generate a local self-signed certificate by running the following in your terminal:

**On Mac / Linux:**
```bash
go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host localhost
```

**On Windows (PowerShell):**
```powershell
$goRoot = go env GOROOT
go run "$goRoot\src\crypto\tls\generate_cert.go" --host localhost
```

This will generate `cert.pem` and `key.pem` in your project root.

### 2. Configure Secrets and Workspaces
Open `configs/tfdriftctl.yaml` and update the security credentials:
- Change `<CHANGE_ME_JWT_SECRET>` to a secure random string.
- Change `<CHANGE_ME_ADMIN_PASSWORD>` to your preferred login password.

Make sure your `workspaces` block points to a valid local `terraform.tfstate` file:
```yaml
workspaces:
  - name: aws-s3-test
    provider: aws
    state:
      backend: local
      path: path/to/your/terraform.tfstate
    regions:
      - us-east-1
```

### 3. Build the Application
Compile the server and the CLI tools (if you are on Windows, append `.exe` to the output files as shown below so Windows knows they are runnable programs):

```bash
# Build the API Server
go build -o bin/drift-server.exe ./cmd/drift-server

# Build the CLI
go build -o bin/tfdriftctl.exe ./cmd/tfdriftctl
```

### 4. Start the Server
Run the API server in your terminal:
```bash
./bin/drift-server.exe -config configs/tfdriftctl.yaml
```

### 5. Authenticate & Trigger a Scan
In a new terminal window, authenticate with your password to receive a JWT token:
```bash
# 1. Login to get your token
curl -k -X POST https://localhost:8443/api/v1/login -d '{"password": "YOUR_ADMIN_PASSWORD"}'

# 2. List your workspaces to get your internal Workspace ID
curl -k -H "Authorization: Bearer YOUR_TOKEN" https://localhost:8443/api/v1/workspaces

# 3. Trigger a scan using the Workspace ID from the previous step
curl -k -H "Authorization: Bearer YOUR_TOKEN" -X POST https://localhost:8443/api/v1/workspaces/YOUR_WORKSPACE_ID/scans
```

*Note: The `POST /scans` endpoint triggers a background scan and returns the raw JSON findings. If you want to view a pretty table format in the terminal, use the CLI instead (see below), or fetch the report manually via `GET /api/v1/scans/SCAN_ID/report?format=table`.*

---

## CLI Usage (Ad-hoc Scans)

If you don't want to run the background server, you can use the CLI tool to perform instant, ad-hoc scans directly against a state file:

```bash
./bin/tfdriftctl.exe scan --state path/to/terraform.tfstate --provider aws --region us-east-1
```
*The CLI will output a table showing exactly what resources are missing, extra, or modified in the live cloud environment.*
