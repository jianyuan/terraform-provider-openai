package openaiparam

import (
	"github.com/openai/openai-go/v3/packages/param"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func FromString(v supertypes.StringValue) param.Opt[string] {
	if v.IsKnown() {
		return param.NewOpt(v.ValueString())
	}
	return param.Opt[string]{}
}

func FromInt64(v supertypes.Int64Value) param.Opt[int64] {
	if v.IsKnown() {
		return param.NewOpt(v.ValueInt64())
	}
	return param.Opt[int64]{}
}
