package http

import (
	"net/http"

	"sprin1/internal/delivery/http/dto"

	"github.com/gin-gonic/gin"
)

func RegisterReviewRoutes(r *gin.Engine, reviewService ReviewService) {
	g := r.Group("/reviews")
	g.POST("/create", createReview(reviewService))
}

func createReview(reviewService ReviewService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.CreateReviewRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		review, err := reviewService.CreateReview(c.Request.Context(),
			body.DealID, body.PvzID, body.AuthorID, body.TargetUserID, body.Rating, body.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, review)
	}
}
