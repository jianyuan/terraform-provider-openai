package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

var _ datasource.DataSource = &ProjectsDataSource{}

func NewProjectsDataSource() datasource.DataSource {
	return &ProjectsDataSource{}
}

type ProjectsDataSource struct {
	baseDataSource
}

type ProjectsDataSourceModel_Project struct {
	Id    types.String `tfsdk:"id"`
	Title types.String `tfsdk:"title"`
}

func (m *ProjectsDataSourceModel_Project) Fill(project apiclient.Project) error {
	m.Id = types.StringValue(project.Id)
	m.Title = types.StringValue(project.Title)
	return nil
}

type ProjectsDataSourceModel struct {
	OrganizationId types.String                      `tfsdk:"organization_id"`
	Projects       []ProjectsDataSourceModel_Project `tfsdk:"projects"`
}

func (m *ProjectsDataSourceModel) Fill(projects []apiclient.Project) error {
	m.Projects = make([]ProjectsDataSourceModel_Project, len(projects))
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
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the organization.",
				Required:            true,
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
						"title": schema.StringAttribute{
							MarkdownDescription: "Human-friendly label for the project, shown in user interfaces.",
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

	httpResp, err := d.client.GetOrganizationProjectsWithResponse(
		ctx,
		data.OrganizationId.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	}

	if err := data.Fill(httpResp.JSON200.Data); err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to unmarshal response: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
