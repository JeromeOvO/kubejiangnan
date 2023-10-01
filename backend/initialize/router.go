package initialize

import (
	"github.com/gin-gonic/gin"
	"kubejiangnan/middleware"
	"kubejiangnan/router"
)

func Routers() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors)
	exampleGroup := router.RouterGroupApp.ExampleRouterGroup
	K8SGroup := router.RouterGroupApp.K8SRouterGroup
	exampleGroup.InitExample(r)
	K8SGroup.InitK8SRouter(r)
	return r
	// r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
