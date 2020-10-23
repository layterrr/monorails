terraform {
  backend "gcs" {
    bucket  = "terraform-admin"
    prefix  = "terraform/state"
  }
}
