package main

import (
	"kubejiangnan/global"
	"kubejiangnan/initialize"
)

/// Project Start

func main() {
	r := initialize.Routers()
	initialize.Viper()
	initialize.K8S()
	panic(r.Run(global.CONF.System.Addr))
}
