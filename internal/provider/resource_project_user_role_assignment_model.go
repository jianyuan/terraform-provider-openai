package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectUserRoleAssignmentResourceModel) Fill(ctx context.Context, data any) diag.Diagnostics {
	switch data := data.(type) {
	case openai.AdminOrganizationProjectUserRoleNewResponse:
		m.UserId = supertypes.NewStringValue(data.User.ID)
		m.RoleId = supertypes.NewStringValue(data.Role.ID)
		return nil
	case openai.AdminOrganizationProjectUserRoleGetResponse:
		m.RoleId = supertypes.NewStringValue(data.ID)
		return nil
	default:
		var diags diag.Diagnostics
		diags.AddError("Unknown type", fmt.Sprintf("Unknown type: %T", data))
		return diags
	}
}

func (r *ProjectUserRoleAssignmentResource) getCreateJSONRequestBody(ctx context.Context, data ProjectUserRoleAssignmentResourceModel) (*openai.AdminOrganizationProjectUserRoleNewParams, diag.Diagnostics) {
	return &openai.AdminOrganizationProjectUserRoleNewParams{
		RoleID: data.RoleId.ValueString(),
	}, nil
}
