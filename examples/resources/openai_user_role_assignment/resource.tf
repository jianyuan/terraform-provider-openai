resource "openai_user_role_assignment" "test" {
  user_id = "user_abc123"
  role_id = "role_01J1F8ROLE01"
}

# Note that prebuilt roles are in the format of role-<role_name>__<resource_type>__<resource_id>
# You can use the function `provider::openai::predefined_role_id` to generate the role_id

# Assign prebuilt owner role to a user
resource "openai_user_role_assignment" "owner" {
  user_id = "user_abc123"
  role_id = provider::openai::predefined_role_id("owner", "org-123") # role-api-organization-owner__api-organization__org-123
}

# Assign prebuilt reader role to a user
resource "openai_user_role_assignment" "reader" {
  user_id = "user_abc123"
  role_id = provider::openai::predefined_role_id("reader", "org-123") # role-api-organization-reader__api-organization__org-123
}
