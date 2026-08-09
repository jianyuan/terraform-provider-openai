package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *ProjectModelPermissionsDataSourceModel) Fill(ctx context.Context, data openai.ProjectModelPermissions) diag.Diagnostics {
	m.Mode = supertypes.NewStringValue(string(data.Mode))
	m.ModelIds = supertypes.NewSetValueOfSlice(ctx, data.ModelIDs)
	return nil
}
