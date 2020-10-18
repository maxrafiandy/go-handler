package handler

import "net/http"

// private constants
const (
	formatDate              string = "2006-01-02"
	formatDateYMD           string = "20060102"
	contentType             string = "Content-Type"
	contentLength           string = "Content-Length"
	imageJPG                string = "image/jpeg"
	imagePNG                string = "image/png"
	invalidContentType      string = "Invalid content type or the request contains empty body"
	invalidDateFormat       string = "Invalid date format"
	decodeFail              string = "Unable to decode file content. The file format is not in jpg neither png"
	noImagePath             string = "assets/no-image.png"
	contentSecurityPolicy   string = "Content-Security-Policy"
	strictTransportSecurity string = "Strict-Transport-Security"

	index string = ""
	subID string = "/{id}"

	id          string = "id"
	restful     string = "rest"
	logicalTrue string = "true"

	get    string = http.MethodGet
	post   string = http.MethodPost
	put    string = http.MethodPut
	patch  string = http.MethodPatch
	delete string = http.MethodDelete

	// 10MB
	defaultMaxMemory int64 = 10 << 20
)

// exported constants
const (
	// MessageOK holds default message for Status Code 200
	MessageOK = "OK"

	// MessageCreated holds default message for Status Code 201
	MessageCreated = "Created"

	// MessageAccepted holds default message for Status Code 202
	MessageAccepted = "Accepted"

	// MessageNoContent holds default message for Status Code 204
	MessageNoContent = "No-content"

	// MessageBadRequest holds default message for Status Code 400
	MessageBadRequest = "Bad request"

	// MessageUnauthorized holds default message for Status Code 401
	MessageUnauthorized = "Unauthorized"

	// MessageForbidden holds default message for Status Code 403
	MessageForbidden = "Forbidden"

	// MessageNotFound holds default message for Status Code 404
	MessageNotFound = "Record not found"

	// MessagePageNotFound holds default message for Status Code 404
	MessagePageNotFound = "Page not found"

	// MessageNotImplemented holds default message for Status Code 404
	MessageNotImplemented = "Not implemented"

	// MessageMethodNotAllowed holds default message for Status Code 405
	MessageMethodNotAllowed = "Method not allowed"

	// MessageConflict holds default message for Status Code 409
	MessageConflict = "Conflict"

	// MessageInternalServerError holds default message for Status Code 500
	MessageInternalServerError = "Internal server error"

	// MessageUpdated holds default message for updated
	MessageUpdated = "Updated"

	// MessageDeleted holds default message for deleted
	MessageDeleted = "Deleted"
)

// private variables
var (
	errNotImplemented Error = Error{
		Description: MessageNotImplemented,
	}
	errBadRequest Error = Error{
		Description: MessageBadRequest,
	}
	errNotFound Error = Error{
		Description: MessageNotFound,
	}
	errPageNotFound Error = Error{
		Description: MessagePageNotFound,
	}
	errConflict Error = Error{
		Description: MessageConflict,
	}
	errForbidden Error = Error{
		Description: MessageForbidden,
	}
	errUnauthorized Error = Error{
		Description: MessageUnauthorized,
	}
	errInternalServerError Error = Error{
		Description: MessageInternalServerError,
	}
	errNotAllowed Error = Error{
		Description: MessageMethodNotAllowed,
	}
)

// public variables
var (
	indexMethods []string = []string{get, post, put, delete, patch}
	subIDMethods []string = []string{get, put, patch, delete}
)
