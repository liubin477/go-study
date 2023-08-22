package gee

import (
  "log"
  "net/http"
  "strings"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type (
  RouterGroup struct {
    prefix      string        // 前缀，如 / 或 /api
    parent      *RouterGroup  // 父组（支持分组嵌套）
    middlewares []HandlerFunc // 应用在该分组上的中间件
    engine      *Engine       // Group 需要映射路由规则，需要有访问 Router 等能力
  }

  Engine struct {
    *RouterGroup                // 将 Engine 作为最顶层的分组
    groups       []*RouterGroup // store all groups
    router       *router
  }
)

// New is the constructor of gee.Engine
func New() *Engine {
  engine := &Engine{router: newRouter()}
  engine.RouterGroup = &RouterGroup{engine: engine}
  engine.groups = []*RouterGroup{engine.RouterGroup}
  return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
  engine := group.engine
  newGroup := &RouterGroup{
    prefix: group.prefix + prefix,
    parent: group,
    engine: engine,
  }
  engine.groups = append(engine.groups, newGroup)
  return newGroup
}

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
  group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
  pattern := group.prefix + comp
  log.Printf("Route %4s - %s", method, pattern)
  group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
  group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
  group.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
  return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  // 通过 URL 前缀判断请求需要用到哪些中间件
  var middlewares []HandlerFunc
  for _, group := range engine.groups {
    if strings.HasPrefix(req.URL.Path, group.prefix) {
      middlewares = append(middlewares, group.middlewares...)
    }
  }
  c := newContext(w, req)
  c.handlers = middlewares
  engine.router.handle(c)
}
