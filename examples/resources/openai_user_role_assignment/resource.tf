resource "openai_user_role_assignment" "test" {
  user_id = "user_abc123"
  role_id = "role_01J1F8ROLE01"
}

# Note that prebuilt roles are in the format of role-<role_name>__<resource_type>__<resource_id>
locals {
  organization_id = "org-123"
}

# Assign prebuilt owner role to a user
resource "openai_user_role_assignment" "owner" {
  user_id = "user_abc123"
  role_id = "role-api-organization-owner__api-organization__${local.organization_id}"
}

# Assign prebuilt reader role to a user
resource "openai_user_role_assignment" "reader" {
  user_id = "user_abc123"
  role_id = "role-api-organization-reader__api-organization__${local.organization_id}"
}
