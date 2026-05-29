package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *GroupUsersDataSourceModel) Fill(ctx context.Context, users []apiclient.GroupUser) diag.Diagnostics {
	items := make([]GroupUsersDataSourceModelUsersItem, len(users))
	for i, user := range users {
		items[i] = GroupUsersDataSourceModelUsersItem{
			Id:    supertypes.NewStringValue(user.Id),
			Email: supertypes.NewStringPointerValue(user.Email),
			Name:  supertypes.NewStringValue(user.Name),
		}
	}
	m.Users = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, items)
	return nil
}
