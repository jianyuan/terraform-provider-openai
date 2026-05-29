package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *UsersDataSourceModel) Fill(ctx context.Context, users []apiclient.User) diag.Diagnostics {
	items := make([]UsersDataSourceModelUsersItem, len(users))
	for i, user := range users {
		items[i] = UsersDataSourceModelUsersItem{
			Id:      supertypes.NewStringValue(user.Id),
			Email:   supertypes.NewStringPointerValue(user.Email),
			Name:    supertypes.NewStringPointerValue(user.Name),
			Role:    supertypes.NewStringPointerValue(user.Role),
			AddedAt: supertypes.NewInt64Value(user.AddedAt),
		}
	}
	m.Users = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, items)
	return nil
}
