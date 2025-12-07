package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectServiceAccountResourceModel) Fill(ctx context.Context, data any) diag.Diagnostics {
	var diag diag.Diagnostics

	switch data := data.(type) {
	case apiclient.ProjectServiceAccountCreateResponse:
		m.Id = supertypes.NewStringValue(data.Id)
		m.Name = supertypes.NewStringValue(data.Name)
		m.Role = supertypes.NewStringValue(string(data.Role))
		m.CreatedAt = supertypes.NewInt64Value(data.CreatedAt)
		m.ApiKeyId = supertypes.NewStringValue(data.ApiKey.Id)
		m.ApiKey = supertypes.NewStringValue(data.ApiKey.Value)
		return nil
	case apiclient.ProjectServiceAccount:
		m.Id = supertypes.NewStringValue(data.Id)
		m.Name = supertypes.NewStringValue(data.Name)
		m.Role = supertypes.NewStringValue(string(data.Role))
		m.CreatedAt = supertypes.NewInt64Value(data.CreatedAt)
		return nil
	default:
		diag.AddError("Unknown type", fmt.Sprintf("Unknown type: %T", data))
		return diag
	}
}

func (r *ProjectServiceAccountResource) getCreateJSONRequestBody(ctx context.Context, data ProjectServiceAccountResourceModel) (apiclient.ProjectServiceAccountCreateRequest, diag.Diagnostics) {
	return apiclient.CreateProjectServiceAccountJSONRequestBody{
		Name: data.Name.ValueString(),
	}, nil
}
