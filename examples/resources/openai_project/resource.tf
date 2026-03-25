resource "openai_project" "example" {
  name = "Example Project"

  # Optional: Create the project with the specified data residency region.
  # Valid values: US, EU, JP, IN, KR, CA, AU, SG
  geography = "US"
}
