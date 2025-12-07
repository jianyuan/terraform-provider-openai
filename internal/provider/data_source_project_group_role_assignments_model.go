package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectGroupRoleAssignmentsDataSourceModel) Fill(ctx context.Context, data []apiclient.AssignedRoleDetails) diag.Diagnostics {
	if data == nil {
		m.Roles = supertypes.NewSetNestedObjectValueOfNull[ProjectGroupRoleAssignmentsDataSourceModelRolesItem](ctx)
	} else {
		items := make([]ProjectGroupRoleAssignmentsDataSourceModelRolesItem, len(data))
		for i, role := range data {
			items[i] = ProjectGroupRoleAssignmentsDataSourceModelRolesItem{
				Id:             supertypes.NewStringValue(role.Id),
				Name:           supertypes.NewStringValue(role.Name),
				Description:    supertypes.NewStringPointerValue(role.Description),
				Permissions:    supertypes.NewSetValueOfSlice(ctx, deduplicate(role.Permissions)),
				PredefinedRole: supertypes.NewBoolValue(role.PredefinedRole),
				ResourceType:   supertypes.NewStringValue(role.ResourceType),
			}
		}
		m.Roles = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, items)
	}
	return nil
}
