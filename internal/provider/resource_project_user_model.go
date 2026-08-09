package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
)

func (r *ProjectUserResource) getNewParams(ctx context.Context, data ProjectUserResourceModel) (*openai.AdminOrganizationProjectUserNewParams, diag.Diagnostics) {
	return &openai.AdminOrganizationProjectUserNewParams{
		Role:   data.Role.ValueString(),
		UserID: openai.String(data.UserId.ValueString()),
	}, nil
}

func (r *ProjectUserResource) getUpdateParams(ctx context.Context, data ProjectUserResourceModel) (*openai.AdminOrganizationProjectUserUpdateParams, diag.Diagnostics) {
	return &openai.AdminOrganizationProjectUserUpdateParams{
		Role: openai.String(data.Role.ValueString()),
	}, nil
}
