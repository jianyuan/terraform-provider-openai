resource "openai_spend_alert" "test" {
  currency         = "USD"
  interval         = "month"
  threshold_amount = 1000000
  notification_channel = {
    type           = "email"
    recipients     = ["finance@example.com"]
    subject_prefix = "OpenAI Spend Alert"
  }
}
