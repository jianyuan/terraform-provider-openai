package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/samber/lo"
)

func (m *OrganizationRoleResourceModel) Fill(ctx context.Context, role openai.Role) diag.Diagnostics {
	m.Id = supertypes.NewStringValue(role.ID)
	m.Name = supertypes.NewStringValue(role.Name)
	m.Description = supertypes.NewStringValue(role.Description)
	m.Permissions = supertypes.NewSetValueOfSlice(ctx, lo.Uniq(role.Permissions))
	return nil
}

func (r *OrganizationRoleResource) getCreateJSONRequestBody(ctx context.Context, data OrganizationRoleResourceModel) (*openai.AdminOrganizationRoleNewParams, diag.Diagnostics) {
	var diags diag.Diagnostics
	return &openai.AdminOrganizationRoleNewParams{
		RoleName:    data.Name.ValueString(),
		Permissions: mergeDiagnostics(data.Permissions.Get(ctx))(&diags),
		Description: openai.String(data.Description.Get()),
	}, diags
}

func (r *OrganizationRoleResource) getUpdateJSONRequestBody(ctx context.Context, data OrganizationRoleResourceModel) (*openai.AdminOrganizationRoleUpdateParams, diag.Diagnostics) {
	var diags diag.Diagnostics
	return &openai.AdminOrganizationRoleUpdateParams{
		RoleName:    openai.String(data.Name.ValueString()),
		Permissions: mergeDiagnostics(data.Permissions.Get(ctx))(&diags),
		Description: openai.String(data.Description.Get()),
	}, diags
}
