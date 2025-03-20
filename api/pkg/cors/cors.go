package cors

import (
	"net/http"
	"strings"
)

const (
	HeaderAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAllowCredentials = "Access-Control-Allow-Credentials"
)

func SetHeaderAllowedMethods(w http.ResponseWriter, methods ...string) {
	w.Header().Set(HeaderAllowMethods, http.MethodOptions+", "+strings.Join(methods, ", "))
}
