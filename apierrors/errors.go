package apierrors

import (
	"net/http"
)

var BadRequest = PublicError{
	HttpCode:  http.StatusBadRequest,
	ErrorCode: "BAD_REQUEST",
	Help:      "An error in the format of your query or in the format of the data is present. The query has not been processed.",
}

var Unauthorized = PublicError{
	HttpCode:  http.StatusUnauthorized,
	ErrorCode: "UNAUTHORIZED",
	Help:      "The use of this endpoint requires authentication. No valid authentication data was provided with your request.",
}

var Forbidden = PublicError{
	HttpCode:  http.StatusForbidden,
	ErrorCode: "FORBIDDEN",
	Help:      "The use of this endpoint is forbidden. You do not have the required permissions to access this resource.",
}

var NotFound = PublicError{
	HttpCode:  http.StatusNotFound,
	ErrorCode: "NOT_FOUND",
	Help:      "The requested resource could not be found.",
}

var InvalidJson = PublicError{
	HttpCode:  http.StatusBadRequest,
	ErrorCode: "INVALID_JSON",
	Help:      "The JSON provided in the request body is invalid.",
}

var EmailNotRegistered = PublicError{
	HttpCode:  http.StatusNotFound,
	ErrorCode: "EMAIL_NOT_REGISTERED",
	Help:      "No account is associated with this email address.",
}
