package pkg

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) queueMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(_uuid).(uuid.UUID)
		user := r.Context().Value(_user).(User)

		if s.Config.QueueLimit > 0 {
			tcfg, cancel := context.WithTimeout(r.Context(), s.Config.WaitQueueTimeout)
			defer cancel()
			s.Logger.Debugf("Waiting queue for user %s req %s", user.Token, id.String())
			if err := s.UserService.WaitQueue(user, tcfg); err != nil {
				s.proxyError(w, r, "request queue timeout", http.StatusServiceUnavailable, err.Error())
				return
			}
			s.Logger.Debugf("User %s queue is free for request %s", user.Token, id.String())

		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) queueFinalizerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.UserService.RemoveQueue(r.Context().Value(_user).(User))
		next.ServeHTTP(w, r)
	})
}
