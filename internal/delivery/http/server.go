package http

import "github.com/gin-gonic/gin"

type Server struct {
	router        *gin.Engine
	userService   UserService
	reviewService ReviewService
	authService   AuthService
	jwtSecret     []byte
}

func NewServer(userService UserService, reviewService ReviewService, authService AuthService, jwtSecret string) *Server {
	r := gin.Default()

	s := &Server{
		router:        r,
		userService:   userService,
		reviewService: reviewService,
		authService:   authService,
		jwtSecret:     []byte(jwtSecret),
	}

	s.RegisterRoutes()

	return s
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
