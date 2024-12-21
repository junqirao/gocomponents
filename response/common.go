package response

import (
	"net/http"
)

// common code define, using biz code from 0-1000
var (
	CodeDefaultSuccess = NewCode(0, "success", http.StatusOK)
	CodeDefaultFailure = NewCode(500, "failed", http.StatusInternalServerError)

	CodeInvalidParameter = NewCode(400, "invalid parameter", http.StatusBadRequest)
	CodeUnauthorized     = NewCode(401, "unauthorized", http.StatusUnauthorized)
	CodePermissionDeny   = NewCode(403, "permission deny", http.StatusForbidden)
	CodeNotFound         = NewCode(404, "not found", http.StatusNotFound)
	CodeConflict         = NewCode(409, "resource conflict", http.StatusConflict)
	CodeTooManyRequests  = NewCode(429, "too many requests", http.StatusTooManyRequests)
	CodeLocked           = NewCode(423, "resource locked", http.StatusLocked)
	CodeUnImplemented    = NewCode(501, "not implemented", http.StatusNotImplemented)
	CodeBadGateway       = NewCode(502, "bad gateway", http.StatusBadGateway)
	CodeUnavailable      = NewCode(503, "service unavailable", http.StatusServiceUnavailable)
	CodeGatewayTimeout   = NewCode(504, "gateway timeout", http.StatusGatewayTimeout)
)

func CodeFromHttpStatus(status int) *Code {
	text := http.StatusText(status)
	if text == "" {
		text = "unknown"
	}
	return NewCode(status, text, status)
}
