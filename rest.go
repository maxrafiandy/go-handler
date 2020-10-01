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
			response(w, MessageMethodNotAllowed, nil, http.StatusMethodNotAllowed)
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
			response(w, MessageMethodNotAllowed, nil, http.StatusMethodNotAllowed)
		}
	}
}

// Created send success response when a record was successfully created
func (c RestContext) Created(data interface{}) {
	response(c.Writer, MessageCreated, data, http.StatusCreated)
}

// Success send success response with result data
func (c RestContext) Success(data interface{}) {
	response(c.Writer, MessageOK, data, http.StatusOK)
}

// NoContent send success response without any content.
// this method usually used with update method.
func (c RestContext) NoContent() {
	response(c.Writer, MessageNoContent, nil, http.StatusNoContent)
}

// BadRequest send general 400-bad request
func (c RestContext) BadRequest(data error) {
	response(c.Writer, MessageBadRequest, data, http.StatusBadRequest)
}

// NotFound send general 200-Succes withoud data.
// this method is equal to PageNotFound() and usually
// used when record was not found in collection instead of
// return a page not found message
func (c RestContext) NotFound() {
	response(c.Writer, MessageNotFound, nil, http.StatusOK)
}

// PageNotFound send general 404-not found.
// this method is equal to NotFound() but returns
// page not found message
func (c RestContext) PageNotFound() {
	response(c.Writer, MessagePageNotFound, nil, http.StatusNotFound)
}

// InternalServerError send general 500-interal server error
func (c RestContext) InternalServerError(data error) {
	response(c.Writer, MessageInternalServerError, data, http.StatusInternalServerError)
}

// MethodNotAllowed send general 405-method not allowed
func (c RestContext) MethodNotAllowed() {
	response(c.Writer, MessageMethodNotAllowed, errMethodNotAllowed, http.StatusMethodNotAllowed)
}

// Get implements Get() function
func (c RestContext) Get() {
	c.MethodNotAllowed()
}

// GetID implements Find() function
func (c RestContext) GetID(id string) {
	c.MethodNotAllowed()
}

// Post implements Post() function
func (c RestContext) Post() {
	c.MethodNotAllowed()
}

// PutID implements Put() function
func (c RestContext) PutID(id string) {
	c.MethodNotAllowed()
}

// Put implements Put() function
func (c RestContext) Put() {
	c.MethodNotAllowed()
}

// PatchID implements PatchID() function
func (c RestContext) PatchID(id string) {
	c.MethodNotAllowed()
}

// Patch implements Patch() function
func (c RestContext) Patch() {
	c.MethodNotAllowed()
}

// DeleteID implements DeleteID() function
func (c RestContext) DeleteID(id string) {
	c.MethodNotAllowed()
}

// Delete implements Delete() function
func (c RestContext) Delete() {
	c.MethodNotAllowed()
}
