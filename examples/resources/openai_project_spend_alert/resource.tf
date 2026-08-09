resource "openai_project_spend_alert" "test" {
  project_id       = "proj_000000000000000000000000"
  currency         = "USD"
  interval         = "month"
  threshold_amount = 1000000
  notification_channel = {
    type           = "email"
    recipients     = ["finance@example.com"]
    subject_prefix = "OpenAI Spend Alert"
  }
}
