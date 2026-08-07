package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func mergeDiagnostics[T any](v T, diagsOut diag.Diagnostics) func(diags *diag.Diagnostics) T {
	return func(diags *diag.Diagnostics) T {
		diags.Append(diagsOut...)
		return v
	}
}
