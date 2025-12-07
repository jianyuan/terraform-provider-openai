resource "openai_project_role" "test" {
  project_id  = "proj_000000000000000000000000"
  name        = "API Project Key Manager"
  description = "Allows managing API keys for the project"
  permissions = [
    "api.organization.projects.api_keys.read",
    "api.organization.projects.api_keys.write",
  ]
}
