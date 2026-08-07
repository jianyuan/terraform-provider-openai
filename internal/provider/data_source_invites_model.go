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
		return InvitesDataSourceModelInvitesItem{
			Id:        supertypes.NewStringValue(invite.ID),
			Email:     supertypes.NewStringValue(invite.Email),
			Role:      supertypes.NewStringValue(string(invite.Role)),
			Status:    supertypes.NewStringValue(string(invite.Status)),
			CreatedAt: supertypes.NewInt64Value(invite.CreatedAt),
			ExpiresAt: (func() supertypes.Int64Value {
				if invite.JSON.ExpiresAt.Valid() {
					return supertypes.NewInt64Value(invite.ExpiresAt)
				}
				return supertypes.NewInt64Null()
			})(),
			AcceptedAt: (func() supertypes.Int64Value {
				if invite.JSON.AcceptedAt.Valid() {
					return supertypes.NewInt64Value(invite.AcceptedAt)
				}
				return supertypes.NewInt64Null()
			})(),
		}
	}))

	return nil
}
