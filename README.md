# go-handler
This handler use gorilla mux structs as default router.

# How to use
- Set enviroment variable ERROR_IMAGE with value /path/to/error/image.png (default is a blank white rectangle)
- Create log directory in your project directory. The log file will write in *.log

# Example
## Map RESTful routes
```
func Route() *mux.Router{
    // create new mux.Router
    r := handler.NewRouter()

    // create route to /example 
    // this route can now accesses with some http verbs: 
    // GET http://localhost:<port>/example
    // POST http://localhost:<port>/example
    // PUT http://localhost:<port>/example
    // PATCH http://localhost:<port>/example
    // DELETE http://localhost:<port>/example
    // GET http://localhost:<port>/example/{id}
    // PUT http://localhost:<port>/example/{id}
    // PATCH http://localhost:<port>/example/{id}
    // DELETE http://localhost:<port>/example/{id}
    sub := handler.Groupr("/example", &example{})

    return r
}

type example struct {
    handler.RestContext
}

func (rest example) Handler(w http.ResponseWriter, r *http.Request) {
	rest.Writer = w
	rest.Request = r
	handler.MapRest(rest, w, r)
}

func (rest example) Get() {
    // go get stuff here
}

func (rest example) Post() {
    // go post stuff here
}

func (rest example) GetID(id string) {
    // go get /{id} stuff here
}
```
## Use standard middleware
```
func Route() *mux.Router{
    // create new mux.Router
    r := handler.NewRouter()

    // user jsonfy middleware to decode to json
    r.Use(handler.JSONfy)

    // create route to /example 
    sub := handler.Groupr("/example", &example{})

    return r
}
```