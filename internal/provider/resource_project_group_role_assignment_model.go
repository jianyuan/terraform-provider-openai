package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectGroupRoleAssignmentResourceModel) Fill(ctx context.Context, data any) diag.Diagnostics {
	switch data := data.(type) {
	case openai.AdminOrganizationProjectGroupRoleNewResponse:
		m.GroupId = supertypes.NewStringValue(data.Group.ID)
		m.RoleId = supertypes.NewStringValue(data.Role.ID)
		return nil
	case openai.AdminOrganizationProjectGroupRoleGetResponse:
		m.RoleId = supertypes.NewStringValue(data.ID)
		return nil
	default:
		var diags diag.Diagnostics
		diags.AddError("Unknown type", fmt.Sprintf("Unknown type: %T", data))
		return diags
	}
}

func (r *ProjectGroupRoleAssignmentResource) getCreateJSONRequestBody(ctx context.Context, data ProjectGroupRoleAssignmentResourceModel) (*openai.AdminOrganizationProjectGroupRoleNewParams, diag.Diagnostics) {
	return &openai.AdminOrganizationProjectGroupRoleNewParams{
		RoleID: data.RoleId.ValueString(),
	}, nil
}
