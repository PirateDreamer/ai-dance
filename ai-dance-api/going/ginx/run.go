package ginx

import (
	"context"

	"github.com/gin-gonic/gin"
)

func Run[T, R any](fn func(context.Context, *gin.Context, T) (*R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var req T
		if err := c.ShouldBind(&req); err != nil {
			Fail(ctx, c, err)
			return
		}
		res, err := fn(ctx, c, req)

		resHandler(ctx, c)
		if err != nil {
			Fail(ctx, c, err)
			return
		}
		Success(ctx, c, res)
	}
}
