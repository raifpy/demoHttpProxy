package pkg

import (
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type TestWriter struct {
	testing.TB
}

func (t TestWriter) Write(b []byte) (int, error) {
	t.Logf("%s", b)
	return len(b), nil
}

var TestServer = &Server{
	Logger:   logrus.New(),
	Database: TestDb,
	Router:   mux.NewRouter(),

	Config: ServerConfig{

		ListenAddr: "127.0.0.1:8080",

		BlacklistedHosts: []string{"scrape.do", "localhost"}, // eeeh: Docker üzerinde koşacağı için büyük oranda sınırlanacak
		WaitQueueTimeout: time.Second * 5,
		QueueLimit:       2,
		WaitQueue:        true,
	},
	Client: DefaultHttpClient, // Cloudflare dns kullanıyor
}

func init() {
	TestServer.UserService = NewUserService(TestServer) // mmh?
}

func TestServerF(t *testing.T) {

	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	TestServer.Logger = l
	TestServer.Handle()
	t.Fatal(TestServer.Listen())
}
