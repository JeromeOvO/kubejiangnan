package initialize

import (
	"github.com/gin-gonic/gin"
	"kubejiangnan/router"
)

func Routers() *gin.Engine {
	r := gin.Default()
	exampleGroup := router.RouterGroupApp.ExampleRouterGroup
	exampleGroup.InitExample(r)
	return r
	// r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
