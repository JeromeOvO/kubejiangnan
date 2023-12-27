package k8s

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"kubejiangnan/global"
	pod_req "kubejiangnan/model/pod/request"
	"kubejiangnan/response"
	"net/http"
)

type PodApi struct {
}

// Because the field attributes modified using the update interface are limited,
// and during the actual modification process, any defined field will be modified
func (*PodApi) UpdatePod(ctx context.Context, pod *coreV1.Pod) error {
	// Update
	_, err := global.KubeConfigSet.CoreV1().Pods(pod.Namespace).Update(ctx, pod, metav1.UpdateOptions{})
	return err
}

func (*PodApi) PatchPod(patchData map[string]interface{}, k8sPod *coreV1.Pod, ctx context.Context) error {
	patchDataBytes, _ := json.Marshal(&patchData)
	_, err := global.KubeConfigSet.CoreV1().Pods(k8sPod.Namespace).Patch(
		ctx,
		k8sPod.Name,
		types.StrategicMergePatchType,
		patchDataBytes,
		metav1.PatchOptions{},
	)
	return err
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

	if msg, err := podService.CreateOrUpdatePod(podReq); err != nil {
		response.FailWithMessage(c, msg)
	} else {
		response.SuccessWithMessage(c, msg)
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
