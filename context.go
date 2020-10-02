package handler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var (
	contextPool sync.Pool
)

type (
	context interface {
		ServeHTTP(http.ResponseWriter, *http.Response)
	}

	// Context context
	Context struct {
		handlers map[string]ContextFunc
		result   interface{}
		Router   *mux.Router
		Writer   http.ResponseWriter
		Request  *http.Request
	}

	// ContextFunc func
	ContextFunc func(*Context) interface{}
)

// New create new Context
func New(middlewares ...mux.MiddlewareFunc) *Context {
	var h = new(Context)

	h.Router = NewRouter()
	h.Router.Use(middlewares...)

	h.handlers = make(map[string]ContextFunc)
	contextPool.New = func() interface{} {
		return &Context{
			Router:   nil,
			Writer:   nil,
			Request:  nil,
			handlers: nil,
			result:   nil,
		}
	}
	return h
}

// Serve call http.ListenAndServe with default setting
func (c *Context) Serve(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), c.Router)
}

// HandlerFunc exec request
func (f ContextFunc) HandlerFunc(w http.ResponseWriter, r *http.Request) interface{} {
	ctx := contextPool.Get().(*Context)
	defer contextPool.Put(ctx)

	ctx.Writer = w
	ctx.Request = r
	ctx.result = f(ctx)

	fmt.Printf("Handler result: %v", ctx.result)

	return ctx.result
}

// ServeHTTP calls f(w, r).
func (f ContextFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.HandlerFunc(w, r)
}

func (c *Context) reset(w http.ResponseWriter, r *http.Request) {
	c.Writer = w
	c.Request = r
}

func (c *Context) add(method, path string, ctx ContextFunc, middlewares []mux.MiddlewareFunc) {
	c.Router.Use(middlewares...)
	c.handlers[method+path] = ctx
}

func (c *Context) addRoute(method, path string, ctx ContextFunc, middlewares []mux.MiddlewareFunc) {
	c.add(method, path, ctx, middlewares)
	c.Router.Handle(path, c.handlers[method+path]).Methods(method)
}

func (c *Context) addRest(method, path string, ctx ContextFunc, middlewares []mux.MiddlewareFunc) {
	c.add(method, path, ctx, middlewares)
	sub := Group(path, c.Router)
	sub.Handle(index, c.handlers[method+path]).Methods(get, post, put, delete, patch)
	sub.Handle(subID, c.handlers[method+path]).Methods(get, put, patch, delete)
}

// REST http request
func (c *Context) REST(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.addRest(restful, path, ctx, middlewares)
}

// GET http request
func (c *Context) GET(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.addRoute(get, path, ctx, middlewares)
}

// POST http request
func (c *Context) POST(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.addRoute(post, path, ctx, middlewares)
}

// PUT http request
func (c *Context) PUT(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.addRoute(put, path, ctx, middlewares)
}

// PATCH http request
func (c *Context) PATCH(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.addRoute(patch, path, ctx, middlewares)
}

// DELETE http request
func (c *Context) DELETE(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.addRoute(delete, path, ctx, middlewares)
}

// Get implements Get() function
func (c *Context) Get() interface{} {
	return c.MethodNotAllowed()
}

// GetID implements GetID() function
func (c *Context) GetID(id string) interface{} {
	return c.MethodNotAllowed()
}

// Post implements Post() function
func (c *Context) Post() interface{} {
	return c.MethodNotAllowed()
}

// PutID implements Put() function
func (c *Context) PutID(id string) interface{} {
	return c.MethodNotAllowed()
}

// Put implements Put() function
func (c *Context) Put() interface{} {
	return c.MethodNotAllowed()
}

// PatchID implements PatchID() function
func (c *Context) PatchID(id string) interface{} {
	return c.MethodNotAllowed()
}

// Patch implements Patch() function
func (c *Context) Patch() interface{} {
	return c.MethodNotAllowed()
}

// DeleteID implements DeleteID() function
func (c *Context) DeleteID(id string) interface{} {
	return c.MethodNotAllowed()
}

// Delete implements Delete() function
func (c *Context) Delete() interface{} {
	return c.MethodNotAllowed()
}

// Created send success response with result data
func (c *Context) Created(data interface{}) interface{} {
	return response(c.Writer, MessageCreated, data, http.StatusCreated)
}

// Success send success response with result data
func (c *Context) Success(data interface{}) interface{} {
	return response(c.Writer, MessageOK, data, http.StatusOK)
}

// NoContent send success response without any content
func (c *Context) NoContent() interface{} {
	return response(c.Writer, MessageNoContent, nil, http.StatusNoContent)
}

// BadRequest send general 400-bad request
func (c *Context) BadRequest(data error) interface{} {
	return response(c.Writer, MessageBadRequest, data, http.StatusBadRequest)
}

// NotFound send general 200-Succes withoud data.
// this method is equal to PageNotFound() and usually
// used when record was not found in collection instead of
// return a page not found message
func (c *Context) NotFound() interface{} {
	return response(c.Writer, MessageNotFound, errNotFound, http.StatusOK)
}

// PageNotFound send general 404-not found.
// this method is equal to NotFound() but returns
// page not found message
func (c *Context) PageNotFound() interface{} {
	return response(c.Writer, MessagePageNotFound, errPageNotFound, http.StatusNotFound)
}

// InternalServerError send general 500-interal server error
func (c *Context) InternalServerError(data error) interface{} {
	return response(c.Writer, MessageInternalServerError, data, http.StatusInternalServerError)
}

// Unauthorized send general 401-unautirized
func (c *Context) Unauthorized(data error) interface{} {
	return response(c.Writer, MessageUnauthorized, data, http.StatusUnauthorized)
}

// Forbidden send general 403-forbidden
func (c *Context) Forbidden(data error) interface{} {
	return response(c.Writer, MessageForbidden, data, http.StatusForbidden)
}

// MethodNotAllowed send general 405-Method not allowed
func (c *Context) MethodNotAllowed() interface{} {
	return response(c.Writer, MessageMethodNotAllowed, nil, http.StatusMethodNotAllowed)
}

// NotImplemented send general 405-Method not allowed
func (c *Context) NotImplemented() interface{} {
	return response(c.Writer, MessageNotImplemented, errNotImplemented, http.StatusNotImplemented)
}
