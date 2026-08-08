package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *UserRoleAssignmentResourceModel) Fill(ctx context.Context, data any) diag.Diagnostics {
	switch data := data.(type) {
	case openai.AdminOrganizationUserRoleNewResponse:
		m.UserId = supertypes.NewStringValue(data.User.ID)
		m.RoleId = supertypes.NewStringValue(data.Role.ID)
		return nil
	case openai.AdminOrganizationUserRoleGetResponse:
		m.RoleId = supertypes.NewStringValue(data.ID)
		return nil
	default:
		var diags diag.Diagnostics
		diags.AddError("Unknown type", fmt.Sprintf("Unknown type: %T", data))
		return diags
	}
}

func (r *UserRoleAssignmentResource) getNewParams(ctx context.Context, data UserRoleAssignmentResourceModel) (*openai.AdminOrganizationUserRoleNewParams, diag.Diagnostics) {
	return &openai.AdminOrganizationUserRoleNewParams{
		RoleID: data.RoleId.ValueString(),
	}, nil
}

func (r *UserRoleAssignmentResource) getUpdateParams(ctx context.Context, data UserRoleAssignmentResourceModel) (*openai.AdminOrganizationUserRoleNewParams, diag.Diagnostics) {
	return r.getNewParams(ctx, data)
}
