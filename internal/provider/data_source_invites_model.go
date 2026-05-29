package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *InvitesDataSourceModel) Fill(ctx context.Context, invites []apiclient.Invite) diag.Diagnostics {
	items := make([]InvitesDataSourceModelInvitesItem, len(invites))
	for i, invite := range invites {
		items[i] = InvitesDataSourceModelInvitesItem{
			Id:         supertypes.NewStringValue(invite.Id),
			Email:      supertypes.NewStringValue(invite.Email),
			Role:       supertypes.NewStringValue(string(invite.Role)),
			Status:     supertypes.NewStringValue(string(invite.Status)),
			CreatedAt:  supertypes.NewInt64Value(invite.CreatedAt),
			ExpiresAt:  supertypes.NewInt64PointerValue(invite.ExpiresAt),
			AcceptedAt: supertypes.NewInt64PointerValue(invite.AcceptedAt),
		}
	}
	m.Invites = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, items)
	return nil
}
