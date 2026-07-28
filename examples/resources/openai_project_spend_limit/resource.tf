resource "openai_project_spend_limit" "example" {
  project_id       = "proj_000000000000000000000000"
  currency         = "USD"
  interval         = "month"
  threshold_amount = 10000
}
