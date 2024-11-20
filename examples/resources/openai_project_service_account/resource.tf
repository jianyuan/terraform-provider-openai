resource "openai_project" "test" {
  name = "my-project"
}

resource "openai_project_service_account" "test" {
  project_id = openai_project.test.id
  name       = "my-service-account"
}

output "service_account_api_key" {
  sensitive = true
  value     = openai_project_service_account.test.api_key
}
