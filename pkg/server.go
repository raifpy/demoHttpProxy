package pkg

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Logger           logrus.FieldLogger
	ListenAddr       string
	BlacklistedHosts []string
	QueueLimit       int
	WaitQueueTimeout time.Duration
	WaitQueue        bool
}

type Server struct {
	Logger      logrus.FieldLogger
	Database    Database
	Router      *mux.Router
	Client      *http.Client
	UserService *UserService
	Config      ServerConfig
}

// Note that HTTP CONNECT not allowed!
func (s *Server) Handle() {
	s.Logger.Infof("Starting server on %s", s.Config.ListenAddr)
	s.Router.PathPrefix("/")
	s.Router.Use(s.signMiddleware, s.authMiddleware, s.parserMiddleware, s.recordMiddleware, s.queueMiddleware, s.filterHandler, s.requestMiddleware, s.queueFinalizerMiddleware, s.reportMiddleware)

	// Bu kütüphaneden hiç hoşlanmadım.
}

func (s *Server) Listen() error {
	server := &http.Server{
		Addr:     s.Config.ListenAddr,
		Handler:  s.Router,
		ErrorLog: log.Default(),
	}
	return server.ListenAndServe()

}
