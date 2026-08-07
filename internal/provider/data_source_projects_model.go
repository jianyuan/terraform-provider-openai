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
		item := ProjectsDataSourceModelProjectsItem{
			Id:        supertypes.NewStringValue(project.ID),
			Name:      supertypes.NewStringValue(project.Name),
			Status:    supertypes.NewStringValue(project.Status),
			CreatedAt: supertypes.NewInt64Value(project.CreatedAt),
		}

		if project.JSON.ExternalKeyID.Valid() {
			item.ExternalKeyId = supertypes.NewStringValue(project.ExternalKeyID)
		} else {
			item.ExternalKeyId = supertypes.NewStringNull()
		}

		if project.JSON.ArchivedAt.Valid() {
			item.ArchivedAt = supertypes.NewInt64Value(project.ArchivedAt)
		} else {
			item.ArchivedAt = supertypes.NewInt64Null()
		}

		return item
	}))

	return nil
}
