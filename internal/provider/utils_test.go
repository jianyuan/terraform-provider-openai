package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestMergeDiagnostics(t *testing.T) {
	var diags diag.Diagnostics

	do := func() (string, diag.Diagnostics) {
		var diags diag.Diagnostics
		diags.AddError("Error summary", "Error detail")
		diags.AddWarning("Warning summary", "Warning detail")
		return "Result", diags
	}

	v := mergeDiagnostics(do())(&diags)

	if v != "Result" {
		t.Errorf("Expected Result, got %s", v)
	}

	var expectedDiags diag.Diagnostics
	expectedDiags.AddError("Error summary", "Error detail")
	expectedDiags.AddWarning("Warning summary", "Warning detail")

	if !diags.Equal(expectedDiags) {
		t.Errorf("Expected %v, got %v", expectedDiags, diags)
	}
}
