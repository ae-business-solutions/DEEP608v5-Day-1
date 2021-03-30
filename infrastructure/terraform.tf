# Setup Terraform & Connect to Terraform Cloud for state storage
# https://www.terraform.io/docs/language/settings/index.html
# https://www.terraform.io/docs/language/settings/backends/remote.html
terraform {
  required_version = ">= 0.14"
  backend "remote" {
    organization = "DEEP608v5"
    workspaces {
      name = "<Insert your Terraform Cloud Workspace Name Here>"
    }
  }
}