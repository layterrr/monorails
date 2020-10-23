module "service-loan" {
  source                  = "../services/loan"
  container_project       = var.container_project
  environment             = var.environment
  project_id              = google_project.project.project_id
  project_number          = google_project.project.number
  region                  = var.region
  private_network         = google_compute_network.private_network.self_link
  service_depends_on      = [google_service_networking_connection.private_vpc_connection, google_project_iam_member.container_access, google_project_iam_binding.cloudsql]
}