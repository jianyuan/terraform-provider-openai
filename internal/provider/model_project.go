package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

type ProjectModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Status     types.String `tfsdk:"status"`
	CreatedAt  types.Int64  `tfsdk:"created_at"`
	ArchivedAt types.Int64  `tfsdk:"archived_at"`
}

func (m *ProjectModel) Fill(project apiclient.Project) error {
	m.Id = types.StringValue(project.Id)
	m.Name = types.StringValue(project.Name)
	m.Status = types.StringValue(string(project.Status))
	m.CreatedAt = types.Int64Value(int64(project.CreatedAt))
	if project.ArchivedAt == nil {
		m.ArchivedAt = types.Int64Null()
	} else {
		m.ArchivedAt = types.Int64Value(int64(*project.ArchivedAt))

	}
	return nil
}
