package pkg

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
)

func (s *Server) authMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Debugf("incoming request: %s", r.RemoteAddr)
		// HTTP2 ile token doğrulaması tcp Conn üzerinden yapılabilir. Database'i rahatlatacaktır. Ancak kuyruk için yeniden istek atılacak. Kuyruk db ile token db ayrı olursa oldukça faydalı olur.
		proxyAuthHeader := r.Header.Get("Proxy-Authorization")

		s.Logger.Debugf("request proxy authorization header: %s", proxyAuthHeader)

		var token string

		if proxyAuthHeader == "" || !strings.HasPrefix(proxyAuthHeader, "Basic ") {
			s.Logger.Debugf("proxy authorization header missing, token query using")

			if token = r.URL.Query().Get("token"); token == "" {
				s.proxyAuthError(w, r, "invalid proxy auth token")
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), _proxyType, "api"))

			q := r.URL.Query()
			q.Del("token")
			r.URL.RawQuery = q.Encode()
		} else {
			encoded, err := base64.StdEncoding.DecodeString(proxyAuthHeader[6:])
			if err != nil {
				s.proxyAuthError(w, r, "invalid proxy authorization header-base64")
				return
			}

			if token, _, _ = strings.Cut(string(encoded), ":"); token == "" || len(token) > 40 {
				// Sql injection kontrol edilmeli!
				s.proxyAuthError(w, r, "invalid proxy authorization header")
				return
			}
			encoded = nil

			r.Header.Del("Proxy-Authorization")

		}

		s.Logger.Debugf("token: %s", token)

		user, err := s.Database.GetToken(r.Context(), token)
		if err != nil {
			s.proxyError(w, nil, "unauthorized", http.StatusUnauthorized)
			return
		}
		if user.CanNotRequest {
			s.proxyError(w, nil, "access denied", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), _user, user)))

	})

}
