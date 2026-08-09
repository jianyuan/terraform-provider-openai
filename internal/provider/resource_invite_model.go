package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
)

func (r *InviteResource) getNewParams(ctx context.Context, data InviteResourceModel) (*openai.AdminOrganizationInviteNewParams, diag.Diagnostics) {
	return &openai.AdminOrganizationInviteNewParams{
		Email: data.Email.ValueString(),
		Role:  openai.AdminOrganizationInviteNewParamsRole(data.Role.ValueString()),
	}, nil
}
