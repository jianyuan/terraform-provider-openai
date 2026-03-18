resource "openai_project" "example" {
  name = "Example Project"
}

resource "openai_project" "eu_project" {
  name      = "EU Data Residency Project"
  geography = "EU"
}
