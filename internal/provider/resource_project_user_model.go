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

func (r *ProjectUserResource) getCreateJSONRequestBody(data ProjectUserResourceModel) apiclient.CreateProjectUserJSONRequestBody {
	return apiclient.CreateProjectUserJSONRequestBody{
		Role:   apiclient.ProjectUserCreateRequestRole(data.Role.ValueString()),
		UserId: data.UserId.ValueString(),
	}
}

func (r *ProjectUserResource) getUpdateJSONRequestBody(data ProjectUserResourceModel) apiclient.ModifyProjectUserJSONRequestBody {
	return apiclient.ModifyProjectUserJSONRequestBody{
		Role: apiclient.ProjectUserUpdateRequestRole(data.Role.ValueString()),
	}
}
