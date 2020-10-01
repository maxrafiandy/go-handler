package handler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var contextPool sync.Pool
var contextRestPool sync.Pool

// New create new Context
func New(middlewares ...mux.MiddlewareFunc) *Context {
	var h = new(Context)

	h.Router = NewRouter()
	h.Router.Use(middlewares...)

	h.handlers = make(map[string]ContextFunc)
	contextPool.New = func() interface{} {
		return &Context{
			Writer:   nil,
			Request:  nil,
			handlers: nil,
		}
	}
	return h
}

func (c *Context) ListenAndServe(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), c.Router)
}

// ContextFunc func
type ContextFunc func(*Context)

type context interface {
	// SetRequest(*http.Request)
	// SetWriter(http.ResponseWriter)
	ServeHTTP(http.ResponseWriter, *http.Response)
}

// Context context
type Context struct {
	handlers map[string]ContextFunc
	Router   *mux.Router
	Writer   http.ResponseWriter
	Request  *http.Request
}

// HandlerFunc exec request
func (f ContextFunc) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	ctx := contextPool.Get().(*Context)
	defer contextPool.Put(ctx)

	ctx.Writer = w
	ctx.Request = r
	f(ctx)
}

// ServeHTTP calls f(w, r).
func (f ContextFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.HandlerFunc(w, r)
}

func (c *Context) add(method, path string, ctx ContextFunc, middlewares []mux.MiddlewareFunc) {
	sub := c.Router

	for _, middleware := range middlewares {
		sub.Use(middleware)
	}

	c.handlers[method+path] = ctx
	sub.Handle(path, c.handlers[method+path]).Methods(method)
}

func (c *Context) rest(method, path string, rest RestHandlers, ctx ContextFunc, middlewares []mux.MiddlewareFunc) {
	sub := c.Router

	for _, middleware := range middlewares {
		sub.Use(middleware)
	}

	c.handlers[method+path] = ctx
	sub.Handle(path, c.handlers[method+path])
}

// REST http request
func (c *Context) REST(path string, rest RestHandlers, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.rest("rest", path, rest, ctx, middlewares)
}

// GET http request
func (c *Context) GET(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.add(get, path, ctx, middlewares)
}

// POST http request
func (c *Context) POST(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.add(post, path, ctx, middlewares)
}

// PUT http request
func (c *Context) PUT(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.add(put, path, ctx, middlewares)
}

// PATCH http request
func (c *Context) PATCH(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.add(patch, path, ctx, middlewares)
}

// DELETE http request
func (c *Context) DELETE(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.add(delete, path, ctx, middlewares)
}

// Success send success response with result data
func (c *Context) Success(data interface{}) {
	response(c.Writer, MessageOK, data, http.StatusOK)
}

// Created send success response with result data
func (c *Context) Created(data interface{}) {
	response(c.Writer, MessageCreated, data, http.StatusCreated)
}

// InternalServerError send general 500-interal server error
func (c *Context) InternalServerError(data error) {
	response(c.Writer, MessageInternalServerError, data, http.StatusInternalServerError)
}

// BadRequest send general 400-bad request
func (c *Context) BadRequest(data error) {
	response(c.Writer, MessageBadRequest, data, http.StatusBadRequest)
}

// Unauthorized send general 401-unautirized
func (c *Context) Unauthorized(data error) {
	response(c.Writer, MessageUnauthorized, nil, http.StatusUnauthorized)
}

// Forbidden send general 405-forbidden
func (c *Context) Forbidden(data error) {
	response(c.Writer, MessageForbidden, data, http.StatusForbidden)
}
