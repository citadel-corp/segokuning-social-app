package userfriends

import (
	"errors"
	"net/http"
)

var (
	ErrorForbidden     = Response{Code: http.StatusForbidden, Message: "Forbidden", Error: errors.New("Forbidden")}
	ErrorUnauthorized  = Response{Code: http.StatusUnauthorized, Message: "Unauthorized", Error: errors.New("Unauthorized")}
	ErrorRequiredField = Response{Code: http.StatusBadRequest, Message: "Required field"}
	ErrorInternal      = Response{Code: http.StatusInternalServerError, Message: "Internal Server Error", Error: errors.New("Internal Server Error")}
	ErrorBadRequest    = Response{Code: http.StatusBadRequest, Message: "Bad Request"}
	ErrorNoRecords     = Response{Code: http.StatusOK, Message: "No records found"}
	ErrorNotFound      = Response{Code: http.StatusNotFound, Message: "No records found"}

	ErrFriendAlreadyExists = Response{Code: http.StatusBadRequest, Message: "Friend had already been added"}
	ErrFriendNotExists     = Response{Code: http.StatusNotFound, Message: "Friend is not found"}
	ErrCannotAddSelf       = Response{Code: http.StatusBadRequest, Message: "Cannot add self as friend"}
	ErrCannotDeleteSelf    = Response{Code: http.StatusBadRequest, Message: "Cannot delete self as friend"}
	ErrNotFriend           = Response{Code: http.StatusBadRequest, Message: "Cannot delete non friend"}
)
