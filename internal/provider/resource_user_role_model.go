package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
)

func (m UserRoleResourceModel) Fill(ctx context.Context, user apiclient.User) diag.Diagnostics {
	return nil
}

func (r *UserRoleResource) getCreateJSONRequestBody(data UserRoleResourceModel) apiclient.ModifyUserJSONRequestBody {
	return apiclient.ModifyUserJSONRequestBody{
		Role: apiclient.UserRoleUpdateRequestRole(data.Role.ValueString()),
	}
}

func (r *UserRoleResource) getUpdateJSONRequestBody(data UserRoleResourceModel) apiclient.ModifyUserJSONRequestBody {
	return apiclient.ModifyUserJSONRequestBody{
		Role: apiclient.UserRoleUpdateRequestRole(data.Role.ValueString()),
	}
}
