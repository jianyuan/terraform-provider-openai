package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/samber/lo"
)

func (m *ProjectsDataSourceModel) Fill(ctx context.Context, projects []openai.Project) diag.Diagnostics {
	m.Projects = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, lo.Map(projects, func(project openai.Project, _ int) ProjectsDataSourceModelProjectsItem {
		return ProjectsDataSourceModelProjectsItem{
			Id:        supertypes.NewStringValue(project.ID),
			Name:      supertypes.NewStringValue(project.Name),
			Status:    supertypes.NewStringValue(project.Status),
			CreatedAt: supertypes.NewInt64Value(project.CreatedAt),
			ExternalKeyId: (func() supertypes.StringValue {
				if project.JSON.ExternalKeyID.Valid() {
					return supertypes.NewStringValue(project.ExternalKeyID)
				}
				return supertypes.NewStringNull()
			})(),
			ArchivedAt: (func() supertypes.Int64Value {
				if project.JSON.ArchivedAt.Valid() {
					return supertypes.NewInt64Value(project.ArchivedAt)
				}
				return supertypes.NewInt64Null()
			})(),
		}
	}))
	return nil
}
