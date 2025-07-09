package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

type UsersDataSourceModel struct {
	Users []UserModel `tfsdk:"users"`
}

func (m *UsersDataSourceModel) Fill(ctx context.Context, users []apiclient.User) (diags diag.Diagnostics) {
	m.Users = make([]UserModel, len(users))
	for i, u := range users {
		diags.Append(m.Users[i].Fill(ctx, u)...)
		if diags.HasError() {
			return
		}
	}
	return
}

var _ datasource.DataSource = &UsersDataSource{}

func NewUsersDataSource() datasource.DataSource {
	return &UsersDataSource{}
}

type UsersDataSource struct {
	baseDataSource
}

func (d *UsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *UsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists all of the users in the organization.",

		Attributes: map[string]schema.Attribute{
			"users": schema.SetNestedAttribute{
				MarkdownDescription: "List of users.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "User ID.",
							Computed:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "The email address of the user.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the user.",
							Computed:            true,
						},
						"role": schema.StringAttribute{
							MarkdownDescription: "Role `owner` or `reader`.",
							Computed:            true,
						},
						"added_at": schema.Int64Attribute{
							MarkdownDescription: "The Unix timestamp (in seconds) of when the user was added.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *UsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UsersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var users []apiclient.User
	params := &apiclient.ListUsersParams{
		Limit: ptr.Ptr(int64(100)),
	}

	for {
		httpResp, err := d.client.ListUsersWithResponse(
			ctx,
			params,
		)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
			return
		} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
			return
		}

		users = append(users, httpResp.JSON200.Data...)

		if !httpResp.JSON200.HasMore {
			break
		}

		params.After = &httpResp.JSON200.LastId
	}

	resp.Diagnostics.Append(data.Fill(ctx, users)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
