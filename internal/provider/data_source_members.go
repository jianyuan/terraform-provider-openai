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

var _ datasource.DataSource = &MembersDataSource{}

func NewMembersDataSource() datasource.DataSource {
	return &MembersDataSource{}
}

type MembersDataSource struct {
	baseDataSource
}

type InvitedMembersDataSourceModel struct {
	Id        types.String `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	IsExpired types.Bool   `tfsdk:"is_expired"`
	Role      types.String `tfsdk:"role"`
}

func (m *InvitedMembersDataSourceModel) Fill(u apiclient.InvitedUser) error {
	m.Id = types.StringValue(u.Id)
	m.Email = types.StringValue(u.Email)
	m.IsExpired = types.BoolValue(u.IsExpired)
	m.Role = types.StringValue(string(u.Role))

	return nil
}

type OrganizationUserDataSourceModel struct {
	Id               types.String `tfsdk:"id"`
	Email            types.String `tfsdk:"email"`
	Name             types.String `tfsdk:"name"`
	Picture          types.String `tfsdk:"picture"`
	IsDefault        types.Bool   `tfsdk:"is_default"`
	IsServiceAccount types.Bool   `tfsdk:"is_service_account"`
	Role             types.String `tfsdk:"role"`
}

func (m *OrganizationUserDataSourceModel) Fill(u apiclient.OrganizationUser) error {
	m.Id = types.StringValue(u.User.Id)
	m.Email = types.StringValue(u.User.Email)
	m.Name = types.StringValue(u.User.Name)
	m.Picture = types.StringPointerValue(u.User.Picture)
	m.IsDefault = types.BoolValue(u.IsDefault)
	m.IsServiceAccount = types.BoolValue(u.IsServiceAccount)
	m.Role = types.StringValue(string(u.Role))

	return nil
}

type MembersDataSourceModel struct {
	InvitedMembers []InvitedMembersDataSourceModel   `tfsdk:"invited_members"`
	Members        []OrganizationUserDataSourceModel `tfsdk:"members"`
}

func (m *MembersDataSourceModel) Fill(invitedUsers []apiclient.InvitedUser, members []apiclient.OrganizationUser) error {
	m.InvitedMembers = make([]InvitedMembersDataSourceModel, len(invitedUsers))
	for i, u := range invitedUsers {
		if err := m.InvitedMembers[i].Fill(u); err != nil {
			return err
		}
	}

	m.Members = make([]OrganizationUserDataSourceModel, len(members))
	for i, u := range members {
		if err := m.Members[i].Fill(u); err != nil {
			return err
		}
	}

	return nil
}

func (d *MembersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_members"
}

func (d *MembersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all users in an organization, including invited users and members.",

		Attributes: map[string]schema.Attribute{
			"invited_members": schema.SetNestedAttribute{
				MarkdownDescription: "List of invited users.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Invite identifier.",
							Computed:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "Email address of the invited user.",
							Computed:            true,
						},
						"is_expired": schema.BoolAttribute{
							MarkdownDescription: "Whether the invite has expired.",
							Computed:            true,
						},
						"role": schema.StringAttribute{
							MarkdownDescription: "Role of the invited user.",
							Computed:            true,
						},
					},
				},
			},
			"members": schema.SetNestedAttribute{
				MarkdownDescription: "List of members.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "User identifier.",
							Computed:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "Email address of the user.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the user.",
							Computed:            true,
						},
						"picture": schema.StringAttribute{
							MarkdownDescription: "URL of the user's profile picture.",
							Computed:            true,
						},
						"is_default": schema.BoolAttribute{
							MarkdownDescription: "Whether this user is the default user for the organization.",
							Computed:            true,
						},
						"is_service_account": schema.BoolAttribute{
							MarkdownDescription: "Whether this user is a service account.",
							Computed:            true,
						},
						"role": schema.StringAttribute{
							MarkdownDescription: "Role of the user.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *MembersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MembersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := d.client.GetOrganizationUsersWithResponse(
		ctx,
		"TODO",
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	}

	if err := data.Fill(httpResp.JSON200.Invited, httpResp.JSON200.Members.Data); err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to unmarshal response: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
