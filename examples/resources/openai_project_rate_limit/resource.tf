resource "openai_project_rate_limit" "example" {
  project_id    = "proj_000000000000000000000000"
  rate_limit_id = "rl-o1-preview"

  max_requests_per_1_minute = 2
  max_tokens_per_1_minute   = 75000
}
