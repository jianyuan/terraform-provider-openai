package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *GroupUsersDataSourceModel) Fill(ctx context.Context, users []apiclient.User) diag.Diagnostics {
	items := make([]GroupUsersDataSourceModelUsersItem, len(users))
	for i, user := range users {
		items[i] = GroupUsersDataSourceModelUsersItem{
			Id:      supertypes.NewStringValue(user.Id),
			Email:   supertypes.NewStringValue(user.Email),
			Name:    supertypes.NewStringValue(user.Name),
			Role:    supertypes.NewStringValue(string(user.Role)),
			AddedAt: supertypes.NewInt64Value(user.AddedAt),
		}
	}
	m.Users = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, items)
	return nil
}
