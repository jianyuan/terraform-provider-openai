package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (r *AdminApiKeyResource) getCreateJSONRequestBody(ctx context.Context, data AdminApiKeyResourceModel) (apiclient.AdminApiKeysCreateJSONRequestBody, diag.Diagnostics) {
	return apiclient.AdminApiKeysCreateJSONRequestBody{
		Name: data.Name.ValueString(),
	}, nil
}

func (m *AdminApiKeyResourceModel) Fill(ctx context.Context, data any) diag.Diagnostics {
	switch v := data.(type) {
	case apiclient.AdminApiKey:
		m.Id = supertypes.NewStringValue(v.Id)
		m.Name = supertypes.NewStringPointerValue(v.Name)
		m.CreatedAt = supertypes.NewInt64Value(v.CreatedAt)
	case apiclient.AdminApiKeyCreateResponse:
		m.Id = supertypes.NewStringValue(v.Id)
		m.Name = supertypes.NewStringPointerValue(v.Name)
		m.CreatedAt = supertypes.NewInt64Value(v.CreatedAt)
		m.ApiKey = supertypes.NewStringValue(v.Value)
	default:
		var diags diag.Diagnostics
		diags.AddError("Unknown type", fmt.Sprintf("Unknown type: %T", data))
		return diags
	}
	return nil
}
