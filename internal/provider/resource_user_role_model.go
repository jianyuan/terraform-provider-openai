package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

func (m UserRoleResourceModel) Fill(ctx context.Context, user apiclient.User) diag.Diagnostics {
	return nil
}

func (r *UserRoleResource) getCreateJSONRequestBody(ctx context.Context, data UserRoleResourceModel) (apiclient.ModifyUserJSONRequestBody, diag.Diagnostics) {
	return apiclient.ModifyUserJSONRequestBody{
		Role: apiclient.UserRoleUpdateRequestRole(data.Role.ValueString()),
	}, nil
}

func (r *UserRoleResource) getUpdateJSONRequestBody(ctx context.Context, data UserRoleResourceModel) (apiclient.ModifyUserJSONRequestBody, diag.Diagnostics) {
	return apiclient.ModifyUserJSONRequestBody{
		Role: apiclient.UserRoleUpdateRequestRole(data.Role.ValueString()),
	}, nil
}
