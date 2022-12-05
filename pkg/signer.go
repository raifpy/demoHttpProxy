package pkg

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) signMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid := uuid.New()
		//s.Logger.Debugf("new inner request: %s %s", r.RemoteAddr, uuid.String())
		w.Header().Set("Proxy-Request-Id", uuid.String())

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), _uuid, uuid)))
	})
}
