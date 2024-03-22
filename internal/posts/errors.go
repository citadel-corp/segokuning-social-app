package posts

import (
	"net/http"
)

var (
	ErrorForbidden     = Response{Code: http.StatusForbidden, Message: "Forbidden"}
	ErrorUnauthorized  = Response{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrorRequiredField = Response{Code: http.StatusBadRequest, Message: "Required field"}
	ErrorInternal      = Response{Code: http.StatusInternalServerError, Message: "Internal Server Error"}
	ErrorBadRequest    = Response{Code: http.StatusBadRequest, Message: "Bad Request"}
	ErrorNoRecords     = Response{Code: http.StatusOK, Message: "No records found"}
	ErrorNotFound      = Response{Code: http.StatusNotFound, Message: "No records found"}
)
