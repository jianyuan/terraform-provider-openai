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

var _ datasource.DataSource = &OrganizationDataSource{}

func NewOrganizationDataSource() datasource.DataSource {
	return &OrganizationDataSource{}
}

type OrganizationDataSource struct {
	baseDataSource
}

type OrganizationDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	IsDefault   types.Bool   `tfsdk:"is_default"`
	Name        types.String `tfsdk:"name"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
}

func (m *OrganizationDataSourceModel) Fill(organization apiclient.Organization) error {
	m.Id = types.StringValue(organization.Id)
	m.IsDefault = types.BoolValue(organization.IsDefault)
	m.Name = types.StringValue(organization.Name)
	m.Title = types.StringValue(organization.Title)
	m.Description = types.StringValue(organization.Description)
	return nil
}

func (d *OrganizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (d *OrganizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve information about an organization.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the organization. If omitted, the default organization is used.",
				Optional:            true,
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
	}
}

func (d *OrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var org *apiclient.Organization

	if data.Id.IsNull() {
		// If the ID is not provided, use the default organization.
		httpResp, err := d.client.GetOrganizationsWithResponse(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
			return
		}

		if httpResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
			return
		}

		for _, organization := range *httpResp.JSON200.Data {
			if organization.IsDefault {
				org = &organization //nolint:exportloopref
				break
			}
		}
	} else {
		// If the ID is provided, use the specified organization.
		httpResp, err := d.client.GetOrganizationWithResponse(ctx, data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
			return
		}

		if httpResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
			return
		}

		org = httpResp.JSON200
	}

	if org == nil {
		resp.Diagnostics.AddError("API Error", "Organization not found")
		return
	}

	if err := data.Fill(*org); err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to unmarshal response: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
