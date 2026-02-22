package http

func (s *Server) RegisterRoutes() {
	s.RegisterUserRoutes()
	s.RegisterReviewRoutes()
	s.RegisterAuthRoutes()
	// s.RegisterListingRoutes()
	// s.RegisterPVZRoutes()

}
