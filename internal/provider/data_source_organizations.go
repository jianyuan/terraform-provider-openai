package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

var _ datasource.DataSource = &OrganizationsDataSource{}

func NewOrganizationsDataSource() datasource.DataSource {
	return &OrganizationsDataSource{}
}

type OrganizationsDataSource struct {
	baseDataSource
}

type OrganizationsDataSourceModel struct {
	Organizations []OrganizationDataSourceModel `tfsdk:"organizations"`
}

func (m *OrganizationsDataSourceModel) Fill(organizations []apiclient.Organization) error {
	m.Organizations = make([]OrganizationDataSourceModel, len(organizations))
	for i, organization := range organizations {
		if err := m.Organizations[i].Fill(organization); err != nil {
			return err
		}
	}
	return nil
}

func (d *OrganizationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations"
}

func (d *OrganizationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all organizations.",

		Attributes: map[string]schema.Attribute{
			"organizations": schema.SetNestedAttribute{
				MarkdownDescription: "List of organizations.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Organization ID used in API requests.",
							Computed:            true,
						},
						"is_default": schema.BoolAttribute{
							MarkdownDescription: "Whether this organization is the default organization for the user.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Internal label for your organization.",
							Computed:            true,
						},
						"title": schema.StringAttribute{
							MarkdownDescription: "Human-friendly label for your organization, shown in user interfaces.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the organization.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *OrganizationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := d.client.GetOrganizationsWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	}

	if err := data.Fill(*httpResp.JSON200.Data); err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to unmarshal response: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
