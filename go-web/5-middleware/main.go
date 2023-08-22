package main

import (
  "log"
  "net/http"
  "time"

  "go-study/go-web/5-middleware/gee"
)

func onlyForV2() gee.HandlerFunc {
  return func(c *gee.Context) {
    t := time.Now()
    time.Sleep(time.Second)
    c.Fail(500, "Internal Server Error")
    log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
  }
}

func main() {
  r := gee.New()
  r.Use(gee.Logger()) // global midlleware
  r.GET("/", func(c *gee.Context) {
    c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
  })

  v2 := r.Group("/v2")
  v2.Use(onlyForV2()) // v2 group middleware
  {
    v2.GET("/hello/:name", func(c *gee.Context) {
      c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
    })
  }

  r.Run(":9999")
}

/*
$ go run go-web/5-middleware/main.go
2023/08/22 14:47:06 Route  GET - /
2023/08/22 14:47:06 Route  GET - /v2/hello/:name

$ curl localhost:9999
<h1>Hello Gee</h1>

>>> log
2023/08/22 14:46:53 [200] / in 0s

$ curl localhost:9999/v2/hello/lb
{"message":"Internal Server Error"}

>>> log
2023/08/22 14:47:09 [500] /v2/hello/lb in 1.0097476s for group v2
2023/08/22 14:47:09 [500] /v2/hello/lb in 1.0097476s
*/
