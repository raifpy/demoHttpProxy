package pkg

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (s *Server) reportMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(_uuid).(uuid.UUID)
		response := r.Context().Value(_response).(*http.Response)

		s.Logger.Printf("response received successfully %s", id.String())

		bd, _ := strconv.Atoi(response.Header.Get("Content-Length"))

		if err := s.Database.UpdateRequest(r.Context(), UserRequest{
			UpdateTime:          time.Now().Unix(),
			Status:              "ok",
			Error:               "",
			ResponseContentType: response.Header.Get("Content-Type"),
			ResponseStatus:      response.Status,
			ResponseSize:        bd,
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
