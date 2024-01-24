package engine

import (
	"encoding/json"
	"fmt"
	"gweb/utils"
	"log"
	"net/http"
)

type H map[string]interface{}

type HandlerFunc func(ctx *Context)

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	Params     map[string]string
	StatusCode int
	handlers   []HandlerFunc
	index      int
	engine     *Engine
}

func EmptyContext() *Context {
	return &Context{}
}

func (ctx *Context) Init(w http.ResponseWriter, req *http.Request) {
	ctx.Path = req.URL.Path
	ctx.Method = req.Method
	ctx.Req = req
	ctx.Writer = w
	ctx.index = -1
	ctx.handlers = nil
	ctx.Params = nil
}

func (ctx *Context) Destroy() {
	ctx.Writer = nil
	ctx.Req = nil
}

func (ctx *Context) Fail(code int, err string) {
	ctx.index = len(ctx.handlers)
	ctx.JSON(code, H{"message": err})
}

func (ctx *Context) Param(key string) string {
	value := ctx.Params[key]
	return value
}

func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key)
}

func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

func (ctx *Context) Status(code int) {
	ctx.StatusCode = code
	ctx.Writer.WriteHeader(code)
}

func (ctx *Context) SetHeader(key string, value string) {
	ctx.Writer.Header().Set(key, value)
}

func (ctx *Context) String(code int, format string, values ...interface{}) {
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.Status(code)
	ctx.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.SetHeader("Content-Type", "application/json")
	ctx.Status(code)
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx.Writer, err.Error(), 500)
	}
}

func (ctx *Context) Data(code int, data []byte) {
	ctx.Status(code)
	ctx.Writer.Write(data)
}

func (ctx *Context) Template(code int, name string, data any) {
	ctx.SetHeader("Content-Type", "text/html")
	ctx.Status(code)
	if err := ctx.engine.htmlTemplates.ExecuteTemplate(ctx.Writer, name, data); err != nil {
		ctx.Fail(500, err.Error())
	}
}

func (ctx *Context) HTML(code int, content string) {
	ctx.SetHeader("Content-Type", "text/html")
	ctx.Status(code)
	ctx.Writer.Write([]byte(content))
}

func (ctx *Context) AddHandlers(f ...HandlerFunc) {
	ctx.handlers = append(ctx.handlers, f...)
}

func (ctx *Context) Next() {
	for ctx.index++; ctx.index < len(ctx.handlers); ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}

func Recovery(ctx *Context) {
	defer func() {
		if err := recover(); err != nil {
			message := fmt.Sprintf("%s", err)
			log.Printf("%s\n\n", utils.Trace(message))
			ctx.Fail(http.StatusInternalServerError, "Internal Server Error")
		}
	}()

	ctx.Next()
}

func NoRoute(ctx *Context) {
	ctx.Fail(404, "Not Found")
}
