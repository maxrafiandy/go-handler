# go-handler
This handler use [gorilla mux](https://github.com/gorilla/mux) structs as default router and [gorm](https://github.com/jinzhu/gorm) as default ORM

# How to use
- Run `go get github.com/maxrafiandy/go-handler`
- Set enviroment variable `ERROR_IMAGE` with value */path/to/error/image.png* (default is a blank white rectangle)
- Create log directory in your project directory. The log file will be written in \*.log (if log directory is not exist it will created automatically by logger func)

# Example
## Routing
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

func (t *testRest) Get() interface{} {
  return t.Success("Hello, REST Get!")
}

func (t *testRest) Post() interface{} {
  return t.Success("Hello, REST Post!")
}

func (t *testRest) GetID(id string) interface{} {
  return t.Success(fmt.Sprintf("Hellon, GET ID %s", id))
}

func (t *testRest) PutID(id string) interface{} {
  return t.Success(fmt.Sprintf("Hellon, PUT ID %s", id))
}

func main() {
  var (
    goHandler *handler.Context
    rest      testRest
  )

  // Create new handler object
  goHandler = handler.New(handler.JSONify, handler.Logging)

  // route to http://host:port/example-get
  // curl -X GET http://host:port/example-get
  goHandler.GET("/example-get", func(ctx *handler.Context) interface{} {
    return ctx.Unauthorized(&handler.Error{Description: "Not authorized"})
  })

  // route to http://host:port/example-post
  // curl -X POST http://host:port/example-post
  goHandler.POST("/example-post", func(ctx *handler.Context) interface{} {
    return ctx.Created("Hello, POST!")
  })

  // route to http://host:port/example-rest with multiple http methods
  // curl -X GET http://host:port/example-rest
  // curl -X POST http://host:port/example-rest
  // curl -X PUT http://host:port/example-rest
  // curl -X PATCH http://host:port/example-rest
  // curl -X DELETE http://host:port/example-rest
  // curl -X GET http://host:port/example-rest/{id}
  // curl -X PUT http://host:port/example-rest/{id}
  // curl -X PATCH http://host:port/example-rest/{id}
  // curl -X DELETE http://host:port/example-rest/{id}
  // if you don't override/implement handler.RestfulHandlers interface
  // it will simply return method not allowed
  goHandler.REST("/example-rest", func(ctx *handler.Context) interface{} {
    return handler.REST(&rest, ctx)
  })

  log.Fatal(goHandler.Serve(8080))
}
```
## Use standard middleware
```
// add JSONify middleware to all http verbs of /example-rest
goHandler.REST("/example-rest", func(ctx *handler.Context) interface{} {
  return handler.REST(&rest, ctx)
}, handler.JSONify)
```
## Accessing database
```
// Connect to database (supported driver: mysql, postgres, mssql).
// set DEBUG_MODE in your enviroment variable to enable console debuging of gorm itself. 
handler.GormAdd("connectionName", handler.NewGormProp("localhost", "3306", "userdb", "password", "dbname", "driver"))

// GormGet returns a *gorm.DB object with preload enabled.
// for more details read official gorm documentation.
db := handler.GormGet("connectionName")
```
