package validate

import (
	"errors"
	pod_req "kubejiangnan/model/pod/request"
)

const (
	IMAGE_PULL_POLICY_IFNOTPRESENT = "IfNotPresent"
)

const (
	RESTART_POLICY_ALWAYS = "Always"
)

type PodValidate struct {
}

func (*PodValidate) Validate(podReq *pod_req.Pod) error {
	// 1. check necessary item
	// 2. init non-default item with default value

	if podReq.Base.Name == "" {
		return errors.New("please define the name of pod")
	}

	if podReq.Base.NameSpace == "" {
		return errors.New("please define the namespace of pod")
	}

	if len(podReq.Containers) == 0 {
		return errors.New("please define the information of pod")
	}

	if len(podReq.InitContainers) > 0 {
		for index, container := range podReq.InitContainers {
			if container.Name == "" {
				return errors.New("existing unnamed container in init-containers")
			}
			if container.Image == "" {
				return errors.New("existing undefined-image container in init-containers")
			}
			if container.ImagePullPolicy == "" {
				podReq.InitContainers[index].ImagePullPolicy = IMAGE_PULL_POLICY_IFNOTPRESENT
			}
		}
	}

	if len(podReq.Containers) > 0 {
		for index, container := range podReq.Containers {
			if container.Name == "" {
				return errors.New("existing unnamed container in init-containers")
			}
			if container.Image == "" {
				return errors.New("existing undefined-image container in init-containers")
			}
			if container.ImagePullPolicy == "" {
				podReq.Containers[index].ImagePullPolicy = IMAGE_PULL_POLICY_IFNOTPRESENT
			}
		}
	}

	if podReq.Base.RestartPolicy == "" {
		podReq.Base.RestartPolicy = RESTART_POLICY_ALWAYS
	}

	return nil
}
