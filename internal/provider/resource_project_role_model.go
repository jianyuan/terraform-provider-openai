package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
)

func (r *ProjectRoleResource) getNewParams(ctx context.Context, data ProjectRoleResourceModel) (*openai.AdminOrganizationProjectRoleNewParams, diag.Diagnostics) {
	var diags diag.Diagnostics
	return &openai.AdminOrganizationProjectRoleNewParams{
		RoleName:    data.Name.ValueString(),
		Permissions: mergeDiagnostics(data.Permissions.Get(ctx))(&diags),
		Description: openai.String(data.Description.ValueString()),
	}, diags
}

func (r *ProjectRoleResource) getUpdateParams(ctx context.Context, data ProjectRoleResourceModel) (*openai.AdminOrganizationProjectRoleUpdateParams, diag.Diagnostics) {
	var diags diag.Diagnostics
	return &openai.AdminOrganizationProjectRoleUpdateParams{
		RoleName:    openai.String(data.Name.ValueString()),
		Permissions: mergeDiagnostics(data.Permissions.Get(ctx))(&diags),
		Description: openai.String(data.Description.ValueString()),
	}, diags
}
