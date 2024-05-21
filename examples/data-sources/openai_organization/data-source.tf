// Get the default organization
data "openai_organization" "default" {
}

// Get an organization by ID
data "openai_organization" "example" {
  id = "org-000000000000000000000000"
}
