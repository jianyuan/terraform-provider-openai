package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (r *AdminApiKeyResource) getCreateJSONRequestBody(ctx context.Context, data AdminApiKeyResourceModel) (apiclient.AdminApiKeysCreateJSONRequestBody, diag.Diagnostics) {
	return apiclient.AdminApiKeysCreateJSONRequestBody{
		Name: data.Name.ValueString(),
	}, nil
}

func (m *AdminApiKeyResourceModel) Fill(ctx context.Context, data apiclient.AdminApiKey) diag.Diagnostics {
	m.Id = supertypes.NewStringValue(data.Id)
	m.Name = supertypes.NewStringPointerValue(data.Name)
	m.CreatedAt = supertypes.NewInt64Value(data.CreatedAt)
	return nil
}

func (m *AdminApiKeyResourceModel) FillFromCreate(ctx context.Context, data apiclient.AdminApiKeyCreateResponse) diag.Diagnostics {
	m.Id = supertypes.NewStringValue(data.Id)
	m.Name = supertypes.NewStringPointerValue(data.Name)
	m.CreatedAt = supertypes.NewInt64Value(data.CreatedAt)
	m.ApiKey = supertypes.NewStringValue(data.Value)
	return nil
}
