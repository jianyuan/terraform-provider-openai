package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *GroupUserResourceModel) Fill(ctx context.Context, data any) diag.Diagnostics {
	switch data := data.(type) {
	case openai.AdminOrganizationGroupUserNewResponse:
		m.GroupId = supertypes.NewStringValue(data.GroupID)
		m.UserId = supertypes.NewStringValue(data.UserID)
		return nil
	case openai.AdminOrganizationGroupUserGetResponse:
		m.UserId = supertypes.NewStringValue(data.ID)
		return nil
	default:
		var diags diag.Diagnostics
		diags.AddError("Unknown type", fmt.Sprintf("Unknown type: %T", data))
		return diags
	}
}

func (r *GroupUserResource) getNewParams(ctx context.Context, data GroupUserResourceModel) (*openai.AdminOrganizationGroupUserNewParams, diag.Diagnostics) {
	return &openai.AdminOrganizationGroupUserNewParams{
		UserID: data.UserId.ValueString(),
	}, nil
}
