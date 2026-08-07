package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/samber/lo"
)

func (m *OrganizationRolesDataSourceModel) Fill(ctx context.Context, roles []openai.Role) diag.Diagnostics {
	m.Roles = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, lo.Map(roles, func(role openai.Role, _ int) OrganizationRolesDataSourceModelRolesItem {
		return OrganizationRolesDataSourceModelRolesItem{
			Id:             supertypes.NewStringValue(role.ID),
			Name:           supertypes.NewStringValue(role.Name),
			Description:    supertypes.NewStringValue(role.Description),
			Permissions:    supertypes.NewSetValueOfSlice(ctx, lo.Uniq(role.Permissions)), // For some reason, the API returns duplicate permissions
			PredefinedRole: supertypes.NewBoolValue(role.PredefinedRole),
			ResourceType:   supertypes.NewStringValue(role.ResourceType),
		}
	}))

	return nil
}
