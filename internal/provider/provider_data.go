package provider

import (
	"github.com/openai/openai-go/v3"
)

type providerData struct {
	client *openai.Client
}
