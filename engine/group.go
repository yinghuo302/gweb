package engine

import (
	"gweb/utils"
	"net/http"
	"path"
	"strings"
)

type RouterGroup struct {
	handlers []HandlerFunc
	basePath string
	engine   *Engine
}

func (gr *RouterGroup) Group(prefix string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		handlers: gr.combineHandlers(handlers),
		basePath: utils.JoinPaths(gr.basePath, prefix),
		engine:   gr.engine,
	}
}

func (gr *RouterGroup) Use(handlers ...HandlerFunc) {
	gr.handlers = append(gr.handlers, handlers...)
}

func (group *RouterGroup) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finalSize := len(group.handlers) + len(handlers)
	mergedHandlers := make([]HandlerFunc, finalSize)
	copy(mergedHandlers, group.handlers)
	copy(mergedHandlers[len(group.handlers):], handlers)
	return mergedHandlers
}

func (gr *RouterGroup) Handle(method string, pattern string, handlers ...HandlerFunc) {
	pattern = utils.JoinPaths(gr.basePath, pattern)
	handlers = gr.combineHandlers(handlers)
	gr.engine.router.AddRoute(method, pattern, handlers)
}

func (gr *RouterGroup) GET(pattern string, handlers ...HandlerFunc) {
	gr.Handle("GET", pattern, handlers...)
}

func (gr *RouterGroup) HEAD(pattern string, handlers ...HandlerFunc) {
	gr.Handle("HEAD", pattern, handlers...)
}

func (engine *RouterGroup) POST(pattern string, handlers ...HandlerFunc) {
	engine.Handle("POST", pattern, handlers...)
}

func (engine *Engine) Any(pattern string, handlers ...HandlerFunc) {
	engine.Handle("ALL", pattern, handlers...)
}

func (group *RouterGroup) Static(relativePath, root string) {
	group.StaticFS(relativePath, http.Dir(root))
}

func (group *RouterGroup) StaticFS(relativePath string, fs http.FileSystem) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := group.createStaticHandler(relativePath, fs)
	urlPattern := path.Join(relativePath, "/*filepath")

	group.GET(urlPattern, handler)
	group.HEAD(urlPattern, handler)
}

func (gr *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := utils.JoinPaths(gr.basePath, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(c *Context) {

		file := c.Param("filepath")
		f, err := fs.Open(file)
		if err != nil {
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}
		f.Close()

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}
