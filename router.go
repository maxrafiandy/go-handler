package handler

import (
	"encoding/json"
	"fmt"
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

// PageNotFound404 will implement the ServeHTTP func
// so than this struct can used to handle the 404
// page not found error
type PageNotFound404 struct{}

// MethodNotAllowed405 will implement the ServeHTTP func
// so than this struct can used to handle the 405
// method not allowed
type MethodNotAllowed405 struct{}

// ServeHTTP impementation of PageNotFound404
func (e PageNotFound404) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentType, "application/json")
	w.Header().Set(contentSecurityPolicy, "default-src 'self'")
	w.Header().Set(strictTransportSecurity, "max-age=31536000")
	w.WriteHeader(http.StatusNotFound)

	encoder := json.NewEncoder(w)
	encoder.Encode(&Response{Message: MessagePageNotFound, Data: nil})
	log := make(map[string]interface{})
	log["message"] = MessagePageNotFound
	log["status"] = http.StatusNotFound
	log["data"] = nil
	Logger(log)
}

// ServeHTTP impementation of MethodNotAllowed405
func (e MethodNotAllowed405) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentType, "application/json")
	w.Header().Set(contentSecurityPolicy, "default-src 'self'")
	w.Header().Set(strictTransportSecurity, "max-age=31536000")
	w.WriteHeader(http.StatusMethodNotAllowed)

	encoder := json.NewEncoder(w)
	encoder.Encode(&Response{Message: MessageMethodNotAllowed, Data: nil})
}

// Group create the sub router of path
func Group(path string, parent *mux.Router) *mux.Router {
	return parent.PathPrefix(path).Subrouter().StrictSlash(false)
}

// Groupr create new sub router of path
// with default restful router
func Groupr(path string, parent *mux.Router, c Restful) *mux.Router {
	sub := Group(path, parent)

	sub.HandleFunc(index, c.Handler).Methods(get, post, put, delete, patch)
	sub.HandleFunc(subID, c.Handler).Methods(get, put, patch, delete)

	return sub
}

// Groupd handle 3 routes: today, singledate, rangedate
func Groupd(path string, r *mux.Router, handler http.HandlerFunc) {
	singleDate := fmt.Sprintf("%s/{start_date}", path)
	rangeDate := fmt.Sprintf("%s/{start_date}/{end_date}", path)

	// handler for today
	r.HandleFunc(path, handler).Methods(get)

	// handler for single date
	r.HandleFunc(singleDate, handler).Methods(get)

	// handler for range of date
	r.HandleFunc(rangeDate, handler).Methods(get)
}
