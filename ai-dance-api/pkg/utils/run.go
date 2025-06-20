package utils

import (
	"context"
	"time"

	"ai-dance-api/pkg/xerr"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
)

type Response struct {
	Code      string `json:"code"`
	Msg       string `json:"msg"`
	Err       string `json:"err"`
	TraceID   string `json:"req_id"`
	Timestamp int64  `json:"t"`
	Data      any    `json:"data"`
}

func ResSuccess(c context.Context, ctx *app.RequestContext, data any) {
	ctx.JSON(200, Response{
		Code:      "0",
		Msg:       "success",
		Err:       "",
		TraceID:   GetTraceID(c),
		Timestamp: time.Now().UnixMilli(),
		Data:      data,
	})
}

func ResFail(c context.Context, ctx *app.RequestContext, err error) {
	ResFailWithData(c, ctx, err, DataHandler(nil))
}

func ResFailWithData(c context.Context, ctx *app.RequestContext, err error, data any) {
	bizErr := xerr.NewCommBizErr("系统繁忙")
	// 错误类型断言
	if v, ok := err.(*xerr.BizError); ok {
		bizErr.Code = v.Code
		bizErr.Msg = v.Msg
	}

	ctx.JSON(200, Response{
		Code:      bizErr.Code,
		Msg:       bizErr.Msg,
		Err:       err.Error(),
		TraceID:   GetTraceID(c),
		Timestamp: time.Now().UnixMilli(),
		Data:      data,
	})
}

type Empty struct{}

func Run[R, T any](fn func(context.Context, *app.RequestContext, R) (*T, error)) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		var req R
		if err := ctx.BindAndValidate(&req); err != nil {
			return
		}
		res, err := fn(c, ctx, req)
		if err != nil {
			ResFailWithData(c, ctx, err, DataHandler(res))
			return
		}
		ResSuccess(c, ctx, DataHandler(res))
	}
}

func GetTraceID(c context.Context) string {
	value := c.Value("trace_id")
	if value != nil {
		return value.(string)
	}
	return uuid.New().String()
}

func DataHandler(data any) any {
	if data == nil {
		return Empty{}
	}
	return data
}
