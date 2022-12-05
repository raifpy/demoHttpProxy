package pkg

import (
	"context"
	"net/http"
	"net/url"
)

func (s *Server) parserMiddleware(next http.Handler) http.Handler {

	// marshaller kullanılabilir

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// var querybased bool
		//  {
		// 	querybased = true
		// }

		//req := r.Clone(r.Context())

		req, err := http.NewRequestWithContext(r.Context(), r.Method, r.URL.String(), r.Body)
		if err != nil {
			s.proxyError(w, r, "unexpected error occurred", http.StatusBadGateway)
			return
		}
		req.Header = r.Header
		query := r.URL.Query()

		if pt, _ := r.Context().Value(_proxyType).(string); pt == "api" {
			if req.URL, err = url.Parse(query.Get("url")); err != nil {
				s.proxyError(w, r, "given wrong or empty url", http.StatusBadRequest, err.Error())
				return
			}
			query.Del("url")
			req.Host = req.URL.Host
			req.URL.RawQuery = query.Encode()
		} else {
			req.URL = r.URL
			req.Host = r.Host
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), _request, req)))
	})
}
