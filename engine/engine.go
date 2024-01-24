package engine

import (
	"gweb/router"
	"net/http"
	"sync"
	"text/template"
)

type Engine struct {
	RouterGroup
	router        router.IRouter
	pool          sync.Pool
	noRoute       []HandlerFunc
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

func New() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			basePath: "/",
			engine:   nil,
			handlers: make([]HandlerFunc, 0),
		},
		router: router.NewTrieRouter(),
		pool: sync.Pool{
			New: func() any {
				return &Context{}
			},
		},
		noRoute: []HandlerFunc{NoRoute},
	}
	engine.engine = engine
	return engine
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := engine.pool.Get().(*Context)
	ctx.Init(w, req)
	params, handlers := engine.router.GetHandlers(ctx.Method, ctx.Path)
	ctx.Params = params
	if handlers == nil {
		ctx.handlers = engine.noRoute
	} else {
		ctx.handlers = handlers.([]HandlerFunc)
	}
	ctx.Next()
	ctx.Destroy()
	engine.pool.Put(ctx)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

func (engine *Engine) SetNoRoute(noRoute ...HandlerFunc) {
	engine.noRoute = noRoute
}
