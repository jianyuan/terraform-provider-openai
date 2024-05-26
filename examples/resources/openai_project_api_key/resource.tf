# Create an API key for the default project
resource "openai_project_api_key" "example" {
  organization_id    = "org-000000000000000000000000"
  service_account_id = "my-service-account"
}

# Create an API key for a specific project
resource "openai_project_api_key" "example" {
  organization_id    = "org-000000000000000000000000"
  project_id         = "proj_000000000000000000000000"
  service_account_id = "my-service-account"
}
