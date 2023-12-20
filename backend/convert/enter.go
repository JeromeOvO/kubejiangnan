package convert

import "kubejiangnan/convert/pod"

type ConvertGroup struct {
	PodConvert pod.PodConvert
}

var ConvertGroupApp = new(ConvertGroup)
