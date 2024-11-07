Ghtml is used to enhance HTML templates in the gin framework.
The gin framework template view is not friendly to layout references, which is also the reason for the birth of ghtml.

* Getting started

```
go get -u github.com/enren-shen/ghtml
```

* Running GHTML

A basic example:

```
package main

import (
	"github.com/enren-shen/ghtml"
	"github.com/gin-gonic/gin"
)
func main() {
	r := gin.New()
	tpl := ghtml.New("./views", ".html")

	//Set default layout
	tpl.Layout("layouts/your-layout.html")

	//Add custom function
	tpl.AddFunc("add", func(a, b int) int {
		return a + b
	})

	//Set real-time loading
	tpl.Reload(true) // or tpl.Reload(false)

	r.HTMLRender = tpl


	//http://localhost:8080/
	r.GET("/", func(ctx *gin.Context) {
		ghtml.Render(ctx, 200, "index.html", nil)
        //or ghtml.Viev(ctx, 200, "index.html", nil)
	})

	r.Run("127.0.0.1:8080")
}
<<<<<<< HEAD

function b(){
    
}
=======
function a(){
}
>>>>>>> bb8cef13aa2a337d45d76e3b0c784787353d4589
