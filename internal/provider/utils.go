package provider

import (
	"fmt"
	"regexp"
	"strings"
)

func Pointer[T any](v T) *T {
	return &v
}

func BuildTwoPartId(a, b string) string {
	return fmt.Sprintf("%s/%s", a, b)
}

func SplitTwoPartId(id, a, b string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected %s/%s", id, a, b)
	}
	return parts[0], parts[1], nil
}

func BuildThreePartId(a, b, c string) string {
	return fmt.Sprintf("%s/%s/%s", a, b, c)
}

func SplitThreePartId(id, a, b, c string) (string, string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("unexpected format of ID (%s), expected %s/%s/%s", id, a, b, c)
	}
	return parts[0], parts[1], parts[2], nil
}

func MaskToRegex(mask string) string {
	// Escape special regex characters in the mask
	escapedMask := regexp.QuoteMeta(mask)
	// Replace \* (escaped asterisks) with .* (regex pattern for any character sequence)
	regexPattern := strings.ReplaceAll(escapedMask, "\\*", ".*")
	return "^" + regexPattern + "$" // Ensure the pattern matches the entire string
}

func MatchStringWithMask(target, mask string) bool {
	regexPattern := MaskToRegex(mask)
	regex := regexp.MustCompile(regexPattern)
	return regex.MatchString(target)
}
