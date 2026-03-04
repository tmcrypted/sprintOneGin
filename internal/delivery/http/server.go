package http

import "github.com/gin-gonic/gin"

type Server struct {
	router        *gin.Engine
	userService   UserService
	reviewService ReviewService
	authService   AuthService
	pvzService    PVZService
	jwtSecret     []byte
}

func NewServer(userService UserService, reviewService ReviewService, authService AuthService, pvzService PVZService, jwtSecret string) *Server {
	r := gin.Default()

	s := &Server{
		router:        r,
		userService:   userService,
		reviewService: reviewService,
		authService:   authService,
		pvzService:    pvzService,
		jwtSecret:     []byte(jwtSecret),
	}

	s.RegisterRoutes()

	return s
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
