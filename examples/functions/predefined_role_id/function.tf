locals {
  # owner role id of organization org-123
  owner_role_id = provider::openai::predefined_role_id("owner", "org-123")

  # reader role id of organization org-123
  reader_role_id = provider::openai::predefined_role_id("reader", "org-123")
}

