package xcode

import (
	"net/http"
)

var (
	UserNotFound = New(http.StatusUnauthorized, "User not login. ")
	TokenInvalid = New(http.StatusPaymentRequired, "Token invalid. ")
)
