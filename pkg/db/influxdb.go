package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/raifpy/proxyApiGateway/pkg"
)

type InfluxDB struct {
	Client            influxdb2.Client
	OrgName           string
	UserBucketName    string
	RequestBucketName string
}

var (
	_ pkg.Database = InfluxDB{}
)

// Call once
func (i InfluxDB) Init() error {
	ctx, cl := context.WithTimeout(context.Background(), time.Second*10)
	defer cl()
	if health, err := i.Client.Health(ctx); err != nil || health.Status != domain.HealthCheckStatusPass {
		if err == nil {
			err = fmt.Errorf("influxdb health check failed %s", health.Status)

		}
		return err
	}

	// org, err := i.Client.OrganizationsAPI().CreateOrganizationWithName(ctx, i.OrgName)
	// if !errors.Is(err, nil) {
	// 	return err
	// }

	org, err := i.Client.OrganizationsAPI().FindOrganizationByName(ctx, i.OrgName)
	if err != nil {
		return err
	}

	if _, err := i.Client.BucketsAPI().CreateBucketWithName(context.Background(), org, i.UserBucketName); err != nil && !(strings.Contains(err.Error(), "already exists")) {
		return err
	}

	if _, err := i.Client.BucketsAPI().CreateBucketWithName(context.Background(), org, i.RequestBucketName); err != nil && !(strings.Contains(err.Error(), "already exists")) {
		return err
	}
	return nil

}

//var _ pkg.Database = &InfluxDB{}

func (i InfluxDB) GetToken(ctx context.Context, token string) (u pkg.User, err error) {
	// I am begginer of influxdata. This query is very bad.
	result, err := i.Client.QueryAPI(i.OrgName).Query(context.Background(), fmt.Sprintf("from(bucket: \"%s\") |> range(start: 0) |> filter(onEmpty: \"drop\", fn: (r) => r.token == \"%s\") |> last()", i.UserBucketName, token))
	if err != nil {
		return u, err
	}
	defer result.Close()
	if !result.Next() {
		return u, errors.New("token not found")
	}
	if err = result.Err(); err != nil {
		return
	}

	u = pkg.UserFromMapI(result.Record().Values())
	return

}

func (i InfluxDB) GetRequest(ctx context.Context, id string) (ur pkg.UserRequest, err error) {
	// I am begginer of influxdb. This query is very ineffecient.

	result, err := i.Client.QueryAPI(i.OrgName).Query(context.Background(), fmt.Sprintf("from(bucket: \"%s\") |> range(start: 0) |> filter(onEmpty: \"drop\", fn: (r) => r.request_id == \"%s\") |> last()", i.RequestBucketName, id))
	if err != nil {
		return ur, err
	}
	defer result.Close()

	if err = result.Err(); err != nil {
		return
	}

	var resultmap = map[string]any{}

	for result.Next() {
		resultmap[result.Record().Field()] = result.Record().Value()
	}
	ur = pkg.UserRequestFromMapI(resultmap)
	return

}

func (i InfluxDB) SetRequest(ctx context.Context, ur pkg.UserRequest) error {
	i.Client.WriteAPIBlocking(i.OrgName, i.RequestBucketName).WriteRecord(context.Background())
	w := i.Client.WriteAPIBlocking(i.OrgName, i.RequestBucketName).WritePoint(ctx,
		influxdb2.NewPoint(ur.RequestId, map[string]string{
			"request_id": ur.RequestId,
			"method":     ur.Method,
			"user":       ur.UserId,
		}, ur.ToMapI(), time.Now()))
	return w
}

func (i InfluxDB) SetUser(ctx context.Context, u pkg.User) error {
	if len(u.Token) < 20 {
		return errors.New("invalid token length")
	}
	w := i.Client.WriteAPIBlocking(i.OrgName, i.UserBucketName).WritePoint(ctx,
		influxdb2.NewPoint(u.Id, map[string]string{
			"id":    u.Id,
			"token": u.Token,
		}, u.ToMapI(), time.Now()))

	return w

}

func (i InfluxDB) UpdateRequest(ctx context.Context, ur pkg.UserRequest) (err error) {
	var ur2 pkg.UserRequest
	if ur2, err = i.GetRequest(ctx, ur.RequestId); err != nil {
		return
	}
	ur2.Status = ur.Status
	ur2.UpdateTime = ur.UpdateTime
	ur2.ResponseContentType = ur.ResponseContentType
	ur2.ResponseSize = ur.ResponseSize
	ur2.ResponseStatus = ur.ResponseStatus
	ur2.Error = ur.Error

	return i.SetRequest(ctx, ur2)
}
