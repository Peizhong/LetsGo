package gateway

import "strings"

func ConvertURL(src string) string {
	trim := strings.TrimPrefix(src, apiPreFix)
	return "http://localhost:8080/" + trim
}
