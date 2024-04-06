resource "google_iam_workload_identity_pool" "main" {
  workload_identity_pool_id = "github"
  display_name              = "GitHub Actions Pool"
  disabled                  = false
  project                   = var.project_id

  depends_on = [module.project-services]
}

resource "google_iam_workload_identity_pool_provider" "main" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.main.workload_identity_pool_id
  workload_identity_pool_provider_id = "github"
  display_name                       = "GitHub Actions Pool Provider"
  disabled                           = false
  attribute_condition                = "assertion.repository_owner == \"${var.github_repo_owner}\""
  attribute_mapping = {
    "google.subject" = "assertion.repository"
  }
  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
  project = var.project_id

  depends_on = [google_iam_workload_identity_pool.main]
}

resource "google_service_account_iam_member" "workload_identity_sa_iam" {
  service_account_id = google_service_account.main.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principal://iam.googleapis.com/${google_iam_workload_identity_pool.main.name}/subject/${var.github_repo_owner}/${var.github_repo_name}"

  depends_on = [google_service_account.main, google_iam_workload_identity_pool_provider.main]
}

resource "google_service_account" "main" {
  account_id   = var.github_repo_name
  display_name = var.github_repo_name
  description  = "GitHub Actions Service Account"
  project      = var.project_id

  depends_on = [google_iam_workload_identity_pool_provider.main]
}
