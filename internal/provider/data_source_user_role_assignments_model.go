package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *UserRoleAssignmentsDataSourceModel) Fill(ctx context.Context, data []apiclient.AssignedRoleDetails) diag.Diagnostics {
	if data == nil {
		m.Roles = supertypes.NewSetNestedObjectValueOfNull[UserRoleAssignmentsDataSourceModelRolesItem](ctx)
	} else {
		items := make([]UserRoleAssignmentsDataSourceModelRolesItem, len(data))
		for i, role := range data {
			items[i] = UserRoleAssignmentsDataSourceModelRolesItem{
				Id:             supertypes.NewStringValue(role.Id),
				Name:           supertypes.NewStringValue(role.Name),
				Description:    supertypes.NewStringPointerValue(role.Description),
				Permissions:    supertypes.NewSetValueOfSlice(ctx, deduplicate(role.Permissions)), // For some reason, the API returns duplicate permissions
				PredefinedRole: supertypes.NewBoolValue(role.PredefinedRole),
				ResourceType:   supertypes.NewStringValue(role.ResourceType),
			}
		}
		m.Roles = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, items)
	}
	return nil
}
