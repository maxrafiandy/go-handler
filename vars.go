package handler

// private constants
const (
	formatDate              string = "2006-01-02"
	formatDateYMD           string = "20060102"
	contentType             string = "Content-Type"
	contentLength           string = "Content-Length"
	imageJPG                string = "image/jpg"
	imagePNG                string = "image/png"
	invalidContentType      string = "Invalid content type"
	invalidDateFormat       string = "Invalid date format"
	decodeFail              string = "Unable to decode file content. The file format is not in jpg neither png"
	noImagePath             string = "assets/no-image.png"
	contentSecurityPolicy   string = "Content-Security-Policy"
	strictTransportSecurity string = "Strict-Transport-Security"

	index string = ""
	subID string = "/{id}"

	id      string = "id"
	get     string = "GET"
	post    string = "POST"
	put     string = "PUT"
	patch   string = "PATCH"
	delete  string = "DELETE"
	restful string = "rest"

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
	MessageNotFound = "Not found"

	// MessagePageNotFound holds default message for Status Code 404
	MessagePageNotFound = "Page not found"

	// MessageNotImplemented holds default message for Status Code 404
	MessageNotImplemented = "Not implemented"

	// MessageMethodNotAllowed holds default message for Status Code 405
	MessageMethodNotAllowed = "Method not allowed"

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
	errNotFound Error = Error{
		Description: MessageNotFound,
	}
	errPageNotFound Error = Error{
		Description: MessagePageNotFound,
	}
)

// public variables
var (
	// GormProps maps database properties
	GormProps map[string]*Gormv1Prop = make(map[string]*Gormv1Prop)
)
