package pkg

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	TestServer.Config.QueueLimit = 3

	go func() {
		if err := http.ListenAndServe("127.0.0.1:6070", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Second * 10)
			w.Write([]byte("selam"))
		})); err != nil {
			t.Error(err)
		}
	}()

	go func() {
		time.Sleep(time.Second)
		for i := 0; i < 600; i++ {
			go func(i int) {
				fmt.Printf("index: %d res: %s\n", i, proxyrequest("http://localhost:6070"))

			}(i)

			time.Sleep(time.Second)
		}

	}()

	TestServerF(t)
}

func proxyrequest(_url string) string {

	res, err := http.Get(fmt.Sprintf("http://localhost:8080?url=%s&token=%s", url.PathEscape(_url), TestDb.Users[0].Token))
	if err != nil {
		return fmt.Sprintf("proxyrequest.error: %v", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return string(body)
}
