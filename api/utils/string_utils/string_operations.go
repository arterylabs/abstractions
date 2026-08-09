package string_utils

// StringArrayContainsString reports whether valueToCheckFor is present in arrayToCheck.
func StringArrayContainsString(arrayToCheck []string, valueToCheckFor string) bool {
	for _, val := range arrayToCheck {
		if val == valueToCheckFor {
			return true
		}
	}
	return false
}
