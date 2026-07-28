resource "openai_spend_limit" "example" {
  currency         = "USD"
  interval         = "month"
  threshold_amount = 10000
}
