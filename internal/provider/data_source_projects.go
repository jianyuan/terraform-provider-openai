package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
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
	Limit           types.Int64    `tfsdk:"limit"`
	Projects        []ProjectModel `tfsdk:"projects"`
}

func (m *ProjectsDataSourceModel) Fill(ctx context.Context, projects []apiclient.Project) (diags diag.Diagnostics) {
	m.Projects = make([]ProjectModel, len(projects))
	for i, project := range projects {
		diags.Append(m.Projects[i].Fill(ctx, project)...)
		if diags.HasError() {
			return
		}
	}
	return
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
			"limit": schema.Int64Attribute{
				MarkdownDescription: "Limit the number of projects to return. Default is to return all projects.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
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
						"external_key_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the customer-managed encryption key used for Enterprise Key Management (EKM). EKM is only available on certain accounts. Refer to the [EKM (External Keys) in the Management API Article](https://help.openai.com/en/articles/20000953-ekm-external-keys-in-the-management-api).",
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
		IncludeArchived: data.IncludeArchived.ValueBoolPointer(),
	}

	// Set the limit for the API request
	if data.Limit.IsNull() {
		params.Limit = ptr.Ptr(int64(100))
	} else {
		requestLimit := data.Limit.ValueInt64()
		if requestLimit > 100 {
			params.Limit = ptr.Ptr(int64(100))
		} else {
			params.Limit = ptr.Ptr(requestLimit)
		}
	}

	for {
		// Recalculate the limit for each request to ensure we don't exceed the desired limit
		if !data.Limit.IsNull() {
			remainingLimit := data.Limit.ValueInt64() - int64(len(projects))
			if remainingLimit <= 0 {
				break
			}
			if remainingLimit > 100 {
				params.Limit = ptr.Ptr(int64(100))
			} else {
				params.Limit = ptr.Ptr(remainingLimit)
			}
		}

		httpResp, err := d.client.ListProjectsWithResponse(
			ctx,
			params,
		)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
			return
		} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
			return
		}

		projects = append(projects, httpResp.JSON200.Data...)

		// If limit is set and we have enough projects, break.
		if !data.Limit.IsNull() && len(projects) >= int(data.Limit.ValueInt64()) {
			projects = projects[:data.Limit.ValueInt64()]
			break
		}

		// If there are no more projects, break.
		if !httpResp.JSON200.HasMore {
			break
		}

		params.After = &httpResp.JSON200.LastId
	}

	resp.Diagnostics.Append(data.Fill(ctx, projects)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
