package router

import (
	"kubejiangnan/router/example"
	"kubejiangnan/router/k8s"
)

type RouterGroup struct {
	ExampleRouterGroup example.ExampleRouter
	K8SRouterGroup     k8s.K8SRouter
}

var RouterGroupApp = new(RouterGroup)
