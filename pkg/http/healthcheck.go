package http

import (
	"github.com/gin-gonic/gin"
	"time"
)

func healthcheck(ctx *gin.Context) {
	Ok(ctx, gin.H{"status": "ok", "time": time.Now()})
}
func Healthcheck(rg *gin.RouterGroup) {
	rg.GET("/", healthcheck)
	rg.GET("/healthcheck", healthcheck)
}
