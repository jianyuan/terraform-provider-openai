package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/samber/lo"
)

func (m *UsersDataSourceModel) Fill(ctx context.Context, users []openai.OrganizationUser) diag.Diagnostics {
	m.Users = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, lo.Map(users, func(user openai.OrganizationUser, _ int) UsersDataSourceModelUsersItem {
		return UsersDataSourceModelUsersItem{
			Id:      supertypes.NewStringValue(user.ID),
			Email:   supertypes.NewStringValue(user.Email),
			Name:    supertypes.NewStringValue(user.Name),
			Role:    supertypes.NewStringValue(user.Role),
			AddedAt: supertypes.NewInt64Value(user.AddedAt),
		}
	}))
	return nil
}
