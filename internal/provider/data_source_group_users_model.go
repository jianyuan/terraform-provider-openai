package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/samber/lo"
)

func (m *GroupUsersDataSourceModel) Fill(ctx context.Context, users []openai.OrganizationGroupUser) diag.Diagnostics {
	m.Users = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, lo.Map(users, func(user openai.OrganizationGroupUser, _ int) GroupUsersDataSourceModelUsersItem {
		return GroupUsersDataSourceModelUsersItem{
			Id:    supertypes.NewStringValue(user.ID),
			Email: supertypes.NewStringValue(user.Email),
			Name:  supertypes.NewStringValue(user.Name),
		}
	}))
	return nil
}
