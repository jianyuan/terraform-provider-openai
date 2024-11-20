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

func (m *ProjectModel) Fill(p apiclient.Project) error {
	m.Id = types.StringValue(p.Id)
	m.Name = types.StringValue(p.Name)
	m.Status = types.StringValue(string(p.Status))
	m.CreatedAt = types.Int64Value(int64(p.CreatedAt))
	if p.ArchivedAt == nil {
		m.ArchivedAt = types.Int64Null()
	} else {
		m.ArchivedAt = types.Int64Value(int64(*p.ArchivedAt))

	}
	return nil
}
