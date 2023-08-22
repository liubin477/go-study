package main

import (
  "net/http"

  "go-study/go-web/7-panic/gee"
)

func main() {
  r := gee.Default()
  r.GET("/", func(c *gee.Context) {
    c.String(http.StatusOK, "Hello Geektutu\n")
  })
  // index out of range for testing Recovery()
  r.GET("/panic", func(c *gee.Context) {
    names := []string{"geektutu"}
    c.String(http.StatusOK, names[100])
  })

  r.Run(":9999")
}

/*
$ go run go-web/7-panic/main.go
2023/08/22 17:19:45 Route  GET - /
2023/08/22 17:19:45 Route  GET - /panic
$ curl localhost:9999/panic
{"message":"Internal Server Error"}

>>> log
2023/08/22 17:20:00 runtime error: index out of range [100] with length 1
Traceback:
        D:/Develop/Go/src/runtime/panic.go:884
        D:/Develop/Go/src/runtime/panic.go:113
        D:/go-test/main.go:17
        D:/go-test/gee/context.go:39
        D:/go-test/gee/recovery.go:37
        D:/go-test/gee/context.go:39
        D:/go-test/gee/logger.go:12
        D:/go-test/gee/context.go:39
        D:/go-test/gee/router.go:100
        D:/go-test/gee/gee.go:93
        D:/Develop/Go/src/net/http/server.go:2948
        D:/Develop/Go/src/net/http/server.go:1992
        D:/Develop/Go/src/runtime/asm_amd64.s:159
2023/08/22 17:20:00 [500] /panic in 1.5262ms
*/
