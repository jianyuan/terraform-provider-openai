terraform {
  required_providers {
    openai = {
      source = "jianyuan/openai"
    }
  }
}

# Configure the OpenAI provider
provider "openai" {
  admin_key = "sk-admin-0000000000000000000000000000000000000000"
}

# Create a project
resource "openai_project" "example" {
  name = "Example Project"
}

# Create a service account for the project
resource "openai_project_service_account" "example" {
  project_id = openai_project.example.id
  name       = "my-service-account"
}

# Output the API key for the service account
output "service_account_api_key" {
  sensitive = true
  value     = openai_project_service_account.example.api_key
}
