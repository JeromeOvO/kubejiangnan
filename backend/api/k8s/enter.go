package k8s

import (
	"kubejiangnan/service"
	"kubejiangnan/validate"
)

type ApiGroup struct {
	PodApi
	NameSpaceApi
}

var podValidate = validate.ValidateGroupApp.PodValidate
var podService = service.ServiceGroupApp.PodServiceGroup.PodService
