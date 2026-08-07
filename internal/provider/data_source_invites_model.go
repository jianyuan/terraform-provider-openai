package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/samber/lo"
)

func (m *InvitesDataSourceModel) Fill(ctx context.Context, invites []openai.Invite) diag.Diagnostics {
	m.Invites = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, lo.Map(invites, func(invite openai.Invite, _ int) InvitesDataSourceModelInvitesItem {
		item := InvitesDataSourceModelInvitesItem{
			Id:        supertypes.NewStringValue(invite.ID),
			Email:     supertypes.NewStringValue(invite.Email),
			Role:      supertypes.NewStringValue(string(invite.Role)),
			Status:    supertypes.NewStringValue(string(invite.Status)),
			CreatedAt: supertypes.NewInt64Value(invite.CreatedAt),
		}

		if invite.JSON.ExpiresAt.Valid() {
			item.ExpiresAt = supertypes.NewInt64Value(invite.ExpiresAt)
		} else {
			item.ExpiresAt = supertypes.NewInt64Null()
		}

		if invite.JSON.AcceptedAt.Valid() {
			item.AcceptedAt = supertypes.NewInt64Value(invite.AcceptedAt)
		} else {
			item.AcceptedAt = supertypes.NewInt64Null()
		}

		return item
	}))

	return nil
}
