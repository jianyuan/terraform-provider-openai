# Create an API key for the default project
resource "openai_project_api_key" "example" {
  service_account_id = "my-service-account"
}

# Create an API key for a specific project
resource "openai_project_api_key" "example" {
  project_id         = "proj_000000000000000000000000"
  service_account_id = "my-service-account"
}

# Create a read-only API key
resource "openai_project_api_key" "example" {
  service_account_id = "my-service-account"
  read_only          = true
}

# Create an API key with specific permissions
resource "openai_project_api_key" "example" {
  service_account_id = "my-service-account"
  permissions {
    models             = "read"
    model_capabilities = "write"
  }
}
