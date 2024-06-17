package provider

import "testing"

func TestMatchStringWithMask(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		target   string
		mask     string
		expected bool
	}{
		{
			target:   "sk-my-secret-key",
			mask:     "sk-my-secret-key",
			expected: true,
		},
		{
			target:   "sk-my-secret-key",
			mask:     "sk-my-se************************************************************key",
			expected: true,
		},
		{
			target:   "sk-my-secret-key",
			mask:     "sk-my-*-key",
			expected: true,
		},
		{
			target:   "sk-my-secret-key",
			mask:     "sk-my-secret-key-2",
			expected: false,
		},
		{
			target:   "sk-my-secret-key",
			mask:     "sk-my-*-key-2",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.target, func(t *testing.T) {
			actual := MatchStringWithMask(tc.target, tc.mask)
			if actual != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
