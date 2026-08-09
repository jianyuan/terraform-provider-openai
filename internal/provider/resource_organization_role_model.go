package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
)

func (r *OrganizationRoleResource) getNewParams(ctx context.Context, data OrganizationRoleResourceModel) (*openai.AdminOrganizationRoleNewParams, diag.Diagnostics) {
	var diags diag.Diagnostics
	return &openai.AdminOrganizationRoleNewParams{
		RoleName:    data.Name.ValueString(),
		Permissions: mergeDiagnostics(data.Permissions.Get(ctx))(&diags),
		Description: openai.String(data.Description.ValueString()),
	}, diags
}

func (r *OrganizationRoleResource) getUpdateParams(ctx context.Context, data OrganizationRoleResourceModel) (*openai.AdminOrganizationRoleUpdateParams, diag.Diagnostics) {
	var diags diag.Diagnostics
	return &openai.AdminOrganizationRoleUpdateParams{
		RoleName:    openai.String(data.Name.ValueString()),
		Permissions: mergeDiagnostics(data.Permissions.Get(ctx))(&diags),
		Description: openai.String(data.Description.ValueString()),
	}, diags
}
