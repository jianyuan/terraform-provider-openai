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

type UserDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Email   types.String `tfsdk:"email"`
	Name    types.String `tfsdk:"name"`
	Role    types.String `tfsdk:"role"`
	AddedAt types.Int64  `tfsdk:"added_at"`
}

func (m *UserDataSourceModel) Fill(u apiclient.User) error {
	m.Id = types.StringValue(u.Id)
	m.Email = types.StringValue(u.Email)
	m.Name = types.StringValue(u.Name)
	m.Role = types.StringValue(string(u.Role))
	m.AddedAt = types.Int64Value(int64(u.AddedAt))
	return nil
}

type UsersDataSourceModel struct {
	Users []UserDataSourceModel `tfsdk:"users"`
}

func (m *UsersDataSourceModel) Fill(users []apiclient.User) error {
	m.Users = make([]UserDataSourceModel, len(users))
	for i, u := range users {
		if err := m.Users[i].Fill(u); err != nil {
			return err
		}
	}
	return nil
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
							MarkdownDescription: "User identifier.",
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
	params := &apiclient.ListUsersParams{}

	for {
		httpResp, err := d.client.ListUsersWithResponse(
			ctx,
			params,
		)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
			return
		}

		if httpResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
			return
		}

		users = append(users, httpResp.JSON200.Data...)

		if !httpResp.JSON200.HasMore {
			break
		}

		params.After = &httpResp.JSON200.LastId
	}

	if err := data.Fill(users); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to fill data: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
