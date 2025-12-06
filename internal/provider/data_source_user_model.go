package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-openai/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *UserDataSourceModel) Fill(ctx context.Context, user apiclient.User) diag.Diagnostics {
	m.Id = supertypes.NewStringValue(user.Id)
	m.Email = supertypes.NewStringValue(user.Email)
	m.Name = supertypes.NewStringValue(user.Name)
	m.Role = supertypes.NewStringValue(string(user.Role))
	m.AddedAt = supertypes.NewInt64Value(user.AddedAt)
	return nil
}
