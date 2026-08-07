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

	if project.JSON.ExternalKeyID.Valid() {
		m.ExternalKeyId = supertypes.NewStringValue(project.ExternalKeyID)
	} else {
		m.ExternalKeyId = supertypes.NewStringNull()
	}

	m.CreatedAt = supertypes.NewInt64Value(project.CreatedAt)

	if project.JSON.ArchivedAt.Valid() {
		m.ArchivedAt = supertypes.NewInt64Value(project.ArchivedAt)
	} else {
		m.ArchivedAt = supertypes.NewInt64Null()
	}

	return nil
}
