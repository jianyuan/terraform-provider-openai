package diag

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func DiagnosticsError(diags diag.Diagnostics) error {
	errList := diags.Errors()
	if len(errList) == 0 {
		return nil
	}

	errs := make([]error, len(errList))
	for i, d := range errList {
		errs[i] = errors.New(FormatDiagnostic(d))
	}

	return errors.Join(errs...)
}

func FormatDiagnostic(d diag.Diagnostic) string {
	msg := d.Summary()

	if detail := d.Detail(); detail != "" {
		msg += "\n\n" + detail
	}

	if p, ok := d.(diag.DiagnosticWithPath); ok {
		msg += "\n" + p.Path().String()
	}

	return msg
}
