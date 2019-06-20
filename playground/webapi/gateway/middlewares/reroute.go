package middlewares

import (
	"strings"
)

const apiPreFix = "/api/"

func ConvertURL(src string) string {
	trim := strings.TrimPrefix(src, apiPreFix)
	return "http://localhost:8080/" + trim
}
