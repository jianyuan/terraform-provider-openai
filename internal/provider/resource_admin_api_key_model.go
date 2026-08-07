package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (r *AdminApiKeyResource) getCreateJSONRequestBody(ctx context.Context, data AdminApiKeyResourceModel) (*openai.AdminOrganizationAdminAPIKeyNewParams, diag.Diagnostics) {
	return &openai.AdminOrganizationAdminAPIKeyNewParams{
		Name: data.Name.ValueString(),
	}, nil
}

func (m *AdminApiKeyResourceModel) Fill(ctx context.Context, data any) diag.Diagnostics {
	switch v := data.(type) {
	case openai.AdminAPIKey:
		m.Id = supertypes.NewStringValue(v.ID)
		m.Name = supertypes.NewStringValue(v.Name)
		m.CreatedAt = supertypes.NewInt64Value(v.CreatedAt)
	case openai.AdminOrganizationAdminAPIKeyNewResponse:
		m.Id = supertypes.NewStringValue(v.ID)
		m.Name = supertypes.NewStringValue(v.Name)
		m.CreatedAt = supertypes.NewInt64Value(v.CreatedAt)
		m.ApiKey = supertypes.NewStringValue(v.Value)
	default:
		var diags diag.Diagnostics
		diags.AddError("Unknown type", fmt.Sprintf("Unknown type: %T", data))
		return diags
	}

	return nil
}
