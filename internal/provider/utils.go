package provider

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/openai/openai-go/v3"
)

func mergeDiagnostics[T any](v T, diagsOut diag.Diagnostics) func(diags *diag.Diagnostics) T {
	return func(diags *diag.Diagnostics) T {
		diags.Append(diagsOut...)
		return v
	}
}

// clientErrorDetail describes a failed API call for a diagnostic. openai.Error's own message
// is only the HTTP status, so this adds the explanation the API returned.
func clientErrorDetail(err error) string {
	if apiErr, ok := errors.AsType[*openai.Error](err); ok && apiErr.Message != "" {
		return fmt.Sprintf("%s: %s", apiErr.Error(), apiErr.Message)
	}
	return err.Error()
}
