package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectResourceModel) Fill(ctx context.Context, project openai.Project) diag.Diagnostics {
	m.Id = supertypes.NewStringValue(project.ID)
	m.Name = supertypes.NewStringValue(project.Name)
	m.Status = supertypes.NewStringValue(project.Status)
	m.ExternalKeyId = (func() supertypes.StringValue {
		if project.JSON.ExternalKeyID.Valid() {
			return supertypes.NewStringValue(project.ExternalKeyID)
		}
		return supertypes.NewStringNull()
	})()
	m.CreatedAt = supertypes.NewInt64Value(project.CreatedAt)
	m.ArchivedAt = (func() supertypes.Int64Value {
		if project.JSON.ArchivedAt.Valid() {
			return supertypes.NewInt64Value(project.ArchivedAt)
		}
		return supertypes.NewInt64Null()
	})()

	return nil
}

func (r *ProjectResource) getCreateJSONRequestBody(ctx context.Context, data ProjectResourceModel) (*openai.AdminOrganizationProjectNewParams, diag.Diagnostics) {
	body := &openai.AdminOrganizationProjectNewParams{
		Name: data.Name.ValueString(),
		ExternalKeyID: (func() param.Opt[string] {
			if data.ExternalKeyId.IsKnown() {
				return openai.String(data.ExternalKeyId.ValueString())
			}
			return param.Opt[string]{}
		})(),
		Geography: (func() param.Opt[string] {
			if data.Geography.IsKnown() {
				return openai.String(data.Geography.ValueString())
			}
			return param.Opt[string]{}
		})(),
	}

	return body, nil
}

func (r *ProjectResource) getUpdateJSONRequestBody(ctx context.Context, data ProjectResourceModel) (*openai.AdminOrganizationProjectUpdateParams, diag.Diagnostics) {
	body := &openai.AdminOrganizationProjectUpdateParams{
		Name: openai.String(data.Name.ValueString()),
		ExternalKeyID: (func() param.Opt[string] {
			if data.ExternalKeyId.IsKnown() {
				return openai.String(data.ExternalKeyId.ValueString())
			}
			return param.Opt[string]{}
		})(),
	}

	return body, nil
}
