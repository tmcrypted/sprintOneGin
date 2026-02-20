package http

import (
	"net/http"

	"sprin1/internal/delivery/http/dto"
	"sprin1/internal/model"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterReviewRoutes() {
	g := s.router.Group("/reviews")
	g.POST("/create", s.AuthMiddleware(), s.createReview())
}

func (s *Server) createReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.CreateReviewRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Автор отзыва всегда берётся из аутентифицированного пользователя, а не из тела запроса.
		userAny, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
			return
		}
		user, ok := userAny.(*model.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user type in context"})
			return
		}
		body.AuthorID = user.ID

		review, err := s.reviewService.CreateReview(c.Request.Context(), body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, review)
	}
}
