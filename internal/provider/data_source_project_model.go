package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectDataSourceModel) Fill(ctx context.Context, project openai.Project) diag.Diagnostics {
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
