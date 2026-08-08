package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectUserResourceModel) Fill(ctx context.Context, data openai.ProjectUser) diag.Diagnostics {
	m.UserId = supertypes.NewStringValue(data.ID)
	m.Role = supertypes.NewStringValue(data.Role)
	return nil
}

func (r *ProjectUserResource) getNewParams(ctx context.Context, data ProjectUserResourceModel) (*openai.AdminOrganizationProjectUserNewParams, diag.Diagnostics) {
	return &openai.AdminOrganizationProjectUserNewParams{
		Role:   data.Role.ValueString(),
		UserID: openai.String(data.UserId.ValueString()),
	}, nil
}

func (r *ProjectUserResource) getUpdateParams(ctx context.Context, data ProjectUserResourceModel) (*openai.AdminOrganizationProjectUserUpdateParams, diag.Diagnostics) {
	return &openai.AdminOrganizationProjectUserUpdateParams{
		Role: openai.String(data.Role.ValueString()),
	}, nil
}
