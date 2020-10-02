package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RestHandlers interface
type RestHandlers interface {
	Post() interface{}
	Get() interface{}
	GetID(string) interface{}
	Put() interface{}
	PutID(string) interface{}
	Patch() interface{}
	PatchID(string) interface{}
	Delete() interface{}
	DeleteID(string) interface{}
	reset(http.ResponseWriter, *http.Request)
}

// REST maps router to appropriate methods
func REST(rest RestHandlers, ctx *Context) interface{} {
	id, withid := mux.Vars(ctx.Request)[id]

	rest.reset(ctx.Writer, ctx.Request)

	if !withid { // route to /
		switch ctx.Request.Method {
		case get:
			return rest.Get()
		case post:
			return rest.Post()
		case put:
			return rest.Put()
		case patch:
			return rest.Patch()
		case delete:
			return rest.Delete()
		}
	} else { // route to /{id:[0-9]+}
		switch ctx.Request.Method {
		case get:
			return rest.GetID(id)
		case put:
			return rest.PutID(id)
		case patch:
			return rest.PatchID(id)
		case delete:
			return rest.DeleteID(id)
		}
	}
	// this line may never reached
	return &Error{Description: "REST EOF!"}
}
