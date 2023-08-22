package main

import (
  "fmt"
  "go-study/go-web/6-template/gee"
  "html/template"
  "net/http"
  "time"
)

type student struct {
  Name string
  Age  int8
}

func FormatAsDate(t time.Time) string {
  year, month, day := t.Date()
  return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
  r := gee.New()
  r.Use(gee.Logger())
  r.SetFuncMap(template.FuncMap{
    "FormatAsDate": FormatAsDate,
  })
  r.LoadHTMLGlob("D:/test/templates/*")
  r.Static("/assets", "D:/test/static")

  stu1 := &student{Name: "LB", Age: 20}
  stu2 := &student{Name: "Jack", Age: 22}

  r.GET("/", func(c *gee.Context) {
    c.HTML(http.StatusOK, "css.tmpl", nil)
  })
  r.GET("/students", func(c *gee.Context) {
    c.HTML(http.StatusOK, "arr.tmpl", gee.H{
      "title":  "LB",
      "stuArr": [2]*student{stu1, stu2},
    })
  })
  r.GET("/date", func(c *gee.Context) {
    c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
      "title": "LB",
      "now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
    })
  })
  r.Run(":9999")
}

/*
$ go run go-web/6-template/main.go
2023/08/22 16:43:13 Route  GET - /assets/*filepath
2023/08/22 16:43:13 Route  GET - /
2023/08/22 16:43:13 Route  GET - /students
2023/08/22 16:43:13 Route  GET - /date
2023/08/22 16:43:27 [200] / in 0s
2023/08/22 16:43:27 [0] /assets/css/main.css in 113.0595ms
2023/08/22 16:43:38 [200] /date in 128.9µs
2023/08/22 16:43:54 [200] /students in 221.2µs
*/
