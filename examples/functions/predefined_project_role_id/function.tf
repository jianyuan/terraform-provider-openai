locals {
  # owner role id of project proj_abc123
  owner_role_id = provider::openai::predefined_project_role_id("owner", "proj_abc123")

  # member role id of project proj_abc123
  member_role_id = provider::openai::predefined_project_role_id("member", "proj_abc123")

  # viewer role id of project proj_abc123
  viewer_role_id = provider::openai::predefined_project_role_id("viewer", "proj_abc123")
}
