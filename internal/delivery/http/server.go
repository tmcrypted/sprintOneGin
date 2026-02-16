package http

import "github.com/gin-gonic/gin"

type Server struct {
	router        *gin.Engine
	userService   UserService
	reviewService ReviewService
}

func NewServer(userService UserService, reviewService ReviewService) *Server {
	r := gin.Default()

	s := &Server{
		router:        r,
		userService:   userService,
		reviewService: reviewService,
	}

	s.RegisterRoutes()

	return s
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
