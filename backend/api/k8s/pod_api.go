package k8s

import (
	"context"
	"github.com/gin-gonic/gin"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"kubejiangnan/global"
	pod_req "kubejiangnan/model/pod/request"
	"kubejiangnan/response"
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

func (*PodApi) GetPodListOrDetail(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Query("name")
	if name != "" {
		podDetail, err := podService.GetPodDetail(namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "Successfully Get Pod Detail!", podDetail)
	} else {
		podList, err := podService.GetPodList(namespace)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "Successfully Get Pod List!", podList)
	}
}

func (*PodApi) DeletePod(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := podService.DeletePod(namespace, name)
	if err != nil {
		response.FailWithMessage(c, "Delete Pod Failed, Detail: "+err.Error())
		return
	}
	response.SuccessWithMessage(c, "Successfully Delete Pod")
}
