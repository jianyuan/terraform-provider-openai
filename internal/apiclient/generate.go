package apiclient

//go:generate curl -sSfL -o api-original.yaml https://raw.githubusercontent.com/openai/openai-openapi/refs/heads/master/openapi.yaml
//go:generate sh -c "go run carvel.dev/ytt/cmd/ytt -f api-original.yaml -f api-patches.yaml > api.yaml"
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config config.yaml api.yaml
