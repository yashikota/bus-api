output "service_account_email" {
  value       = google_service_account.main.email
  description = "Service Account Email"
}

output "wif_provider_id" {
  value       = google_iam_workload_identity_pool_provider.main.name
  description = "Workload Identity Pool Provider ID"
}
