resource "openai_project_user" "example" {
  project_id = "proj_000000000000000000000000"
  user_id    = "user-000000000000000000000000"
  role       = "member"
}
