package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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

func (m *ProjectModel) Fill(ctx context.Context, p apiclient.Project) (diags diag.Diagnostics) {
	m.Id = types.StringValue(p.Id)
	m.Name = types.StringValue(p.Name)
	m.Status = types.StringValue(string(p.Status))
	m.CreatedAt = types.Int64Value(p.CreatedAt)
	m.ArchivedAt = types.Int64PointerValue(p.ArchivedAt)
	return
}
