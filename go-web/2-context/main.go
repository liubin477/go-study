package main

import (
  "net/http"

  "go-study/go-web/2-context/gee"
)

func main() {
  r := gee.New()
  r.GET("/", func(c *gee.Context) {
    c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
  })
  r.GET("/hello", func(c *gee.Context) {
    c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
  })
  r.POST("/login", func(c *gee.Context) {
    c.JSON(http.StatusOK, gee.H{
      "username": c.PostForm("username"),
      "password": c.PostForm("password"),
    })
  })

  r.Run(":9999")
}

/*
$ go run go-web/2-context/main.go
2023/08/21 15:59:34 Route  GET - /
2023/08/21 15:59:34 Route  GET - /hello
2023/08/21 15:59:34 Route POST - /login
$ curl -i localhost:9999
HTTP/1.1 200 OK
Content-Type: text/html
Date: Mon, 21 Aug 2023 08:03:14 GMT
Content-Length: 18
<h1>Hello Gee</h1>
$ curl localhost:9999/hello?name=lb
hello lb, you're at /hello
$ curl localhost:9999/login -X POST -d "username=lb&password=123"
{"password":"123","username":"lb"}
*/
