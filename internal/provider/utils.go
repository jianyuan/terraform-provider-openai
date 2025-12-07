package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// deduplicate removes duplicates from a slice of any comparable type T.
func deduplicate[T comparable](input []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0, len(input))

	for _, v := range input {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

func mergeDiagnostics[T any](v T, diagsOut diag.Diagnostics) func(diags *diag.Diagnostics) T {
	return func(diags *diag.Diagnostics) T {
		diags.Append(diagsOut...)
		return v
	}
}

func getBool(v any) bool {
	switch v := v.(type) {
	case bool:
		return v
	case *bool:
		if v == nil {
			return false
		}
		return *v
	default:
		panic("unknown type")
	}
}

func getString(v any) string {
	switch v := v.(type) {
	case string:
		return v
	case *string:
		if v == nil {
			return ""
		}
		return *v
	default:
		panic("unknown type")
	}
}
