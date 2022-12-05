package pkg

import (
	"net/http"
)

// Custom DNS resolver (DefaultHTTPClient) kullanıldığı zaman ihtiyacımız yok.
func (s *Server) filterHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		request, _ := r.Context().Value(_request).(*http.Request)
		if request == nil {
			s.proxyError(w, r, "unexcepted error occurred", http.StatusBadGateway, "context request missing")
			return
		}
		s.Logger.Debugf("filtering the request %s", request.URL.String())

		if s.checkHostIsBlacklisted(request.URL.Host) {
			s.proxyError(w, r, "request denied", http.StatusBadGateway, "host is blacklisted")
			s.Logger.Debugf("blacklisted request %s aborted", request.URL.String())
			return
		}

		next.ServeHTTP(w, r)
	})

}

func (s *Server) checkHostIsBlacklisted(host string) bool {
	for i := 0; i < len(s.Config.BlacklistedHosts); i++ {
		if host == s.Config.BlacklistedHosts[i] {
			return true
		}
	}
	return false
	// go:for-range var olan listeyi kopyalarak çalışır. Sürekli güncellenen bir liste için kullanımı makul; bu durumda gereksiz bellek kullanımına sebep olacak.

}
