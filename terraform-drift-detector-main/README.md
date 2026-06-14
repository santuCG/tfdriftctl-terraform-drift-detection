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
The API server runs securely over HTTPS. Generate a local self-signed certificate by running:
```bash
go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host localhost
```
*(For Windows PowerShell, you can use the provided script: `.\scripts\generate-cert.ps1`)*

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
Compile the server and the CLI tools:
```bash
# Build the API Server
go build -o bin/drift-server ./cmd/drift-server

# Build the CLI
go build -o bin/tfdriftctl ./cmd/tfdriftctl
```

### 4. Start the Server
Run the API server in your terminal:
```bash
./bin/drift-server -config configs/tfdriftctl.yaml
```

### 5. Authenticate & Trigger a Scan
In a new terminal window, authenticate with your password to receive a JWT token:
```bash
# Login
curl -k -X POST https://localhost:8443/api/v1/login -d '{"password": "YOUR_ADMIN_PASSWORD"}'
```

Copy the `token` from the response, and use it to trigger a drift scan for your workspace:
```bash
# Replace YOUR_TOKEN with the JWT from the previous step
curl -k -H "Authorization: Bearer YOUR_TOKEN" -X POST https://localhost:8443/api/v1/workspaces/aws-s3-test/scans
```

---

## CLI Usage (Ad-hoc Scans)

If you don't want to run the background server, you can use the CLI tool to perform instant, ad-hoc scans directly against a state file:

```bash
./bin/tfdriftctl scan --state path/to/terraform.tfstate --provider aws --region us-east-1
```
*The CLI will output a table showing exactly what resources are missing, extra, or modified in the live cloud environment.*
