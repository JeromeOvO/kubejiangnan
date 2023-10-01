package global

import (
	"k8s.io/client-go/kubernetes"
	"kubejiangnan/config"
)

var (
	CONF          config.Server
	KubeConfigSet *kubernetes.Clientset
)
