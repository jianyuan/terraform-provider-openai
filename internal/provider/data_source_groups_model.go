package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *GroupsDataSourceModel) Fill(ctx context.Context, groups []apiclient.GroupResponse) diag.Diagnostics {
	items := make([]GroupsDataSourceModelGroupsItem, len(groups))
	for i, group := range groups {
		items[i] = GroupsDataSourceModelGroupsItem{
			Id:            supertypes.NewStringValue(group.Id),
			Name:          supertypes.NewStringValue(group.Name),
			IsScimManaged: supertypes.NewBoolValue(group.IsScimManaged),
			CreatedAt:     supertypes.NewInt64Value(group.CreatedAt),
		}
	}
	m.Groups = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, items)
	return nil
}
