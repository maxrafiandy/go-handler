# go-handler
This handler use [gorilla mux](https://github.com/gorilla/mux) structs as default router and [gorm](https://github.com/jinzhu/gorm) as default ORM

# How to use
- Run `go get github.com/maxrafiandy/go-handler`
- (Optional) Set enviroment variable `ERROR_IMAGE` with value of */path/to/error/image.png*. Otherwise SendImage(path) returning page not found

# Examples
## Standard Route
```
package main

import (
  "fmt"
  "log"

  "github.com/maxrafiandy/go-handler"
)

func main() {
  var goHandler *handler.Context

  // Create new handler object
  goHandler = handler.New(handler.JSONify, handler.Logging)

  // route to http://localhost:8080/example-get
  // curl -X GET http://localhost:8080/example-get
  goHandler.GET("/example-get", func(ctx *handler.Context) interface{} {
    return ctx.Unauthorized(&handler.Error{Description: "Not authorized"})
  })

  // route to http://localhost:8080/example-post
  // curl -X POST http://localhost:8080/example-post
  goHandler.POST("/example-post", func(ctx *handler.Context) interface{} {
    return ctx.Created("Hello, POST!")
  })

  // route to http://localhost:8080/example-image
  // curl -X GET http://localhost:8080/example-image
  goHandler.GET("/example-get", func(ctx *handler.Context) interface{} {
    return ctx.SendImage("assets/example.png"})
  })

  log.Fatal(goHandler.Serve(8080))
}
```
## Restful Route
```
// create a derivied struc of handler.Context
type testRest struct {
  handler.Context
}

// Get handle http://localhost:8080/example-rest with http method GET
func (t *testRest) Get() interface{} {
  return t.Success("Hello, REST Get!")
}

// Post handle http://localhost:8080/example-rest with http method POST
func (t *testRest) Post() interface{} {
  return t.Success("Hello, REST Post!")
}

// GetID handle http://localhost:8080/example-rest/{id} with http method GET
func (t *testRest) GetID(id string) interface{} {
  return t.Success(fmt.Sprintf("Hello, GET ID %s", id))
}

// Put handle http://localhost:8080/example-rest with http method Put
func (t *testRest) Put() interface{} {
  return t.Success("Hello, REST Post!")
}

// PutID handle http://localhost:8080/example-rest/{id} with http method PUT
func (t *testRest) PutID(id string) interface{} {
  return t.Success(fmt.Sprintf("Hello, PUT ID %s", id))
}

// Patch handle http://localhost:8080/example-rest with http method PATCH
func (t *testRest) Patch() interface{} {
  return t.Success("Hello, REST Patch!")
}

// PatchID handle http://localhost:8080/example-rest/{id} with http method PATCH
func (t *testRest) PatchID(id string) interface{} {
  return t.Success(fmt.Sprintf("Hello, PATCH ID %s", id))
}

// Delete handle http://localhost:8080/example-rest with http method DELETE
func (t *testRest) Delete() interface{} {
  return t.Success("Hello, REST Delete!))
}

// DeleteID handle http://localhost:8080/example-rest/{id} with http method DELETE
func (t *testRest) DeleteID(id string) interface{} {
  return t.Success(fmt.Sprintf("Hello, DELETE ID %s", id))
}

func main() {
  var (
    goHandler *handler.Context
    rest testRest
  )

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
## Adding sub router
```
// create api router with CSP middleware
api := goHandler.SubRouter("/api/v1", handler.AddCSP)

// route to http://localhost:8080/api/v1/example-post
api.POST("/example-post", func(ctx *handler.Context) interface{} {
  // do something
  return ctx.Success("Hello, world!")
})

// route to http://localhost:8080/api/v1/example-get
api.GET("/example-get", func(ctx *handler.Context) interface{} {
  // do something
  return ctx.Success("Hello, world!")
})

// create public router
public := goHandler.SubRouter("/public")

// route to http://localhost:8080/public/example-get
public.GET("/example-get", func(ctx *handler.Context) interface{} {
  // do something
  return ctx.Success("Hello, world!")
})

```
## Getting URL variables
```
goHandler.GET("/example-get/{username}", func(ctx *handler.Context) interface{} {
  result := fmt.Sprintf("Hello, %s!", ctx.Vars["username"])

  return ctx.Success(result)
}, handler.JSONify)
```
## Accessing database
```
// Connect to database (supported driver: mysql, postgres, mssql).
// set DEBUG_MODE in your enviroment variable to enable console query debuging.
// default value of DEBUG_MODE is FALSE
handler.GormAdd("connectionName", handler.NewGormProp("localhost", "3306", "userdb", "password", "dbname", "driver"))

// GormGet returns a *gorm.DB object with preload enabled.
// for more details read official gorm documentation.
db := handler.GormGet("connectionName")
```
# Credits to
All dependencies libraries listed in go.mod file