package k8s

import (
	"kubejiangnan/convert"
	"kubejiangnan/validate"
)

type ApiGroup struct {
	PodApi
	NameSpaceApi
}

var podValidate = validate.ValidateGroupApp.PodValidate
var podConvert = convert.ConvertGroupApp.PodConvert
