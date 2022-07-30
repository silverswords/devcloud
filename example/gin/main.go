package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverswords/devcloud/pkg/framework"
	dcgin "github.com/silverswords/devcloud/pkg/http/gin"
	"github.com/silverswords/devcloud/pkg/workflows"
	"github.com/zeromicro/go-zero/core/conf"
)

func main() {
	appconfig := AppConfig{}
	conf.MustLoad("./config.yaml", &appconfig)
	log.Printf("%+v", appconfig)

	ginctx := dcgin.GetGinCtx()
	log.Printf("%+v", ginctx)

	ginctx.GET("/v1/greeter", Greeter)
	workflows.RunWorker()
	framework.Start()
}

// Greeter handler
// @Summary Greeter
// @Id 1
// @Tags Hello
// @version 1.0
// @Param name query string true "name"
// @produce application/json
// @Success 200 {object} GreeterResponse
// @Router /v1/greeter [get]
func Greeter(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, &GreeterResponse{
		Message: fmt.Sprintf("Hello %s!", ctx.Query("name")),
	})
}

type GreeterResponse struct {
	Message string
}
