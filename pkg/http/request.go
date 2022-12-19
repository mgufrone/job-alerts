package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mgufrone/go-utils/try"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	RequestKey string = "RequestKey"
)

type Request interface {
	Validate(ctx context.Context) error
}

type RequestConstructor func() Request

func MustSanitize(o RequestConstructor) func(ctx *gin.Context) {
	return Sanitize(o, true)
}
func Sanitize(o RequestConstructor, shouldValidate bool) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		out := o()
		err := try.RunOrError(func() error {
			return ctx.BindHeader(out)
		}, func() error {
			return ctx.BindQuery(out)
		}, func() error {
			return ctx.BindUri(out)
		}, func() error {
			if ctx.Request.Body != nil && (ctx.Request.Method == "POST" || ctx.Request.Method == "PUT") {
				return ctx.Bind(out)
			}
			return nil
		}, func() error {
			if shouldValidate {
				return out.Validate(ctx)
			}
			return nil
		})
		if err != nil {
			Error(ctx, status.Errorf(codes.InvalidArgument, err.Error()))
			return
		}
		ctx.Set(RequestKey, out)
	}
}
