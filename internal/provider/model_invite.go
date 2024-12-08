package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

type InviteModel struct {
	Id         types.String `tfsdk:"id"`
	Email      types.String `tfsdk:"email"`
	Role       types.String `tfsdk:"role"`
	Status     types.String `tfsdk:"status"`
	InvitedAt  types.Int64  `tfsdk:"invited_at"`
	ExpiresAt  types.Int64  `tfsdk:"expires_at"`
	AcceptedAt types.Int64  `tfsdk:"accepted_at"`
}

func (m *InviteModel) Fill(i apiclient.Invite) error {
	m.Id = types.StringValue(i.Id)
	m.Email = types.StringValue(i.Email)
	m.Role = types.StringValue(string(i.Role))
	m.Status = types.StringValue(string(i.Status))
	m.InvitedAt = types.Int64Value(int64(i.InvitedAt))
	m.ExpiresAt = types.Int64Value(int64(i.ExpiresAt))
	if i.AcceptedAt == nil {
		m.AcceptedAt = types.Int64Null()
	} else {
		m.AcceptedAt = types.Int64Value(int64(*i.AcceptedAt))
	}
	return nil
}
