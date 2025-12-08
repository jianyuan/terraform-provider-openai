package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectsDataSourceModel) Fill(ctx context.Context, projects []apiclient.Project) diag.Diagnostics {
	items := make([]ProjectsDataSourceModelProjectsItem, len(projects))
	for i, project := range projects {
		items[i] = ProjectsDataSourceModelProjectsItem{
			Id:            supertypes.NewStringValue(project.Id),
			Name:          supertypes.NewStringValue(project.Name),
			Status:        supertypes.NewStringValue(string(project.Status)),
			ExternalKeyId: supertypes.NewStringPointerValue(project.ExternalKeyId),
			CreatedAt:     supertypes.NewInt64Value(project.CreatedAt),
			ArchivedAt:    supertypes.NewInt64PointerValue(project.ArchivedAt),
		}
	}
	m.Projects = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, items)
	return nil
}
