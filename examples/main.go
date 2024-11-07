package main

import (
	"github.com/enren-shen/ghtml"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	tpl := ghtml.New("./views", ".html")

	//Set default layout
	tpl.Layout("layouts/base.html")

	//Add custom function
	tpl.AddFunc("add", func(a, b int) int {
		return a + b
	})

	//Set real-time loading
	tpl.Reload(true) // or tpl.Reload(false)

	r.HTMLRender = tpl

	//1 Use default layout
	//http://localhost:8080/
	r.GET("/", func(ctx *gin.Context) {
		ghtml.Render(ctx, 200, "index.html", nil)
	})

	//2 Use custom layout
	//http://localhost:8080/pink
	r.GET("/pink", func(ctx *gin.Context) {
		ghtml.SetLayout(ctx, "layouts/color.html")
		ghtml.Render(ctx, 200, "color/pink.html", gin.H{
			"title": "pink",
		})
	})

	//3 Do not use layout and call custom functions
	//http://localhost:8080/func
	r.GET("/func", func(ctx *gin.Context) {
		ghtml.SetLayout(ctx, "")
		ghtml.Render(ctx, 200, "color/func.html", nil)
	})

	//4 Group and use the same layout

	user := r.Group("/user", func(ctx *gin.Context) {
		ghtml.SetLayout(ctx, "layouts/user.html")
	})

	{
		//http://localhost:8080/user/create
		user.GET("/create", func(ctx *gin.Context) {
			ghtml.Render(ctx, 200, "user/create.html", nil)
		})

		//http://localhost:8080/user/update
		user.GET("/update", func(ctx *gin.Context) {
			ghtml.Render(ctx, 200, "user/update.html", nil)
		})

		//http://localhost:8080/user/list
		user.GET("/list", func(ctx *gin.Context) {
			ghtml.SetLayout(ctx, "layouts/list.html")
			ghtml.View(ctx, "user/list.html", nil)
		})
	}

	r.Run("127.0.0.1:8080")
}
