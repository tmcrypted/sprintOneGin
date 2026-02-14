package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, userService UserService, reviewService ReviewService) {
	RegisterUserRoutes(r, userService)
	RegisterListingRoutes(r)
	RegisterPVZRoutes(r)
	RegisterReviewRoutes(r, reviewService)
}
