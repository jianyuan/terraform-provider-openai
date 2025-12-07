package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *GroupRoleResourceModel) Fill(ctx context.Context, data any) diag.Diagnostics {
	switch data := data.(type) {
	case apiclient.GroupRoleAssignment:
		m.GroupId = supertypes.NewStringValue(data.Group.Id)
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

func (r *GroupRoleResource) resourceMatch(data GroupRoleResourceModel, roleAssignment apiclient.AssignedRoleDetails) bool {
	return data.RoleId.ValueString() == roleAssignment.Id
}

func (r *GroupRoleResource) getCreateJSONRequestBody(ctx context.Context, data GroupRoleResourceModel) (apiclient.AssignGroupRoleJSONRequestBody, diag.Diagnostics) {
	return apiclient.AssignGroupRoleJSONRequestBody{
		RoleId: data.RoleId.ValueString(),
	}, nil
}
