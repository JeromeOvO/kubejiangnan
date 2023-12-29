package pod

import (
	"context"
	"errors"
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"kubejiangnan/global"
	pod_req "kubejiangnan/model/pod/request"
	pod_res "kubejiangnan/model/pod/response"
	"strings"
)

type PodService struct {
}

func (*PodService) DeletePod(namespace string, name string) error {
	ctx := context.TODO()
	return global.KubeConfigSet.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (*PodService) GetPodList(namespace string) (_ []pod_res.PodListItem, err error) {
	ctx := context.TODO()
	list, err := global.KubeConfigSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		errMsg := fmt.Sprintf("Pod[namespace=%s] Getting Pod List Failed, Details: %s", namespace, err.Error())
		err = errors.New(errMsg)
		return
	}

	podList := make([]pod_res.PodListItem, 0)
	for _, item := range list.Items {
		podList = append(podList, podConvert.PodK8s2ItemRes(item))
	}

	return podList, err
}

func (*PodService) GetPodDetail(namespace string, name string) (podReq pod_req.Pod, err error) {
	ctx := context.TODO()
	podApi := global.KubeConfigSet.CoreV1().Pods(namespace)
	k8sGetPod, err := podApi.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		errMsg := fmt.Sprintf("Pod[namespace=%s, name=%s] Getting Pod Detail Failed, Details: %s", namespace, name, err.Error())
		err = errors.New(errMsg)
		return
	}

	// turn k8s pod into pod request
	podReq = podConvert.K8s2ReqConvert.PodK8s2Req(*k8sGetPod)
	return
}

func (*PodService) CreateOrUpdatePod(podReq pod_req.Pod) (msg string, err error) {
	k8sPod := podConvert.Req2K8sConvert.PodReq2K8s(podReq)
	ctx := context.TODO()
	// [No]update [No]patch [Yes]delete+create
	podApi := global.KubeConfigSet.CoreV1().Pods(k8sPod.Namespace)
	if k8sGetPod, err := podApi.Get(ctx, k8sPod.Name, metav1.GetOptions{}); err == nil {
		// validate the parameters of the pod
		k8sPodCopy := *k8sPod
		k8sPodCopy.Name = k8sPod.Name + "-validate"
		_, err := podApi.Create(ctx, &k8sPodCopy, metav1.CreateOptions{
			DryRun: []string{metav1.DryRunAll}, // DryRyun: just check the response, will not actually submit the data to k8s
		})

		if err != nil {
			errMsg := fmt.Sprintf("Pod[namespace=%s, name=%s] Updating Failed, Details: %s", k8sPod.Namespace, k8sPod.Name, err.Error())
			return errMsg, err
		}
		// delete
		err = podApi.Delete(ctx, k8sPod.Name, metav1.DeleteOptions{})
		if err != nil {
			errMsg := fmt.Sprintf("Pod[namespace=%s, name=%s] Updating Failed, Details: %s", k8sPod.Namespace, k8sPod.Name, err.Error())
			return errMsg, err
		}

		// Create a watcher to monitor the status of pod
		var labelSelectors []string
		for k, v := range k8sGetPod.Labels {
			labelSelectors = append(labelSelectors, fmt.Sprintf("%s=%s", k, v))
		}
		// label format: app=test,app2=test2, use labels to select the pod which is needed to be watched
		watcher, err := podApi.Watch(ctx, metav1.ListOptions{
			LabelSelector: strings.Join(labelSelectors, ","),
		})
		if err != nil {
			errMsg := fmt.Sprintf("Pod[namespace=%s, name=%s] Updating Failed, Details: %s", k8sPod.Namespace, k8sPod.Name, err.Error())
			return errMsg, err
		}

		for event := range watcher.ResultChan() {
			k8sPodChan := event.Object.(*coreV1.Pod)

			// Prevent the pod from being deleted before creating a watcher and starting to listen.
			if _, err := podApi.Get(ctx, k8sPod.Name, metav1.GetOptions{}); k8serror.IsNotFound(err) {
				if createdPod, err := podApi.Create(ctx, k8sPod, metav1.CreateOptions{}); err != nil {
					errMsg := fmt.Sprintf("Pod[namespace=%s, name=%s] Updating Failed, Details: %s", k8sPod.Namespace, k8sPod.Name, err.Error())
					return errMsg, err
				} else {
					successMsg := fmt.Sprintf("Pod[namespace=%s, name=%s] Updating Successfully", createdPod.Namespace, createdPod.Name)
					return successMsg, err
				}
			}

			switch event.Type {
			case watch.Deleted:
				if k8sPodChan.Name != k8sPod.Name {
					continue
				}
				// create
				if createdPod, err := podApi.Create(ctx, k8sPod, metav1.CreateOptions{}); err != nil {
					errMsg := fmt.Sprintf("Pod[namespace=%s, name=%s] Updating Failed, Details: %s", k8sPod.Namespace, k8sPod.Name, err.Error())
					return errMsg, err
				} else {
					successMsg := fmt.Sprintf("Pod[namespace=%s, name=%s] Updating Successfully", createdPod.Namespace, createdPod.Name)
					return successMsg, err
				}
			}
		}
		return "", nil
	} else {
		// Create
		if createdPod, err := podApi.Create(ctx, k8sPod, metav1.CreateOptions{}); err != nil {
			errMsg := fmt.Sprintf("Pod[namespace=%s, name=%s] Creaing Failed, Details: %s", k8sPod.Namespace, k8sPod.Name, err.Error())
			return errMsg, err
		} else {
			successMsg := fmt.Sprintf("Pod[namespace=%s, name=%s] Creating Successfully", createdPod.Namespace, createdPod.Name)
			return successMsg, err
		}
	}
}
