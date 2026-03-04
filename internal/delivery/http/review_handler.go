package http

import (
	"net/http"
	"strconv"

	"sprin1/internal/delivery/http/dto"
	"sprin1/internal/model"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterReviewRoutes() {
	g := s.router.Group("/reviews")
	g.POST("/create", s.AuthMiddleware, s.createReview)
	g.DELETE("/:id", s.AuthMiddleware, s.deleteReview)
	g.GET("", s.AuthMiddleware, s.getReviews)
}

func (s *Server) createReview(c *gin.Context) {

	var body dto.CreateReviewRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	// Автор всегда берётся из контекста, а не из тела запроса.
	body.AuthorID = user.ID

	review, err := s.reviewService.CreateReview(c.Request.Context(), body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, review)

}

func (s *Server) deleteReview(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := s.reviewService.DeleteReview(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *Server) getReviews(c *gin.Context) { // тут Id передавать в query норм или стрем?
	var query dto.GetReviewsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 20
	}

	if query.PvzID == 0 && query.TargetUserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pvz_id or target_user_id is required"})
		return
	}

	items, total, err := s.reviewService.GetReviews(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"total": total,
		"page":  query.Page,
		"limit": query.Limit,
	})
}
