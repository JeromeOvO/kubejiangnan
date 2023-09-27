package router

import "kubejiangnan/router/example"

type RouterGroup struct {
	ExampleRouterGroup example.ExampleRouter
}

var RouterGroupApp = new(RouterGroup)
