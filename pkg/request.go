package pkg

import (
	"context"
	"io"
	"net"
	"net/http"
)

var DefaultHttpClient = &http.Client{
	CheckRedirect: http.DefaultClient.CheckRedirect,
	Transport: &http.Transport{
		Proxy:             nil, // kullanılabilir :))
		ForceAttemptHTTP2: true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			d := &net.Dialer{
				Resolver: &net.Resolver{
					PreferGo:     true,
					StrictErrors: false,

					Dial: func(ctx context.Context, network, address string) (net.Conn, error) {

						return net.Dial("udp", "1.1.1.1:53") // CF DNS
					},
				},
			}
			return d.DialContext(ctx, network, addr)
		},
	},
}

func (s *Server) requestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		request, _ := r.Context().Value(_request).(*http.Request)
		if request == nil {
			s.proxyError(w, r, "unexcepted error occurred", http.StatusBadGateway, "context request missing.")
			return
		}
		s.Logger.Debugf("requesting to %s", request.Host)

		//yeni, _ := http.NewRequestWithContext(r.Context(), "GET", "https://httpbin.org/anything", nil)

		response, err := s.Client.Do(request) // TODO Request with context
		if response != nil && response.Body != nil {
			defer response.Body.Close()
		}
		if err != nil {
			s.proxyError(w, r, err.Error(), http.StatusBadRequest, "client.do:", err.Error())
			return
		}
		s.Logger.Debugf("response %s", response.Status)
		w.WriteHeader(response.StatusCode)
		w.Header().Set("Proxy-Status-Code", "200")
		for key := range w.Header() {
			w.Header().Add(key, w.Header().Get(key))
		}

		size, err := io.Copy(w, response.Body)
		if err != nil {
			s.proxyError(w, r, err.Error(), http.StatusBadGateway, "client.copy:", err.Error())
			return
		}

		s.Logger.Debugf("response %d bytes shared with client", size)

		if next != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), _response, response)))
		}
	})
}
