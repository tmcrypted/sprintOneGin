package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterReviewRoutes(r *gin.Engine, reviewService ReviewService) {
	g := r.Group("/reviews")
	g.POST("/create", createReview(reviewService))
}

func createReview(reviewService ReviewService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			DealID       int64   `json:"deal_id" binding:"required"`
			PvzID        int64   `json:"pvz_id" binding:"required"`
			AuthorID     int64   `json:"author_id" binding:"required"`
			TargetUserID int64   `json:"target_user_id" binding:"required"`
			Rating       int     `json:"rating" binding:"required"`
			Body         *string `json:"body"`
		}
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
