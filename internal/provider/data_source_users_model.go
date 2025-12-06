package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *UsersDataSourceModel) Fill(ctx context.Context, users []apiclient.User) diag.Diagnostics {
	if users == nil {
		m.Users = supertypes.NewSetNestedObjectValueOfNull[UsersDataSourceModelUsersItem](ctx)
	} else {
		items := make([]UsersDataSourceModelUsersItem, len(users))
		for i, user := range users {
			items[i] = UsersDataSourceModelUsersItem{
				Id:      supertypes.NewStringValue(user.Id),
				Email:   supertypes.NewStringValue(user.Email),
				Name:    supertypes.NewStringValue(user.Name),
				Role:    supertypes.NewStringValue(string(user.Role)),
				AddedAt: supertypes.NewInt64Value(user.AddedAt),
			}
		}
		m.Users = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, items)
	}
	return nil
}
