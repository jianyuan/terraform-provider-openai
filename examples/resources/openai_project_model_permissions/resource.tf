resource "openai_project_model_permissions" "example" {
  project_id = "proj_000000000000000000000000"
  mode       = "allow_list"
  model_ids  = ["gpt-4.1", "o3"]
}
