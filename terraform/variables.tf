variable "billing_account_id" {
  type = string
}

variable "project_id" {
  type = string
}

variable "region" {
  type    = string
  default = "asia-northeast1"
}

variable "zone" {
  type    = string
  default = "asia-northeast1-b"
}

variable "cr_name" {
  type        = string
  description = "cloud run service name"
}

variable "ar_image_url" {
  type        = string
  description = "artifact registry image url"
}

variable "ar_repository_name" {
  type = string
}

variable "ar_repository_description" {
  type = string
}

variable "github_repo_owner" {
  type    = string
  default = "yashikota"
}

variable "github_repo_org" {
  type = string
}

variable "github_repo_name" {
  type = string
}
