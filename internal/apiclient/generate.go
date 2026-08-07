package apiclient

//go:generate curl -s -o api.json https://raw.githubusercontent.com/openai/openai-openapi/refs/heads/main/openapi.json
//go:generate sh -c "jq -f api-mods.jq api.json > api.tmp && mv api.tmp api.json"
//go:generate go tool oapi-codegen -config config.yaml api.json
