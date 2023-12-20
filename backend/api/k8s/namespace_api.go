package k8s

import (
	"context"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubejiangnan/global"
	namespace_reponse "kubejiangnan/model/namespace/response"
	"kubejiangnan/response"
)

type NameSpaceApi struct {
}

func (*NameSpaceApi) GetNamespaceList(c *gin.Context) {

	ctx := context.TODO()
	list, err := global.KubeConfigSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	namespaceList := make([]namespace_reponse.NameSpace, 0)
	for _, item := range list.Items {
		namespaceList = append(namespaceList, namespace_reponse.NameSpace{
			Name:              item.Name,
			CreationTimeStamp: item.CreationTimestamp.Unix(),
			Status:            string(item.Status.Phase),
		})
	}

	response.SuccessWithDetailed(c, "Successfully Get", namespaceList)
}
