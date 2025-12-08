package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &PredefinedProjectRoleIdFunction{}

func NewPredefinedProjectRoleIdFunction() function.Function {
	return &PredefinedProjectRoleIdFunction{}
}

type PredefinedProjectRoleIdFunction struct{}

func (f *PredefinedProjectRoleIdFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "predefined_project_role_id"
}

func (f *PredefinedProjectRoleIdFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Get the ID of a predefined project role",
		Description: "Returns the ID of a predefined project role.",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "role",
				Description: "The role of the predefined role. `owner`, `member` or `viewer`.",
				Validators: []function.StringParameterValidator{
					stringvalidator.OneOf("owner", "member", "viewer"),
				},
			},
			function.StringParameter{
				Name:        "project_id",
				Description: "The ID of the project.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *PredefinedProjectRoleIdFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var role, projectId string

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &role, &projectId))

	output := fmt.Sprintf("role-api-project-%s__api-project__%s", role, projectId)

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, output))
}
