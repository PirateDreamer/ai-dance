package ginx

import (
	"a-dance-api/going/xerr"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Err  string `json:"err"`
	Data any    `json:"data"`
}

func FailWith(ctx context.Context, c *gin.Context, err error, data any) {

	var resErr *xerr.BizErr
	switch e := err.(type) {
	case *xerr.BizErr:
		resErr = e
	default:
		resErr = xerr.NewBizErr("1", "服务不见了")
	}

	// 判断错误类型
	c.JSON(200, Response{
		Code: resErr.Code,
		Msg:  resErr.Msg,
		Err:  err.Error(),
		Data: data,
	})
}

func Fail(ctx context.Context, c *gin.Context, err error) {
	FailWith(ctx, c, err, nil)
}

func Success(ctx context.Context, c *gin.Context, data any) {
	c.JSON(200, Response{
		Code: "0",
		Msg:  "Success",
		Err:  "",
		Data: data,
	})
}

func resHandler(ctx context.Context, c *gin.Context) {
	c.Request.Header.Set("x-timestamp", fmt.Sprintf("%d", time.Now().UnixMilli()))
	c.Request.Header.Set("x-trace-id", "")
}
