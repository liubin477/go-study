package gee

import (
  "log"
  "net/http"
  "path"
  "strings"
  "text/template"
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
    *RouterGroup
    groups []*RouterGroup
    router *router
    // html render
    htmlTemplates *template.Template // 将所有的模板加载进内存
    funcMap       template.FuncMap   // 所有的自定义模板渲染函数
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

// serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
  handler := group.createStaticHandler(relativePath, http.Dir(root))
  urlPattern := path.Join(relativePath, "/*filepath")
  // Register GET handlers
  group.GET(urlPattern, handler)
}

// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
  absolutePath := path.Join(group.prefix, relativePath)
  fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
  return func(c *Context) {
    file := c.Param("filepath")
    // Check if file exists and/or if we have permission to access it
    if _, err := fs.Open(file); err != nil {
      c.Status(http.StatusNotFound)
      return
    }

    fileServer.ServeHTTP(c.Writer, c.Req)
  }
}

// 设置自定义渲染函数
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
  engine.funcMap = funcMap
}

// 加载模板
func (engine *Engine) LoadHTMLGlob(pattern string) {
  engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
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
  c.engine = engine
  engine.router.handle(c)
}
