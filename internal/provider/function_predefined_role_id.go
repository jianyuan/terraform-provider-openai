package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &PredefinedRoleIdFunction{}

func NewPredefinedRoleIdFunction() function.Function {
	return &PredefinedRoleIdFunction{}
}

type PredefinedRoleIdFunction struct{}

func (f *PredefinedRoleIdFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "predefined_role_id"
}

func (f *PredefinedRoleIdFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Get the ID of a predefined role",
		Description: "Returns the ID of a predefined role.",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "role",
				Description: "The role of the predefined role. `owner` or `reader`.",
				Validators: []function.StringParameterValidator{
					stringvalidator.OneOf("owner", "reader"),
				},
			},
			function.StringParameter{
				Name:        "organization_id",
				Description: "The ID of the organization.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *PredefinedRoleIdFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var role, organizationId string

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &role, &organizationId))

	output := fmt.Sprintf("role-api-organization-%s__api-organization__%s", role, organizationId)

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, output))
}
