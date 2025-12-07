package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *GroupUserResourceModel) Fill(ctx context.Context, data any) diag.Diagnostics {
	var diag diag.Diagnostics

	switch data := data.(type) {
	case apiclient.GroupUserAssignment:
		m.GroupId = supertypes.NewStringValue(data.GroupId)
		m.UserId = supertypes.NewStringValue(data.UserId)
		return nil
	case apiclient.User:
		m.UserId = supertypes.NewStringValue(data.Id)
		return nil
	default:
		diag.AddError("Unknown type", fmt.Sprintf("Unknown type: %T", data))
		return diag
	}
}

func (r *GroupUserResource) resourceMatch(data GroupUserResourceModel, user apiclient.User) bool {
	return data.UserId.ValueString() == user.Id
}

func (r *GroupUserResource) getCreateJSONRequestBody(ctx context.Context, data GroupUserResourceModel) (apiclient.AddGroupUserJSONRequestBody, diag.Diagnostics) {
	return apiclient.AddGroupUserJSONRequestBody{
		UserId: data.UserId.ValueString(),
	}, nil
}
