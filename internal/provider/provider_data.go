package provider

import (
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	"github.com/openai/openai-go/v3"
)

type providerData struct {
	client   *apiclient.ClientWithResponses
	clientV2 *openai.Client
}
