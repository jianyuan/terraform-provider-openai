package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *InviteResourceModel) Fill(ctx context.Context, invite openai.Invite) diag.Diagnostics {
	m.Id = supertypes.NewStringValue(invite.ID)
	m.Email = supertypes.NewStringValue(invite.Email)
	m.Role = supertypes.NewStringValue(string(invite.Role))
	m.Status = supertypes.NewStringValue(string(invite.Status))
	m.CreatedAt = supertypes.NewInt64Value(invite.CreatedAt)
	m.ExpiresAt = (func() supertypes.Int64Value {
		if invite.JSON.ExpiresAt.Valid() {
			return supertypes.NewInt64Value(invite.ExpiresAt)
		}
		return supertypes.NewInt64Null()
	})()
	m.AcceptedAt = (func() supertypes.Int64Value {
		if invite.JSON.AcceptedAt.Valid() {
			return supertypes.NewInt64Value(invite.AcceptedAt)
		}
		return supertypes.NewInt64Null()
	})()

	return nil
}

func (r *InviteResource) getCreateJSONRequestBody(ctx context.Context, data InviteResourceModel) (*openai.AdminOrganizationInviteNewParams, diag.Diagnostics) {
	return &openai.AdminOrganizationInviteNewParams{
		Email: data.Email.ValueString(),
		Role:  openai.AdminOrganizationInviteNewParamsRole(data.Role.ValueString()),
	}, nil
}
