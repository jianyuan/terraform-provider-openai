package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *InviteResourceModel) Fill(ctx context.Context, invite apiclient.Invite) diag.Diagnostics {
	m.Id = supertypes.NewStringValue(invite.Id)
	m.Email = supertypes.NewStringValue(invite.Email)
	m.Role = supertypes.NewStringValue(string(invite.Role))
	m.Status = supertypes.NewStringValue(string(invite.Status))
	m.InvitedAt = supertypes.NewInt64Value(invite.InvitedAt)
	m.ExpiresAt = supertypes.NewInt64Value(invite.ExpiresAt)
	m.AcceptedAt = supertypes.NewInt64PointerValue(invite.AcceptedAt)
	return nil
}

func (r *InviteResource) getCreateJSONRequestBody(data InviteResourceModel) apiclient.InviteUserJSONRequestBody {
	return apiclient.InviteUserJSONRequestBody{
		Email: data.Email.ValueString(),
		Role:  apiclient.InviteRequestRole(data.Role.ValueString()),
	}
}
