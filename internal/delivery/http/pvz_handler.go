package http

import (
	"net/http"
	"strconv"

	"sprin1/internal/delivery/http/dto"
	"sprin1/internal/model"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterPVZRoutes() {
	g := s.router.Group("/pvz")
	g.POST("/create", s.AuthMiddleware, s.createPVZ)
	g.POST("/moderate", s.AuthMiddleware, s.ModeratorMiddleware, s.moderatePVZ)
	g.GET("/:id", s.AuthMiddleware, s.getPVZ)
	g.GET("/all", s.AuthMiddleware, s.ModeratorMiddleware, s.getAllPVZ)
}

func (s *Server) getPVZ(c *gin.Context) {
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

	pvz, err := s.pvzService.GetPVZ(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pvz)
}

func (s *Server) moderatePVZ(c *gin.Context) {
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

	var body dto.ModeratePVZRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pvz, err := s.pvzService.ModeratePVZ(c.Request.Context(), user.ID, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pvz)
}

func (s *Server) createPVZ(c *gin.Context) {
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

	var body dto.CreatePVZRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pvz, err := s.pvzService.CreatePVZ(c.Request.Context(), user.ID, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pvz)
}

func (s *Server) getAllPVZ(c *gin.Context) {
	var query dto.GetAllPVZQuery
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

	items, total, err := s.pvzService.GetAllPVZ(c.Request.Context(), query)
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
