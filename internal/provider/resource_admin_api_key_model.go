package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (r *AdminApiKeyResource) getCreateJSONRequestBody(data AdminApiKeyResourceModel) apiclient.AdminApiKeysCreateJSONRequestBody {
	return apiclient.AdminApiKeysCreateJSONRequestBody{
		Name: data.Name.ValueString(),
	}
}

func (m *AdminApiKeyResourceModel) Fill(ctx context.Context, data apiclient.AdminApiKey) diag.Diagnostics {
	m.Id = supertypes.NewStringValue(data.Id)
	m.Name = supertypes.NewStringValue(data.Name)
	m.CreatedAt = supertypes.NewInt64Value(data.CreatedAt)

	if data.Value != nil {
		m.ApiKey = supertypes.NewStringPointerValue(data.Value)
	}

	return nil
}
