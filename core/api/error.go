package api

import (
	"fmt"
)

const (
	// ErrorCodeDatabaseFailure error code for database failure
	ErrorCodeDatabaseFailure = "Key_DBQueryFailure"
	// ErrorCodeInternalError error code for internal error
	ErrorCodeInternalError = "Key_InternalError"
	// ErrorCodeInvalidFields error code for invalid fields
	ErrorCodeInvalidFields = "Key_InvalidFields"
	// ErrorCodeInvalidRequestPayload error code for invalid request payload
	ErrorCodeInvalidRequestPayload = "Key_InvalidRequestPayload"
	// ErrorCodeResourceNotFound error code for invalid request payload
	ErrorCodeResourceNotFound = "Key_ResourceNotFound"
	// ErrorCodeUnexpected error code for invalid request payload
	ErrorCodeUnexpected = "Key_Unexpected"
	// ErrorCodePaymentGateway error code for payment gateway failure
	ErrorCodePaymentGateway = "Key_PaymentGatewayFailure"
)

// NewHTTPResourceNotFound creates an new instance of HTTP Error
func NewHTTPResourceNotFound(resourceName,
	resourceValue, errorMessage string) HTTPResourceNotFound {
	return HTTPResourceNotFound{
		ErrorCodeResourceNotFound,
		resourceName,
		resourceValue,
		errorMessage,
	}
}

// HTTPResourceNotFound represents HTTP 404 error
type HTTPResourceNotFound struct {
	ErrorKey      string `json:"error"`
	ResourceName  string `json:"resource_name"`
	ResourceValue string `json:"resource_value"`
	ErrorMessage  string `json:"message"`
}

// Error returns the error string
func (e HTTPResourceNotFound) Error() string {
	return e.ErrorKey
}

// NewHTTPError creates an new instance of HTTP Error
func NewHTTPError(err, message string) HTTPError {
	return HTTPError{ErrorKey: err, ErrorMessage: message}
}

// HTTPError Represent an error to be sent back on response
type HTTPError struct {
	ErrorKey     string `json:"error"`
	ErrorMessage string `json:"message"`
}

// Error returns the error string
func (err HTTPError) Error() string {
	return fmt.Sprintf("key : %s message %s", err.ErrorKey, err.ErrorMessage)
}

type Error interface {
	Error() string
}
