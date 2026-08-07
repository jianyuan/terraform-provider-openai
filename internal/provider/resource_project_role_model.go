package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/samber/lo"
)

func (m *ProjectRoleResourceModel) Fill(ctx context.Context, role openai.Role) diag.Diagnostics {
	m.Id = supertypes.NewStringValue(role.ID)
	m.Name = supertypes.NewStringValue(role.Name)
	m.Description = supertypes.NewStringValue(role.Description)
	m.Permissions = supertypes.NewSetValueOfSlice(ctx, lo.Uniq(role.Permissions))
	return nil
}

func (r *ProjectRoleResource) getCreateJSONRequestBody(ctx context.Context, data ProjectRoleResourceModel) (*openai.AdminOrganizationProjectRoleNewParams, diag.Diagnostics) {
	var diags diag.Diagnostics
	return &openai.AdminOrganizationProjectRoleNewParams{
		RoleName:    data.Name.ValueString(),
		Permissions: mergeDiagnostics(data.Permissions.Get(ctx))(&diags),
		Description: openai.String(data.Description.ValueString()),
	}, diags
}

func (r *ProjectRoleResource) getUpdateJSONRequestBody(ctx context.Context, data ProjectRoleResourceModel) (*openai.AdminOrganizationProjectRoleUpdateParams, diag.Diagnostics) {
	var diags diag.Diagnostics
	return &openai.AdminOrganizationProjectRoleUpdateParams{
		RoleName:    openai.String(data.Name.ValueString()),
		Permissions: mergeDiagnostics(data.Permissions.Get(ctx))(&diags),
		Description: openai.String(data.Description.ValueString()),
	}, diags
}
