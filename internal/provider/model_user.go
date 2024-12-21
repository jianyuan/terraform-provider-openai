package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

type UserModel struct {
	Id      types.String `tfsdk:"id"`
	Email   types.String `tfsdk:"email"`
	Name    types.String `tfsdk:"name"`
	Role    types.String `tfsdk:"role"`
	AddedAt types.Int64  `tfsdk:"added_at"`
}

func (m *UserModel) Fill(ctx context.Context, u apiclient.User) (diags diag.Diagnostics) {
	m.Id = types.StringValue(u.Id)
	m.Email = types.StringValue(u.Email)
	m.Name = types.StringValue(u.Name)
	m.Role = types.StringValue(string(u.Role))
	m.AddedAt = types.Int64Value(int64(u.AddedAt))
	return
}
