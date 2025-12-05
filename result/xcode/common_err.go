package xcode

import (
	"net/http"
)

const (
	SERVER_COMMON_ERROR = 100001
	REQUEST_PARAM_ERROR = 100002
	DB_ERROR            = 100003
)

var (
	UserNotFound = New(http.StatusUnauthorized, "User not login. ")
	TokenInvalid = New(http.StatusPaymentRequired, "Token invalid. ")
)
