package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

type InvitesDataSourceModel struct {
	Invites []InviteModel `tfsdk:"invites"`
}

func (m *InvitesDataSourceModel) Fill(invites []apiclient.Invite) error {
	m.Invites = make([]InviteModel, len(invites))
	for i, invite := range invites {
		if err := m.Invites[i].Fill(invite); err != nil {
			return err
		}
	}
	return nil
}

var _ datasource.DataSource = &InvitesDataSource{}

func NewInvitesDataSource() datasource.DataSource {
	return &InvitesDataSource{}
}

type InvitesDataSource struct {
	baseDataSource
}

func (d *InvitesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_invites"
}

func (d *InvitesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists all of the invites in the organization.",

		Attributes: map[string]schema.Attribute{
			"invites": schema.SetNestedAttribute{
				MarkdownDescription: "List of invites.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Invite ID.",
							Computed:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "The email address of the individual to whom the invite was sent.",
							Computed:            true,
						},
						"role": schema.StringAttribute{
							MarkdownDescription: "`owner` or `reader`.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "`accepted`, `expired`, or `pending`.",
							Computed:            true,
						},
						"invited_at": schema.Int64Attribute{
							MarkdownDescription: "The Unix timestamp (in seconds) of when the invite was sent.",
							Computed:            true,
						},
						"expires_at": schema.Int64Attribute{
							MarkdownDescription: "The Unix timestamp (in seconds) of when the invite expires.",
							Computed:            true,
						},
						"accepted_at": schema.Int64Attribute{
							MarkdownDescription: "The Unix timestamp (in seconds) of when the invite was accepted.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *InvitesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InvitesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var invites []apiclient.Invite
	params := &apiclient.ListInvitesParams{
		Limit: ptr.Ptr(100),
	}

	for {
		httpResp, err := d.client.ListInvitesWithResponse(
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

		invites = append(invites, httpResp.JSON200.Data...)

		if httpResp.JSON200.HasMore == nil || !*httpResp.JSON200.HasMore {
			break
		}

		params.After = httpResp.JSON200.LastId
	}

	if err := data.Fill(invites); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to fill data: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
