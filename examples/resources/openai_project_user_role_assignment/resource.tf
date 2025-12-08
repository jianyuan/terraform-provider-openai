# Create a project
resource "openai_project" "test" {
  name = "Project name"
}

# Create a project role
resource "openai_project_role" "test" {
  project_id  = openai_project.test.id
  name        = "API Project Key Manager"
  description = "Allows managing API keys for the project"
  permissions = [
    "api.organization.projects.api_keys.read",
    "api.organization.projects.api_keys.write"
  ]
}

# Assign a user to a project
resource "openai_project_user" "test" {
  project_id = openai_project.test.id
  user_id    = "user_abc123"
  role       = "member"
}

# Assign a role to a user in a project
resource "openai_project_user_role_assignment" "test" {
  project_id = openai_project_user.test.project_id
  user_id    = openai_project_user.test.user_id
  role_id    = openai_project_role.test.id
}
