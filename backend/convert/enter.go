package convert

import "kubejiangnan/convert/pod"

type ConvertGroup struct {
	PodConvert pod.PodConvertGroup
}

var ConvertGroupApp = new(ConvertGroup)
