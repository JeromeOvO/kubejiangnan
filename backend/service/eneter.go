package service

import "kubejiangnan/service/pod"

type ServiceGroup struct {
	PodServiceGroup pod.PodServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
