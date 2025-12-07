resource "openai_organization_role" "test" {
  name        = "API Group Manager"
  description = "Allows managing organization groups"
  permissions = ["api.groups.read", "api.groups.write"]
}
