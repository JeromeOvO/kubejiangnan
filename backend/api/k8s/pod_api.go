package k8s

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubejiangnan/global"
	pod_req "kubejiangnan/model/pod/request"
	"kubejiangnan/response"
	"net/http"
)

type PodApi struct {
}

func (*PodApi) CreateOrUpdatePod(c *gin.Context) {
	var podReq pod_req.Pod

	if err := c.ShouldBind(&podReq); err != nil {
		response.FailWithMessage(c, "Parameter Parsing Failed, Details: "+err.Error())
		return
	}

	if err := podValidate.Validate(&podReq); err != nil {
		response.FailWithMessage(c, "Parameter Checking Failed, Details: "+err.Error())
		return
	}

	k8sPod := podConvert.PodReq2K8s(podReq)
	ctx := context.TODO()
	if createdPod, err := global.KubeConfigSet.CoreV1().Pods(k8sPod.Namespace).Create(ctx, k8sPod, metav1.CreateOptions{}); err != nil {
		errMsg := fmt.Sprintf("Pod[%s-%s] Creating Failed, Details: %s", k8sPod.Namespace, k8sPod.Name, err.Error())
		response.FailWithMessage(c, errMsg)
		return
	} else {
		successMsg := fmt.Sprintf("Pod[%s-%s] Creating Successfully", createdPod.Namespace, createdPod.Name)
		response.SuccessWithMessage(c, successMsg)
	}
}

func (*PodApi) GetPodList(c *gin.Context) {
	ctx := context.TODO()
	list, err := global.KubeConfigSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{})

	if err != nil {
		fmt.Print(err.Error())
	}

	for _, item := range list.Items {
		fmt.Println(item.Namespace, item.Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
