package handler

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

type (
	// Context context
	Context struct {
		handlers map[string]ContextFunc
		result   interface{}
		Router   *mux.Router
		Writer   http.ResponseWriter
		Request  *http.Request
		Vars     map[string]string // Vars
	}

	// ContextFunc func, this func implement http.Handler
	ContextFunc func(*Context) interface{}
)

// New create new Context
func New(middlewares ...mux.MiddlewareFunc) *Context {
	var h = new(Context)

	h.Router = NewRouter()
	if len(middlewares) > 0 {
		h.Router.Use(middlewares...)
	}

	h.handlers = make(map[string]ContextFunc)
	return h
}

// DB returns database instance
func (c *Context) DB(alias string) (*gorm.DB, error) {
	return GetGormDB(alias)
}

// Use sets middlewares chain within context Router
func (c *Context) Use(middlewares ...mux.MiddlewareFunc) {
	if len(middlewares) > 0 {
		c.Router.Use(middlewares...)
	}
}

// Serve call http.ListenAndServe with default setting
func (c *Context) Serve(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), c.Router)
}

// ServeWith call http.ListenAndServe with default negroni classic middleware
func (c *Context) ServeWith(port int, router http.Handler) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

// HandlerFunc execute request chain
func (f ContextFunc) HandlerFunc(w http.ResponseWriter, r *http.Request) interface{} {
	var ctx Context

	ctx.reset(w, r)
	ctx.Vars = mux.Vars(r)
	ctx.result = f(&ctx)

	switch ctx.result.(type) {
	case error, Error, *Error:
		fmt.Printf("[go-handler] error: {%v}\n", ctx.result)
	default:
		fmt.Println("[go-handler] info: request complete")
	}

	return ctx.result
}

// ServeHTTP calls f(w, r).
func (f ContextFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.HandlerFunc(w, r)
}

func (c *Context) reset(w http.ResponseWriter, r *http.Request) {
	c.setWriter(w)
	c.setRequest(r)
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

	var sub *mux.Router
	sub = Group(path, c.Router)
	sub.Handle(index, c.handlers[method+path]).Methods(indexMethods...)
	sub.Handle(subID, c.handlers[method+path]).Methods(subIDMethods...)
}

// SetRequest set http.Request
func (c *Context) setRequest(r *http.Request) {
	c.Request = r
}

// SetWriter set http.ResponseWriter
func (c *Context) setWriter(w http.ResponseWriter) {
	c.Writer = w
}

// REST map request as http RESTful resource
func (c *Context) REST(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.addRest(restful, path, ctx, middlewares)
}

// SubRouter create sub router and set ctx
func (c *Context) SubRouter(path string, middlewares ...mux.MiddlewareFunc) *Context {
	var (
		subRouter  *mux.Router
		newContext *Context
	)

	subRouter = Group(path, c.Router)
	subRouter.Use(middlewares...)

	newContext = New()
	newContext.Router = subRouter
	return newContext
}

// GET handle http GET request
func (c *Context) GET(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.addRoute(get, path, ctx, middlewares)
}

// POST handle http POST request
func (c *Context) POST(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.addRoute(post, path, ctx, middlewares)
}

// PUT handle http PUT request
func (c *Context) PUT(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.addRoute(put, path, ctx, middlewares)
}

// PATCH handle http PATCH request
func (c *Context) PATCH(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.addRoute(patch, path, ctx, middlewares)
}

// DELETE handle http DELETE request
func (c *Context) DELETE(path string, ctx ContextFunc, middlewares ...mux.MiddlewareFunc) {
	c.addRoute(delete, path, ctx, middlewares)
}

// FormData parse the incoming POST body into "form" struct
// handle application/json and application/x-www-form-urlencoded
// and multipart/form-data
func (c *Context) FormData(form interface{}) error {
	var contentType string
	contentType = c.Request.Header.Get(contentType)

	if strings.Contains(contentType, "application/json") {
		decoder := json.NewDecoder(c.Request.Body)
		decoder.Decode(form)
		return nil
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		c.Request.ParseForm()
		decoder := schema.NewDecoder()
		decoder.Decode(form, c.Request.Form)
		return nil
	} else if strings.Contains(contentType, "multipart/form-data") {
		c.Request.ParseMultipartForm(defaultMaxMemory)
		decoder := schema.NewDecoder()
		decoder.Decode(form, c.Request.Form)
		return nil
	} else {
		custError := &Error{
			Description: invalidContentType,
		}
		Logger(custError)
		return custError
	}
}

// FormFile returns the first file for the provided form key.
// FormFile calls ParseMultipartForm and ParseForm if necessary.
func (c *Context) FormFile(field string) (multipart.File, *multipart.FileHeader, error) {
	return c.Request.FormFile(field)
}

// PostFormValue returns the first value for the named component of the POST,
// PATCH, or PUT request body. URL query parameters are ignored.
// PostFormValue calls ParseMultipartForm and ParseForm if necessary and ignores
// any errors returned by these functions.
// If key is not present, PostFormValue returns the empty string.
func (c *Context) PostFormValue(key string) string {
	if c.Request.PostForm == nil {
		c.Request.ParseMultipartForm(defaultMaxMemory)
	}
	if vs := c.Request.PostForm[key]; len(vs) > 0 {
		return vs[0]
	}
	return ""
}

// PostFormValues returns the all value for the named component of the POST,
// PATCH, or PUT request body. URL query parameters are ignored.
// PostFormValue calls ParseMultipartForm and ParseForm if necessary and ignores
// any errors returned by these functions.
// If key is not present, PostFormValue returns the empty array.
func (c *Context) PostFormValues(key string) []string {
	if c.Request.PostForm == nil {
		c.Request.ParseMultipartForm(defaultMaxMemory)
	}
	return c.Request.PostForm[key]
}

// DecodeURLQuery parse the incoming URL query into struct Urlq.
// Returns true if everyting went well, otherwise false.
func (c *Context) DecodeURLQuery() (urlQuery URLQuery, err error) {
	// get url's query
	decoder := schema.NewDecoder()
	decoder.Decode(&urlQuery, c.Request.URL.Query())
	if err := urlQuery.Validate(); err != nil {
		return urlQuery, err
	}

	return urlQuery, nil
}

// Get implements Get() of RestHandlers
// Returns 405-Not allowed if derivied RestHandlers
// object does not override Get() function
func (c *Context) Get() interface{} {
	return c.MethodNotAllowed()
}

// GetID implements GetID() of RestHandlers
// Returns 405-Not allowed if derivied RestHandlers
// object does not override GetID() function
func (c *Context) GetID(id string) interface{} {
	return c.MethodNotAllowed()
}

// Post implements Post() of RestHandlers
// Returns 405-Not allowed if derivied RestHandlers
// object does not override Post() function
func (c *Context) Post() interface{} {
	return c.MethodNotAllowed()
}

// PutID implements PutID() of RestHandlers
// Returns 405-Not allowed if derivied RestHandlers
// object does not override PutID() function
func (c *Context) PutID(id string) interface{} {
	return c.MethodNotAllowed()
}

// Put implements Put() of RestHandlers
// Returns 405-Not allowed if derivied RestHandlers
// object does not override Put() function
func (c *Context) Put() interface{} {
	return c.MethodNotAllowed()
}

// PatchID implements PatchID() of RestHandlers
// Returns 405-Not allowed if derivied RestHandlers
// object does not override PatchID() function
func (c *Context) PatchID(id string) interface{} {
	return c.MethodNotAllowed()
}

// Patch implements Patch() of RestHandlers
// Returns 405-Not allowed if derivied RestHandlers
// object does not override Patch() function
func (c *Context) Patch() interface{} {
	return c.MethodNotAllowed()
}

// DeleteID implements DeleteID() of RestHandlers
// Returns 405-Not allowed if derivied RestHandlers
// object does not override DeleteID() function
func (c *Context) DeleteID(id string) interface{} {
	return c.MethodNotAllowed()
}

// Delete implements Delete() of RestHandlers
// Returns 405-Not allowed if derivied RestHandlers
// object does not override Delete() function
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

func (c *Context) Write(message string, data interface{}, status int) interface{} {
	return response(c.Writer, message, data, status)
}

// SendImage returns image in response body
func (c *Context) SendImage(path string) interface{} {
	var err error
	if err = WriteImage(path, c.Writer); err != nil {
		return DescError(err)
	}
	return "send image: " + path
}
