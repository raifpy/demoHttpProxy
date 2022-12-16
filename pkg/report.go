package pkg

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (s *Server) reportMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(_uuid).(uuid.UUID)
		response := r.Context().Value(_response).(*http.Response)

		s.Logger.Printf("response received successfully %s", id.String())

		if err := s.Database.UpdateRequest(context.Background(), UserRequest{
			UpdateTime:          time.Now(),
			Status:              "ok",
			Error:               "",
			ResponseContentType: response.Header.Get("Content-Type"),
			ResponseStatus:      response.Status,
			ResponseSize:        response.ContentLength,
			RequestId:           id.String(),
		}); err != nil {
			s.Logger.Errorf("failed to update request %s: %s", id.String(), err.Error())
		}

		if next != nil {
			next.ServeHTTP(w, r)
		}

		// TODO: queue -
	})
}
