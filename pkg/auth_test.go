package pkg

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

type testResponseWriter struct {
	testing.TB
	io.Writer
	h http.Header
}

func (s testResponseWriter) Header() http.Header {
	return s.h
}

func (s testResponseWriter) WriteHeader(statusCode int) {
	s.h.Set("status", strconv.Itoa(statusCode))
}

func TestAuth(t *testing.T) {
	server := TestServer
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	logger.Out = TestWriter{
		TB: t,
	}
	logger.Formatter = &logrus.JSONFormatter{}
	server.Logger = logger

	testreq := httptest.NewRequest("GET", "http://localhost?token=123", nil)
	testreq.Header.Set("User-Agent", ".")

	//testreq.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(TestDb.Users[0].Token+":burasi_kontrol_edilmiyor")))
	testreq.URL.RawQuery = "token=" + TestDb.Users[0].Token + "&url=https://httpbin.org/anything"
	responsewriter := testResponseWriter{
		TB:     t,
		Writer: &bytes.Buffer{},
		h:      http.Header{},
	}

	server.authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Proxy-Status", "200")
		w.WriteHeader(http.StatusOK)
		spew.Dump(r)

	})).ServeHTTP(responsewriter, testreq)

	if responsewriter.Header().Get("Proxy-Status") != "200" {
		t.Fatalf("unexpected status code: %s", responsewriter.Header().Get("status"))
	}
}
