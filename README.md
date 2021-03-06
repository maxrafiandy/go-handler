# go-handler
This handler use [gorilla mux](https://github.com/gorilla/mux) structs as default router and [gorm](https://github.com/jinzhu/gorm) as default ORM

# How to use
- Run `go get github.com/maxrafiandy/go-handler`
- (Optional) Set enviroment variable `ERROR_IMAGE` with value of */path/to/error/image.png*. Otherwise SendImage(path) returning page not found
- (Optional) Set enviroment variable `DEBUG_MODE` with value of *true* or *false* to print gorm query log. Default value is *false*

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
  )

  // all method are accepted in goHandler.REST func, in order to map
  // resource of testRest instance (inherit from handler.Context) to
  // correct methods, we need to call handler.REST  
  goHandler.REST("/example-rest", func(ctx *handler.Context) interface{} {
    // Data race warning!!! Don't passed ctx as REST's resHandler.
	  // to avoid multiple request accessing same *ctx resource
	  // the REST's RestHandler (first argument) must be a new instance to derified *handler.Context
    var rest testRest
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
### Gorm v2
```
// mysql
config := handler.NewGormConfig("user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
handler.ConnectMysql("connectionAlias", config)

// postges
config := handler.NewGormConfig("user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai")
handler.ConnectPostgres("connectionAlias", nil)

// mssql
config := handler.NewGormConfig("sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm")
handler.ConnectMssal("connectionAlias", config)

// return index "connectionAlias" instance
db := handler.GetGormDB("connectionAlias")
```
### Gorm v1
```
// Connect to database (supported driver: mysql, postgres, mssql).
// set DEBUG_MODE in your enviroment variable to enable console query debuging.
// default value of DEBUG_MODE is FALSE. Do some configuration within config var
config := handler.NewGormv1Config("localhost", "3306", "userdb", "password", "dbname", "driver")
handler.AddGormv1("connectionAlias", config)

// GetGormDBv1 returns a *gorm.DB object with preload enabled.
// for more details read official gorm documentation.
db := handler.GetGormDBv1("connectionAlias")
```
### Redis
```
// create new redis option
// do some other configuration within options var
options := handler.NewRedisOptions("localhost", "6379", "", 0)
handler.AddRedis("connectionAlias", options)

// get redis.client instance
rd := handler.GetRedis("connectionAlias")
```
# Credits to
All dependencies libraries listed in go.mod file
