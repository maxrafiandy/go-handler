# go-handler
This handler use gorilla mux structs as default router.

# How to use
- Set enviroment variable ERROR_IMAGE with value /path/to/error/image.png (default is a blank white rectangle)
- Create log directory in your project directory. The log file will write in *.log

# Example
## Map RESTful routes
```
package main

import (
	"fmt"
	"log"

	"github.com/maxrafiandy/go-handler"
)

type testRest struct {
	handler.Context
}

func (t *testRest) Get() {
	t.Success("Hello, REST Get!")
}

func (t *testRest) Post() {
	t.Success("Hello, REST Post!")
}

func (t *testRest) GetID(id string) {
	t.Success(fmt.Sprintf("Hellon, GET ID %s", id))
}

func (t *testRest) PutID(id string) {
	t.Success(fmt.Sprintf("Hellon, PUT ID %s", id))
}

func main() {
	var (
		goHandler *handler.Context
		rest      testRest
	)

	goHandler = handler.New(handler.JSONfy)

	goHandler.GET("/example-get", func(ctx *handler.Context) {
		ctx.Success("Hello, GET!")
	})

	goHandler.POST("/example-post", func(ctx *handler.Context) {
		ctx.Created("Hello, POST!")
	})

	goHandler.REST("/example-rest", func(ctx *handler.Context) {
		handler.REST(&rest, ctx)
	})

	log.Fatal(goHandler.Serve(8080))
}
```
## Use standard middleware
```
goHandler.REST("/example-rest", func(ctx *handler.Context) {
		handler.REST(&rest, ctx)
}, handler.JSONify)
```