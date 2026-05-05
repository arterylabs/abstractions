package security_utils

import "strings"

func SanitizeSQLString(value string) string {
	return strings.Replace(value, "'", "''", -1)
}
