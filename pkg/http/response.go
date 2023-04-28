package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	errors3 "mgufrone.dev/job-alerts/pkg/errors"
	"net/http"
)

type DataResponse struct {
	Data interface{} `json:"data"`
}
type DataWithTotalResponse struct {
	DataResponse
	Total uint `json:"total"`
}
type ErrorObj struct {
	Error interface{}
}

func Error(ctx *gin.Context, err error) {
	var er1 errors3.AppError
	if ok := errors.As(err, &er1); ok {
		ctx.AbortWithStatusJSON(er1.Code, err)
		return
	}
	var er2 errors3.ValidationFieldError
	if ok := errors.As(err, &er2); ok {
		ctx.AbortWithStatusJSON(400, gin.H{"error": gin.H{er2.Field: er2.Parent.Error()}})
		return
	}
	if er3, ok := status.FromError(err); ok {
		errCode := 500
		maps := map[codes.Code]int{
			// server side
			codes.Internal:         http.StatusInternalServerError,
			codes.Canceled:         http.StatusInternalServerError,
			codes.DeadlineExceeded: http.StatusGatewayTimeout,
			codes.Unavailable:      http.StatusServiceUnavailable,
			// client side
			codes.NotFound:           404,
			codes.InvalidArgument:    400,
			codes.AlreadyExists:      409,
			codes.FailedPrecondition: http.StatusPreconditionFailed,
			codes.OutOfRange:         400,
			codes.Unauthenticated:    401,
		}
		if v, ok := maps[er3.Code()]; ok {
			errCode = v
		}
		ctx.AbortWithStatusJSON(errCode, gin.H{"error": er3.Message()})
		return
	}
	ctx.AbortWithStatusJSON(500, err)
}
func Unauthorized(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
}
func NotFound(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not found"})
}

func OkWithTotal(ctx *gin.Context, result interface{}, total uint) {
	ctx.JSON(http.StatusOK, &DataWithTotalResponse{
		DataResponse: DataResponse{
			Data: result,
		},
		Total: total,
	})
}

func Ok(ctx *gin.Context, result interface{}) {
	ctx.JSON(http.StatusOK, &DataResponse{Data: result})
}

func Deleted(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, nil)
}
