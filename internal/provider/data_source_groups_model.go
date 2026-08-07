package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/samber/lo"
)

func (m *GroupsDataSourceModel) Fill(ctx context.Context, groups []openai.Group) diag.Diagnostics {
	m.Groups = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, lo.Map(groups, func(group openai.Group, _ int) GroupsDataSourceModelGroupsItem {
		return GroupsDataSourceModelGroupsItem{
			Id:            supertypes.NewStringValue(group.ID),
			Name:          supertypes.NewStringValue(group.Name),
			IsScimManaged: supertypes.NewBoolValue(group.IsScimManaged),
			CreatedAt:     supertypes.NewInt64Value(group.CreatedAt),
		}
	}))

	return nil
}
