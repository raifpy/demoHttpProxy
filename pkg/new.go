package pkg

import "github.com/gorilla/mux"

func New(c ServerConfig, db Database) (s *Server) {
	s = &Server{
		Logger:   c.Logger,
		Database: db,
		Router:   mux.NewRouter(),
		Client:   c.Client,
		//UserService: NewUserService(s),
		Config: c,
	}
	s.UserService = NewUserService(s)
	s.Handle()
	return

}
