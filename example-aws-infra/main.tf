terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# Configure the AWS Provider
provider "aws" {
  region = "us-east-1"
}

# Create a randomized bucket name to avoid conflicts
resource "random_id" "bucket_suffix" {
  byte_length = 4
}

# Create a basic S3 Bucket
resource "aws_s3_bucket" "test_bucket" {
  bucket = "tfdriftctl-test-bucket-${random_id.bucket_suffix.hex}"

  tags = {
    Environment = "Test"
    ManagedBy   = "Terraform"
    Project     = "tfdriftctl"
  }
}

# Add Basic Versioning
resource "aws_s3_bucket_versioning" "test_bucket_versioning" {
  bucket = aws_s3_bucket.test_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}

# Output the bucket name and the region so we know what to look for
output "bucket_name" {
  value = aws_s3_bucket.test_bucket.id
}

output "aws_region" {
  value = "us-east-1"
}
