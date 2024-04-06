module "project-services" {
  source  = "terraform-google-modules/project-factory/google//modules/project_services"
  version = "~> 14.5"

  project_id = var.project_id
  activate_apis = [
    "servicemanagement.googleapis.com",
    "iam.googleapis.com",
    "cloudapis.googleapis.com",
    "artifactregistry.googleapis.com",
    "run.googleapis.com",
  ]

  depends_on = [google_project.my_project]
}
