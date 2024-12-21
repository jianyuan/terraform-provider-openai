package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = &InviteDataSource{}

func NewInviteDataSource() datasource.DataSource {
	return &InviteDataSource{}
}

type InviteDataSource struct {
	baseDataSource
}

func (d *InviteDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_invite"
}

func (d *InviteDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves an invite.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Invite ID.",
				Required:            true,
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
	}
}

func (d *InviteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InviteModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := d.client.RetrieveInviteWithResponse(
		ctx,
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	}

	if httpResp.JSON200 == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to read, got empty response")
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
