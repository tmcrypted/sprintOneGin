package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterListingRoutes(r *gin.Engine) {
	g := r.Group("/listings")
	_ = g
}
