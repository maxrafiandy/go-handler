package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Restful interface
type Restful interface {
	Handler(http.ResponseWriter, *http.Request)
}

// RestContext inherits Restful as parent
type RestContext struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

var (
	errMethodNotAllowed Error = Error{
		Description: MessageMethodNotAllowed,
	}
)

// RestHandlers interface
type RestHandlers interface {
	Post()
	Get()
	GetID(string)
	Put()
	PutID(string)
	Patch()
	PatchID(string)
	Delete()
	DeleteID(string)
}

// MapRest maps router to appropriate methods
func MapRest(rest RestHandlers, w http.ResponseWriter, r *http.Request) {
	id, withid := mux.Vars(r)[id]

	if !withid { // route to /
		switch r.Method {
		case get:
			rest.Get()
		case post:
			rest.Post()
		case put:
			rest.Put()
		case patch:
			rest.Patch()
		case delete:
			rest.Delete()
		default:
			methodNotAllowed(w)
		}
	} else { // route to /{id:[0-9]+}
		switch r.Method {
		case get:
			rest.GetID(id)
		case put:
			rest.PutID(id)
		case patch:
			rest.PatchID(id)
		case delete:
			rest.DeleteID(id)
		default:
			methodNotAllowed(w)
		}
	}
}

// InitRest initialize restful writer and request
func (rest RestContext) InitRest(h RestHandlers, w http.ResponseWriter, r *http.Request) {
	rest.Writer = w
	rest.Request = r
	MapRest(h, w, r)
}

// Created send success response when a record was successfully created
func (rest RestContext) Created(data interface{}) {
	response(rest.Writer, MessageCreated, data, http.StatusCreated)
}

// Success send success response with result data
func (rest RestContext) Success(data interface{}) {
	response(rest.Writer, MessageOK, data, http.StatusOK)
}

// NoContent send success response without any content.
// this method usually used with update method.
func (rest RestContext) NoContent() {
	response(rest.Writer, MessageNoContent, nil, http.StatusNoContent)
}

// BadRequest send general 400-bad request
func (rest RestContext) BadRequest(data interface{}) {
	response(rest.Writer, MessageBadRequest, data, http.StatusBadRequest)
}

// NotFound send general 404-Not found.
// this method is equal to PageNotFound() and usually
// used when record was not found in collection instead of
// return a page not found message
func (rest RestContext) NotFound() {
	response(rest.Writer, MessageNotFound, nil, http.StatusNotFound)
}

// InternalServerError send general 500-interal server error
func (rest RestContext) InternalServerError(data interface{}) {
	response(rest.Writer, MessageInternalServerError, data, http.StatusInternalServerError)
}

// PageNotFound send general 404-not found.
// this method is equal to NotFound() but returns
// page not found message
func (rest RestContext) PageNotFound() {
	response(rest.Writer, MessagePageNotFound, nil, http.StatusNotFound)
}

// MethodNotAllowed send general 405-method not allowed
func (rest RestContext) MethodNotAllowed() {
	response(rest.Writer, MessageMethodNotAllowed, errMethodNotAllowed, http.StatusMethodNotAllowed)
}

// Get implements Get() function
func (rest RestContext) Get() {
	rest.MethodNotAllowed()
}

// GetID implements Find() function
func (rest RestContext) GetID(id string) {
	rest.MethodNotAllowed()
}

// Post implements Post() function
func (rest RestContext) Post() {
	rest.MethodNotAllowed()
}

// PutID implements Put() function
func (rest RestContext) PutID(id string) {
	rest.MethodNotAllowed()
}

// Put implements Put() function
func (rest RestContext) Put() {
	rest.MethodNotAllowed()
}

// PatchID implements PatchID() function
func (rest RestContext) PatchID(id string) {
	rest.MethodNotAllowed()
}

// Patch implements Patch() function
func (rest RestContext) Patch() {
	rest.MethodNotAllowed()
}

// DeleteID implements DeleteID() function
func (rest RestContext) DeleteID(id string) {
	rest.MethodNotAllowed()
}

// Delete implements Delete() function
func (rest RestContext) Delete() {
	rest.MethodNotAllowed()
}
