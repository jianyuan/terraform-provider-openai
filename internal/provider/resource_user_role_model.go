package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
)

func (m UserRoleResourceModel) Fill(ctx context.Context, user openai.OrganizationUser) diag.Diagnostics {
	return nil
}

func (r *UserRoleResource) getCreateJSONRequestBody(ctx context.Context, data UserRoleResourceModel) (*openai.AdminOrganizationUserUpdateParams, diag.Diagnostics) {
	return r.getUpdateJSONRequestBody(ctx, data)
}

func (r *UserRoleResource) getUpdateJSONRequestBody(ctx context.Context, data UserRoleResourceModel) (*openai.AdminOrganizationUserUpdateParams, diag.Diagnostics) {
	return &openai.AdminOrganizationUserUpdateParams{
		Role: openai.String(data.Role.ValueString()),
	}, nil
}
