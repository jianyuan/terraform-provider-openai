package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *OrganizationRoleResourceModel) Fill(ctx context.Context, role apiclient.Role) diag.Diagnostics {
	m.Id = supertypes.NewStringValue(role.Id)
	m.Name = supertypes.NewStringValue(role.Name)
	m.Description = supertypes.NewStringPointerValue(role.Description)
	m.Permissions = supertypes.NewSetValueOfSlice(ctx, deduplicate(role.Permissions))
	return nil
}

func (r *OrganizationRoleResource) resourceMatch(data OrganizationRoleResourceModel, role apiclient.Role) bool {
	return data.Id.ValueString() == role.Id
}

func (r *OrganizationRoleResource) getCreateJSONRequestBody(ctx context.Context, data OrganizationRoleResourceModel) (apiclient.CreateRoleJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics
	return apiclient.CreateRoleJSONRequestBody{
		RoleName:    data.Name.ValueString(),
		Permissions: mergeDiagnostics(data.Permissions.Get(ctx))(&diags),
		Description: data.Description.ValueStringPointer(),
	}, diags
}

func (r *OrganizationRoleResource) getUpdateJSONRequestBody(ctx context.Context, data OrganizationRoleResourceModel) (apiclient.UpdateRoleJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics
	return apiclient.UpdateRoleJSONRequestBody{
		RoleName:    data.Name.ValueStringPointer(),
		Permissions: ptr.Ptr(mergeDiagnostics(data.Permissions.Get(ctx))(&diags)),
		Description: data.Description.ValueStringPointer(),
	}, diags
}
