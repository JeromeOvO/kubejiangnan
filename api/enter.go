package api

import "kubejiangnan/api/example"

type ApiGroup struct {
	ExampleApiGroup example.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
