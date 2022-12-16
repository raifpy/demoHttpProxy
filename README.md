# proxyApiGateway

Demo project. Not ready for production!

## build

Ensure that docker instance is available.
```bash
git clone https://github.com/raifpy/proxyApiGateway
cd proxyApiGateway/scripts
make
make run
```

## usage

```
Usage:
  ./cmd/demo [OPTIONS]

Application Options:
      --verbose                   Enable verbose output
  -d, --docker-socket=            Docker unix socket to find container local ip
  -i, --influxdb2-address=        Adress to access influxdb2 [REQUIRED]
  -t, --influxdb2-token=          Token to access influxdb2 [REQUIRED]
  -o, --influxdb2-org=            OrgName for influxdb2 [REQUIRED]
      --influxdb2-user-bucket=    Bucket name for users (default: user)
      --influxdb2-request-bucket= Bucket name for requests (default: request)
  -b, --black-list=               Blacklist to block host request (default: localhost)
  -l, --listen-addr=              Listen port (default: 0.0.0.0:8080)
  -q, --enable-queue              Enable Queue
      --queue-limit=              Queue for requests of the same token (default: 10)
      --queue-wait-timeout=       Queue wait timeout to drop waiting request. MS (default: 5000)

Help Options:
  -h, --help                      Show this help message

```

Configure:
```
type ServerConfig struct {
	Logger           logrus.FieldLogger
	Client           *http.Client
	ListenAddr       string
	BlacklistedHosts []string
	QueueLimit       int
	WaitQueueTimeout time.Duration
	WaitQueue        bool
}
```

Default credentials
```
username: admin
password: 12345678
influxdb2-token: l7ZBNSgx7h4nJvPB1WKLn8pMX3_NdABMNTfJ1QcP_DuQ545qApj9ao7yXAgCeLmChfv4urUMtepBX5M8FQ5zDw==
influxdb2-org: api_v1
proxy-token: 45ec72df20bf09a726ac65ffdd5ef652fb8b5ba06f-test
```

Usage: 

```
/opt/demo/demo -i $INFLUX_HOST -t $INFLUX_TOKEN -o $INFLUX_ORG --verbose -d $DOCKER_SOCK_PATH -b localhost -l $LISTEN_ADDR -q
```

plain
```
curl -v "http://localhost:8080/?token=45ec72df20bf09a726ac65ffdd5ef652fb8b5ba06f-test&url=https://httpbin.org/anything"
```
proxy mode
```
curl -x -k "http://45ec72df20bf09a726ac65ffdd5ef652fb8b5ba06f-test:@localhost:8080" "http://httpbin.org/anything"
```

CONNECT proxy type is not supporting. You can't request a https web site while proxy serving http over requesting proxy mode.