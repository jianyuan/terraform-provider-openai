package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectServiceAccountResourceModel) Fill(ctx context.Context, data any) diag.Diagnostics {
	switch data := data.(type) {
	case openai.AdminOrganizationProjectServiceAccountNewResponse:
		m.Id = supertypes.NewStringValue(data.ID)
		m.Name = supertypes.NewStringValue(data.Name)
		m.Role = supertypes.NewStringValue(string(data.Role))
		m.CreatedAt = supertypes.NewInt64Value(data.CreatedAt)
		m.ApiKeyId = supertypes.NewStringValue(data.APIKey.ID)
		m.ApiKey = supertypes.NewStringValue(data.APIKey.Value)
		return nil
	case openai.ProjectServiceAccount:
		m.Id = supertypes.NewStringValue(data.ID)
		m.Name = supertypes.NewStringValue(data.Name)
		m.Role = supertypes.NewStringValue(string(data.Role))
		m.CreatedAt = supertypes.NewInt64Value(data.CreatedAt)
		return nil
	default:
		var diags diag.Diagnostics
		diags.AddError("Unknown type", fmt.Sprintf("Unknown type: %T", data))
		return diags
	}
}

func (r *ProjectServiceAccountResource) getNewParams(ctx context.Context, data ProjectServiceAccountResourceModel) (*openai.AdminOrganizationProjectServiceAccountNewParams, diag.Diagnostics) {
	return &openai.AdminOrganizationProjectServiceAccountNewParams{
		Name: data.Name.ValueString(),
	}, nil
}

func (r *ProjectServiceAccountResource) getUpdateParams(ctx context.Context, data ProjectServiceAccountResourceModel) (*openai.AdminOrganizationProjectServiceAccountUpdateParams, diag.Diagnostics) {
	return &openai.AdminOrganizationProjectServiceAccountUpdateParams{
		Name: openai.String(data.Name.ValueString()),
	}, nil
}
