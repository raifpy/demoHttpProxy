package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jessevdk/go-flags"

	"github.com/raifpy/proxyApiGateway/pkg"
	"github.com/raifpy/proxyApiGateway/pkg/db"
	"github.com/sirupsen/logrus"
)

var demoFlags struct {
	Verbose                   bool     `long:"verbose" description:"Enable verbose output"`
	DockerSocket              string   `long:"docker-socket" short:"d" description:"Docker unix socket to find container local ip"`
	Influxdb2Address          string   `long:"influxdb2-address" description:"Adress to access influxdb2 [REQUIRED]" short:"i" required:"true"`
	Infludb2Token             string   `long:"influxdb2-token" short:"t" description:"Token to access influxdb2 [REQUIRED]" required:"true"`
	InfluxdbOrg               string   `long:"influxdb2-org" short:"o" description:"OrgName for influxdb2 [REQUIRED]" required:"true"`
	InfluxdbUserBucketName    string   `long:"influxdb2-user-bucket" description:"Bucket name for users" default:"user"`
	InfluxdbRequestBucketName string   `long:"influxdb2-request-bucket" description:"Bucket name for requests" default:"request"`
	BlackList                 []string `long:"black-list" short:"b" description:"Blacklist to block host request" default:"localhost"`
	ListenAddr                string   `long:"listen-addr" short:"l" description:"Listen port" default:"0.0.0.0:8080" required:"true"`
	EnableQueue               bool     `long:"enable-queue" description:"Enable Queue" short:"q"`
	QueueLimit                int      `long:"queue-limit" description:"Queue for requests of the same token" default:"10" required:"true"`
	WaitQueueTimeout          int      `long:"queue-wait-timeout" description:"Queue wait timeout to drop waiting request. MS" default:"5000" required:"true"`
}

func main() {

	flag := flags.NewParser(&demoFlags, flags.Default|flags.IgnoreUnknown)
	flag.Name = os.Args[0]
	extra, err := flag.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n%s", err, flag.Usage)
		os.Exit(1)
	}

	logger := logrus.New()
	if demoFlags.Verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	cfg := pkg.ServerConfig{
		ListenAddr:       demoFlags.ListenAddr,
		BlacklistedHosts: demoFlags.BlackList,
		QueueLimit:       demoFlags.QueueLimit,
		WaitQueueTimeout: time.Duration(demoFlags.WaitQueueTimeout) * time.Millisecond,
		Logger:           logger,

		WaitQueue: demoFlags.EnableQueue,
		Client:    pkg.DefaultHttpClient,
	}

	if cfg.WaitQueue && cfg.QueueLimit <= 0 || cfg.WaitQueueTimeout <= 0 {
		fmt.Fprintf(os.Stderr, "error: wrong queue configurations\n%s", flag.Usage)
		os.Exit(1)
	}

	if demoFlags.DockerSocket != "" {
		if demoFlags.DockerSocket[0] == '/' {
			demoFlags.DockerSocket = "unix://" + demoFlags.DockerSocket
		}

		logger.Debugf("Connecting to the docker")
		client, err := client.NewClientWithOpts(client.WithHost(demoFlags.DockerSocket))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n%s", err, flag.Usage)
			os.Exit(1)
		}

		containers, err := client.ContainerList(context.Background(), types.ContainerListOptions{Quiet: true})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n%s", err, flag.Usage)
			os.Exit(1)
		}
		logger.Debugf("%d container running", len(containers))
		for _, c := range containers {

			for _, n := range c.NetworkSettings.Networks {

				logger.Debugf("container %s ip %s adding into the blacklist", c.ID, n.IPAddress)
				cfg.BlacklistedHosts = append(cfg.BlacklistedHosts, n.IPAddress)  // Not enough addresses
				cfg.BlacklistedHosts = append(cfg.BlacklistedHosts, n.Aliases...) // Not enough protect!
				// http://proxy.local //

			}
		}

		cfg.Client = http.DefaultClient

	}
	logger.Debugf("Connecting to the InfluxDB2 server")
	database := db.InfluxDB{
		Client:            influxdb2.NewClient(demoFlags.Influxdb2Address, demoFlags.Infludb2Token),
		OrgName:           demoFlags.InfluxdbOrg,
		UserBucketName:    demoFlags.InfluxdbUserBucketName,
		RequestBucketName: demoFlags.InfluxdbRequestBucketName,
	}
	if err := database.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n%s", err, flag.Usage)
		os.Exit(1)
	}

	if len(extra) != 0 {
		switch extra[0] {
		case "set-user":
			gofakeit.Seed(time.Now().Unix())
			var user = pkg.User{
				Id:    gofakeit.Username(),
				Token: uuid.NewString(),
			}
			if len(extra) != 1 {
				user.Token = extra[1]
			}

			_user, err := database.GetToken(context.Background(), user.Token)
			if err != nil {
				if err = database.SetUser(context.Background(), user); err != nil {
					fmt.Printf(`{"error":"%v"}`, err)
				} else {
					user = _user
				}

			} else {
				r, _ := json.Marshal(user)
				fmt.Printf("%s", r)
			}

			os.Exit(0)

		}
	}

	logger.Debugf("listening on %s", demoFlags.ListenAddr)

	if err := pkg.New(cfg, database).Listen(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
