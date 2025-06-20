package internal

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/spf13/viper"
)

// api规范：api/是否需要鉴权/服务/模块/接口名字

func InitHttp(conf *viper.Viper) {
	h := server.Default()

	// h.Use(middleware.LoginAuth())

	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})

	h.Spin()
}
