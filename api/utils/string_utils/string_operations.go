package string_utils

// checks if string value  is present in the array to check
func StringArrayContainsString(arrayToCheck []string, valueToCheckFor string) bool {
	for _, val := range arrayToCheck {
		if val == valueToCheckFor {
			return true
		}
	}
	return false
}

func RemoveStringFromArray(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func ConvertStringToBoolean(s string) bool {
	if s == "true" {
		return true
	} else if s == "false" {
		return false
	}
	// Handle other cases or return a default value if needed
	return false
}

func FindStringElementInterfaceArray(list []interface{}, element string) bool {
	itemExists := false
	for _, item := range list {
		if element == item {
			itemExists = true
		}
	}
	return itemExists
}

func DeduplicateArray(array []string) []string {
	var deduplicatedArray = []string{}
	for _, item := range array {
		itemExists := false
		for _, ditem := range deduplicatedArray {
			if ditem == item {
				itemExists = true
			}
		}
		if !itemExists {
			deduplicatedArray = append(deduplicatedArray, item)
		}
	}
	return deduplicatedArray
}
