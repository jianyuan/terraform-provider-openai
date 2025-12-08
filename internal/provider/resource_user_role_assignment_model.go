package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *UserRoleAssignmentResourceModel) Fill(ctx context.Context, data any) diag.Diagnostics {
	switch data := data.(type) {
	case apiclient.UserRoleAssignment:
		m.UserId = supertypes.NewStringValue(data.User.Id)
		m.RoleId = supertypes.NewStringValue(data.Role.Id)
		return nil
	case apiclient.AssignedRoleDetails:
		m.RoleId = supertypes.NewStringValue(data.Id)
		return nil
	default:
		var diags diag.Diagnostics
		diags.AddError("Unknown type", fmt.Sprintf("Unknown type: %T", data))
		return diags
	}
}

func (r *UserRoleAssignmentResource) resourceMatch(data UserRoleAssignmentResourceModel, roleAssignment apiclient.AssignedRoleDetails) bool {
	return data.RoleId.ValueString() == roleAssignment.Id
}

func (r *UserRoleAssignmentResource) getCreateJSONRequestBody(ctx context.Context, data UserRoleAssignmentResourceModel) (apiclient.AssignUserRoleJSONRequestBody, diag.Diagnostics) {
	return apiclient.AssignUserRoleJSONRequestBody{
		RoleId: data.RoleId.ValueString(),
	}, nil
}

func (r *UserRoleAssignmentResource) getUpdateJSONRequestBody(ctx context.Context, data UserRoleAssignmentResourceModel) (apiclient.AssignUserRoleJSONRequestBody, diag.Diagnostics) {
	return apiclient.AssignUserRoleJSONRequestBody{
		RoleId: data.RoleId.ValueString(),
	}, nil
}
