variable "project_name" {}
variable "project_id" {}
variable "billing_account" {}
variable "org_id" {}
variable "region" {}
variable "environment" {}
variable "container_project" {}

provider "google" {
  region = var.region
  project = var.project_id
  version = "3.23.0"
}

provider "google-beta" {
  region = var.region
  project = var.project_id
  version = "3.23.0"
}

resource "google_project" "project" {
  provider  = google
  name            = var.project_name
  project_id      = var.project_id
  billing_account = var.billing_account
  org_id          = var.org_id
}

resource "google_project_service" "service" {
  for_each = toset([
    "appengine.googleapis.com",
    "bigquery.googleapis.com",
    "bigquerystorage.googleapis.com",
    "cloudapis.googleapis.com",
    "cloudbilling.googleapis.com",
    "cloudbuild.googleapis.com",
    "clouddebugger.googleapis.com",
    "cloudfunctions.googleapis.com",
    "cloudkms.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "cloudscheduler.googleapis.com",
    "cloudshell.googleapis.com",
    "cloudtasks.googleapis.com",
    "cloudtrace.googleapis.com",
    "compute.googleapis.com",
    "container.googleapis.com",
    "containerregistry.googleapis.com",
    "dataflow.googleapis.com",
    "datastore.googleapis.com",
    "deploymentmanager.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
    "iap.googleapis.com",
    "logging.googleapis.com",
    "monitoring.googleapis.com",
    "oslogin.googleapis.com",
    "pubsub.googleapis.com",
    "resourceviews.googleapis.com",
    "run.googleapis.com",
    "secretmanager.googleapis.com",
    "servicedirectory.googleapis.com",
    "servicemanagement.googleapis.com",
    "servicenetworking.googleapis.com",
    "serviceusage.googleapis.com",
    "sql-component.googleapis.com",
    "sqladmin.googleapis.com",
    "stackdriver.googleapis.com",
    "storage-api.googleapis.com",
    "storage-component.googleapis.com",
  ])

  service = each.key

  project            = google_project.project.project_id
  disable_on_destroy = false
}

resource "google_project_iam_member" "container_access" {
  project = var.container_project
  role    = "roles/storage.objectViewer"
  member = "serviceAccount:service-${google_project.project.number}@serverless-robot-prod.iam.gserviceaccount.com"
}

resource "google_project_iam_binding" "cloudsql" {
  project = google_project.project.project_id
  role    = "roles/cloudsql.client"

  members = [
    "serviceAccount:service-${google_project.project.number}@serverless-robot-prod.iam.gserviceaccount.com"
  ]
}

resource "google_project_iam_binding" "cloudrun_admin_cloud_build" {
  project = google_project.project.project_id
  role    = "roles/run.admin"

  members = [
    "serviceAccount:${google_project.project.number}@cloudbuild.gserviceaccount.com"
  ]
}

resource "google_project_iam_binding" "cloudrun_storage_token_creator" {
  project = google_project.project.project_id
  role    = "roles/iam.serviceAccountTokenCreator"

  members = [
    "serviceAccount:${google_project.project.number}-compute@developer.gserviceaccount.com"
  ]
}

data "google_compute_default_service_account" "default" {
  project = google_project.project.project_id
}

resource "google_service_account_iam_member" "cloud_build_run_as_cloud_run" {
  service_account_id = data.google_compute_default_service_account.default.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_project.project.number}@cloudbuild.gserviceaccount.com"
}

resource "google_project_iam_binding" "cloud_build_service_account_access" {
  project = google_project.project.project_id
  role    = "roles/cloudbuild.builds.builder"

  members = [
    "serviceAccount:${google_project.project.number}@cloudbuild.gserviceaccount.com",
    "serviceAccount:github-actions@${container_project}.iam.gserviceaccount.com", # allow github actions to create cloud builds
  ]
}

resource "google_project_iam_binding" "cloud_build_logs_access" {
  project = google_project.project.project_id
  role    = "roles/viewer"

  members = [
    "serviceAccount:${google_project.project.number}@cloudbuild.gserviceaccount.com",
    "serviceAccount:github-actions@${container_project}.iam.gserviceaccount.com", # allow github actions to view cloud build logs
  ]
}

resource "google_project_iam_binding" "cloud_sql_editor" {
  project = google_project.project.project_id
  role    = "roles/cloudsql.editor"

  members = [
    "serviceAccount:github-actions@${container_project}.iam.gserviceaccount.com", # allow github actions to create database backups
  ]
}

output "project_id" {
  value = google_project.project.project_id
}
