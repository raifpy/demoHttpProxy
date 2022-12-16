package pkg

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (s *Server) proxyError(w http.ResponseWriter, r *http.Request, status string, statusCode int, desc ...string) {
	s.Logger.Debugf("proxy error: %s %v", status, desc)
	w.Header().Set("Proxy-Status-Code", strconv.Itoa(statusCode))
	if r != nil {
		id := r.Context().Value(_uuid).(uuid.UUID)

		go s.Database.UpdateRequest(context.Background(), UserRequest{
			UpdateTime: time.Now(),
			RequestId:  id.String(),
			Status:     "error",
			Error:      fmt.Sprintf("%s:%v", status, desc),
		})
	}
	
	http.Error(w, status, statusCode)
}

func (s *Server) proxyAuthError(w http.ResponseWriter, r *http.Request, reason string) {
	s.Logger.Debugf("request aborted: %s %s", r.RemoteAddr, reason)
	w.Header().Set("Proxy-Authenticate:", "Basic")
	s.proxyError(w, nil, "proxy auth required", http.StatusProxyAuthRequired)

	// if s.Config.AbortOnProxyAuthRequired {
	// 	if h, _ := w.(http.Hijacker); h != nil {
	// 		if conn, _, _ := h.Hijack(); conn != nil {
	// 			conn.Close()
	// 		}
	// 	}
	// } // Useless
}
