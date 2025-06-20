package middleware

import (
	"ai-dance-api/pkg/utils"
	"ai-dance-api/pkg/xerr"
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

func LoginAuth() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		uri := ctx.GetRequest().Path()
		// 判断是否需要鉴权
		uriArr := strings.Split(string(uri), "/")
		if len(uriArr) != 5 {
			utils.ResFail(c, ctx, xerr.NewCommBizErr("接口不存在"))
			return
		}
		switch uriArr[1] {
		case "auth-required":
		case "no-auth":
		case "auth-optional":
		default:
			utils.ResFail(c, ctx, xerr.NewCommBizErr("接口不存在"))
			return
		}
	}
}
