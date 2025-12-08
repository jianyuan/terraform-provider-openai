resource "openai_project" "test" {
  name = "Project name"
}

resource "openai_project_role" "test" {
  project_id  = openai_project.test.id
  name        = "API Project Key Manager"
  description = "Allows managing API keys for the project"
  permissions = [
    "api.organization.projects.api_keys.read",
    "api.organization.projects.api_keys.write"
  ]
}

resource "openai_group" "test" {
  name = "Support Team"
}

resource "openai_project_group_role_assignment" "test" {
  project_id = openai_project.test.id
  group_id   = openai_group.test.id
  role_id    = openai_project_role.test.id
}
