package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterPVZRoutes(r *gin.Engine) {
	g := r.Group("/pvz")
	_ = g
}
