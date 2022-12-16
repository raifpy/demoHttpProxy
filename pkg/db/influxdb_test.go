package db

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/raifpy/proxyApiGateway/pkg"
)

var influxclient influxdb2.Client
var idb InfluxDB

var influxdb = InfluxDB{
	OrgName:           "api_v1",
	UserBucketName:    "user",
	RequestBucketName: "request",
}

func init() {
	rand.Seed(time.Now().Unix())

	idb = InfluxDB{
		Client:            influxdb2.NewClient("http://localhost:8086", "5YWc-oU0wW8SRVqVKCX0o8giMJAyIHPXrSWyjzhIoZg_q6lO4AxAHfznRAIA4MoAlbaSEILo8tM6Zcd2yMgP0w=="),
		OrgName:           "api_v1",
		UserBucketName:    "user",
		RequestBucketName: "request",
	}
	if err := idb.Init(); err != nil {
		panic("influxdb init error: " + err.Error())
	}
}

func TestWriterUser(t *testing.T) {

	setuser := pkg.User{
		Id:    gofakeit.BeerName(),
		Token: uuid.NewString(),
	}

	if err := idb.SetUser(context.Background(), setuser); err != nil {
		t.Fatal(err)
	}

	ureq := pkg.UserRequest{
		UserId:    setuser.Id,
		Ip:        "localhost",
		Method:    "OS",
		Status:    "pending",
		RequestId: uuid.NewString(),
		URL:       "http://localhost:1010",
		BodySize:  1,
		InitTime:  time.Now(),
	}
	fmt.Printf("ureq.RequestId: %v\n", ureq.RequestId)
	fmt.Printf("ureq.InitTime: %v\n", ureq.InitTime)
	fmt.Printf("ureq.InitTime.Unix(): %v\n", ureq.InitTime.Unix())

	if err := idb.SetRequest(context.Background(), ureq); err != nil {
		t.Fatal(err)
	}

	req, err := idb.GetRequest(context.Background(), ureq.RequestId)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("first: %+v\n", req)
	time.Sleep(time.Second)
	req.Status = "done"
	req.UpdateTime = time.Now()
	req.InitTime = time.Now()

	if err := idb.UpdateRequest(context.Background(), req); err != nil {
		t.Fatal(err)
	}

	requpdated, err := idb.GetRequest(context.Background(), ureq.RequestId)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("requpdated: %+v\n", requpdated)
	if requpdated.Status != "done" {
		t.Fatal("Unexpected")
	}

}
