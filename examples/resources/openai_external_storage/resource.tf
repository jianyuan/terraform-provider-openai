resource "openai_project" "example" {
  name = "Example Project"
}

# Store the project's data in a Google Cloud Storage bucket. Use `aws`, `azure` or `oci`
# instead for the other clouds.
resource "openai_external_storage" "example" {
  project_id = openai_project.example.id

  gcp = {
    bucket                           = "example-openai-data"
    workload_identity_project_number = "123456789012"
    workload_identity_pool_id        = "openai-pool"
    workload_identity_provider_id    = "openai-provider"
  }
}
