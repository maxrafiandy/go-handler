package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

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
	reset(http.ResponseWriter, *http.Request)
}

// REST maps router to appropriate methods
func REST(rest RestHandlers, ctx *Context) {
	id, withid := mux.Vars(ctx.Request)[id]

	rest.reset(ctx.Writer, ctx.Request)

	if !withid { // route to /
		switch ctx.Request.Method {
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
			response(ctx.Writer, MessageMethodNotAllowed, nil, http.StatusMethodNotAllowed)
		}
	} else { // route to /{id:[0-9]+}
		switch ctx.Request.Method {
		case get:
			rest.GetID(id)
		case put:
			rest.PutID(id)
		case patch:
			rest.PatchID(id)
		case delete:
			rest.DeleteID(id)
		default:
			response(ctx.Writer, MessageMethodNotAllowed, nil, http.StatusMethodNotAllowed)
		}
	}
}
