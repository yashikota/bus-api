resource "google_artifact_registry_repository" "my-repo" {
  location      = var.region
  repository_id = var.ar_repository_name
  description   = var.ar_repository_description
  format        = "DOCKER"

  depends_on = [module.project-services]
}
