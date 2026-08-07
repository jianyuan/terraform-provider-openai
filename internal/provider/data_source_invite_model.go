package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *InviteDataSourceModel) Fill(ctx context.Context, invite openai.Invite) diag.Diagnostics {
	m.Id = supertypes.NewStringValue(invite.ID)
	m.Email = supertypes.NewStringValue(invite.Email)
	m.Role = supertypes.NewStringValue(string(invite.Role))
	m.Status = supertypes.NewStringValue(string(invite.Status))
	m.CreatedAt = supertypes.NewInt64Value(invite.CreatedAt)

	if invite.JSON.ExpiresAt.Valid() {
		m.ExpiresAt = supertypes.NewInt64Value(invite.ExpiresAt)
	} else {
		m.ExpiresAt = supertypes.NewInt64Null()
	}

	if invite.JSON.AcceptedAt.Valid() {
		m.AcceptedAt = supertypes.NewInt64Value(invite.AcceptedAt)
	} else {
		m.AcceptedAt = supertypes.NewInt64Null()
	}

	return nil
}
