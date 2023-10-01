package example

import (
	"github.com/gin-gonic/gin"
	"kubejiangnan/response"
)

type ExampleApi struct {
}

func (*ExampleApi) ExampleTest(c *gin.Context) {
	response.SuccessWithDetailed(c, "Request Data Success", map[string]string{
		"message": "pong",
	})
}
