package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectResourceModel) Fill(ctx context.Context, project apiclient.Project) diag.Diagnostics {
	m.Id = supertypes.NewStringValue(project.Id)
	m.Name = supertypes.NewStringValue(project.Name)
	m.Status = supertypes.NewStringValue(string(project.Status))
	m.ExternalKeyId = supertypes.NewStringPointerValue(project.ExternalKeyId)
	m.CreatedAt = supertypes.NewInt64Value(project.CreatedAt)
	m.ArchivedAt = supertypes.NewInt64PointerValue(project.ArchivedAt)
	return nil
}

func (r *ProjectResource) getCreateJSONRequestBody(ctx context.Context, data ProjectResourceModel) (apiclient.CreateProjectJSONRequestBody, diag.Diagnostics) {
	return apiclient.CreateProjectJSONRequestBody{
		Name:          data.Name.ValueString(),
		ExternalKeyId: data.ExternalKeyId.ValueStringPointer(),
	}, nil
}

func (r *ProjectResource) getUpdateJSONRequestBody(ctx context.Context, data ProjectResourceModel) (apiclient.ModifyProjectJSONRequestBody, diag.Diagnostics) {
	return apiclient.ModifyProjectJSONRequestBody{
		Name: data.Name.ValueString(),
	}, nil
}
