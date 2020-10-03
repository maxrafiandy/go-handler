package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter with default handler
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.NotFoundHandler = PageNotFound404{}
	r.MethodNotAllowedHandler = MethodNotAllowed405{}
	return r
}

type (
	// PageNotFound404 will implement the ServeHTTP func
	// so than this struct can used to handle the 404
	// page not found error
	PageNotFound404 struct{}

	// MethodNotAllowed405 will implement the ServeHTTP func
	// so than this struct can used to handle the 405
	// method not allowed
	MethodNotAllowed405 struct{}
)

// ServeHTTP impementation of PageNotFound404
func (e PageNotFound404) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentType, "application/json")
	w.WriteHeader(http.StatusNotFound)
	encoder := json.NewEncoder(w)
	encoder.Encode(&Response{Message: MessagePageNotFound, Data: nil})
}

// ServeHTTP impementation of MethodNotAllowed405
func (e MethodNotAllowed405) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentType, "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	encoder := json.NewEncoder(w)
	encoder.Encode(&Response{Message: MessageMethodNotAllowed, Data: nil})
}

// Group create the sub router of path
func Group(path string, parent *mux.Router) *mux.Router {
	return parent.PathPrefix(path).Subrouter()
}
