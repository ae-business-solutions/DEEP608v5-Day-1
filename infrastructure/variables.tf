# Local Variables
# https://www.terraform.io/docs/language/values/locals.html

locals {
  # GCP & Kubernetes
  gcp_region    = "us-central1"
  gcp_zone      = "us-central1-c"
  gcp_project   = var.gcp_project_id
  gke_name      = "deep608v5-01"
  gke_node_size = "e2-medium"
  redis_name    = "redis"
  acme_contact  = "deep608v5-lab@aebs.com"
  
  # Services
  block_page_dns      = "block.edl.${var.domain}"
  request_unblock_dns = "request.edl.${var.domain}"
  edl_admin_dns       = "admin.edl.${var.domain}"
  edl_dns             = "list.edl.${var.domain}"

  # Okta
  okta_auth_url = "https://${var.okta_domain}/oauth2/default"
}

# Input Variables
# https://www.terraform.io/docs/language/values/variables.html

# Lab Variables
variable "gcp_project_id" {
  type        = string
  description = "GCP Project Id"
}
variable "domain" {
  type        = string
  description = "Lab Domain"
}

# Okta Variables
variable "okta_domain" {
  type        = string
  description = "Okta Domain"
}
variable "okta_user_app_client_id" {
  type        = string
  description = "Okta User App Client Id"
}
variable "okta_user_app_client_secret" {
  type        = string
  description = "Okta User App Client Secret"
}
variable "okta_admin_app_client_id" {
  type        = string
  description = "Okta Admin App Client Id"
}
variable "okta_admin_app_client_secret" {
  type        = string
  description = "Okta Admin App Client Secret"
}