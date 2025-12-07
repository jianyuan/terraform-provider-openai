package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectUserResourceModel) Fill(ctx context.Context, data apiclient.ProjectUser) diag.Diagnostics {
	m.UserId = supertypes.NewStringValue(data.Id)
	m.Role = supertypes.NewStringValue(string(data.Role))
	return nil
}

func (r *ProjectUserResource) getCreateJSONRequestBody(ctx context.Context, data ProjectUserResourceModel) (apiclient.CreateProjectUserJSONRequestBody, diag.Diagnostics) {
	return apiclient.CreateProjectUserJSONRequestBody{
		Role:   apiclient.ProjectUserCreateRequestRole(data.Role.ValueString()),
		UserId: data.UserId.ValueString(),
	}, nil
}

func (r *ProjectUserResource) getUpdateJSONRequestBody(ctx context.Context, data ProjectUserResourceModel) (apiclient.ModifyProjectUserJSONRequestBody, diag.Diagnostics) {
	return apiclient.ModifyProjectUserJSONRequestBody{
		Role: apiclient.ProjectUserUpdateRequestRole(data.Role.ValueString()),
	}, nil
}
