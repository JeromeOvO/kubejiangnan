package pod

import "kubejiangnan/convert"

type PodServiceGroup struct {
	PodService
}

var podConvert = convert.ConvertGroupApp.PodConvert
