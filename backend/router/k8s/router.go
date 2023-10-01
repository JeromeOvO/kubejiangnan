package k8s

import (
	"github.com/gin-gonic/gin"
	"kubejiangnan/api"
)

type K8SRouter struct {
}

func (*K8SRouter) InitK8SRouter(r *gin.Engine) {
	group := r.Group("/k8s")
	apiGroup := api.ApiGroupApp.K8SApiGroup
	group.GET("/listPod", apiGroup.GetPodList)
}
