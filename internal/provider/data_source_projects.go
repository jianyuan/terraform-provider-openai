package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	"github.com/jianyuan/terraform-provider-openai/internal/ptr"
)

var _ datasource.DataSource = &ProjectsDataSource{}

func NewProjectsDataSource() datasource.DataSource {
	return &ProjectsDataSource{}
}

type ProjectsDataSource struct {
	baseDataSource
}

type ProjectsDataSourceModel struct {
	IncludeArchived types.Bool     `tfsdk:"include_archived"`
	Projects        []ProjectModel `tfsdk:"projects"`
}

func (m *ProjectsDataSourceModel) Fill(projects []apiclient.Project) error {
	m.Projects = make([]ProjectModel, len(projects))
	for i, project := range projects {
		if err := m.Projects[i].Fill(project); err != nil {
			return err
		}
	}
	return nil
}

func (d *ProjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *ProjectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all projects in an organization.",

		Attributes: map[string]schema.Attribute{
			"include_archived": schema.BoolAttribute{
				MarkdownDescription: "Include archived projects. Default is `false`.",
				Optional:            true,
			},
			"projects": schema.SetNestedAttribute{
				MarkdownDescription: "List of projects.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Project ID.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the project. This appears in reporting.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "Status `active` or `archived`.",
							Computed:            true,
						},
						"created_at": schema.Int64Attribute{
							MarkdownDescription: "The Unix timestamp (in seconds) of when the project was created.",
							Computed:            true,
						},
						"archived_at": schema.Int64Attribute{
							MarkdownDescription: "The Unix timestamp (in seconds) of when the project was archived or `null`.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *ProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var projects []apiclient.Project
	params := &apiclient.ListProjectsParams{
		Limit:           ptr.Ptr(100),
		IncludeArchived: data.IncludeArchived.ValueBoolPointer(),
	}

	for {
		httpResp, err := d.client.ListProjectsWithResponse(
			ctx,
			params,
		)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
			return
		}

		if httpResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
			return
		}

		projects = append(projects, httpResp.JSON200.Data...)

		if !httpResp.JSON200.HasMore {
			break
		}

		params.After = &httpResp.JSON200.LastId
	}

	if err := data.Fill(projects); err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to unmarshal response: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
