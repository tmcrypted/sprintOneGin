package http

import (
	"net/http"

	"sprin1/internal/delivery/http/dto"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterReviewRoutes() {
	g := s.router.Group("/reviews")
	g.POST("/create", s.AuthMiddleware, s.createReview())
}

func (s *Server) createReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.CreateReviewRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		review, err := s.reviewService.CreateReview(c.Request.Context(), body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, review)
	}
}
