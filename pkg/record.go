package pkg

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (s *Server) recordMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(_uuid).(uuid.UUID)
		user := r.Context().Value(_user).(User)
		request := r.Context().Value(_request).(*http.Request)

		if err := s.Database.SetRequest(r.Context(), UserRequest{
			URL:       request.URL.String(),
			Ip:        r.RemoteAddr,
			Method:    request.Method,
			Status:    "pending",
			RequestId: id.String(),
			UserId:    user.Id,
			BodySize:  r.ContentLength,
			InitTime:  time.Now(),
		}); err != nil {
			s.Logger.Errorf("request save failed: %v %s", err, id.String())
			s.proxyError(w, r, "server error please try again", http.StatusInternalServerError, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}
