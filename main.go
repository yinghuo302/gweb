package main

import (
	"gweb/engine"
	"net/http"
)

func main() {
	r := engine.New()
	r.GET("/", func(c *engine.Context) {
		c.HTML(http.StatusOK, "<h1>Hello GWeb</h1>")
	})
	// 获取Query参数，路径参数
	r.GET("/hello", func(c *engine.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	r.GET("/goodbye/:name", func(c *engine.Context) {
		c.JSON(http.StatusOK, engine.H{
			"msg":  "goodbye",
			"name": c.Param("name"),
		})
	})
	r.GET("/file/*filepath", func(c *engine.Context) {
		c.JSON(http.StatusOK, engine.H{
			"msg":  "getFile",
			"path": c.Param("filepath"),
		})
	})
	// 分组控制
	use_rec := r.Group("/panic")
	// 添加中间件，故障恢复
	use_rec.Use(engine.Recovery)
	use_rec.GET("/test", func(ctx *engine.Context) {
		panic("test")
	})
	r.Run(":9999")
}
