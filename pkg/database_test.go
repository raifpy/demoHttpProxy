package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"
)

type DbSim struct {
	Users []User
}

func (d DbSim) GetToken(_ context.Context, t string) (User, error) {
	time.Sleep(time.Millisecond * 100) // Database delay?
	for _, u := range d.Users {
		if u.Token == t {
			return u, nil
		}
	}
	return User{}, errors.New("not exists")
}

func (d DbSim) ListenUserRequests(context.Context, User) chan UserRequest {
	panic("not implemented")
}

func (d DbSim) SetRequest(_ context.Context, r UserRequest) error {
	file, err := os.Create("/tmp/" + r.RequestId + ".json")
	if err != nil {
		return nil
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(r)
}

func (d DbSim) UpdateUser(context.Context, User) error {
	panic("not implemented")

}

func (d DbSim) UpdateRequest(ctx context.Context, r UserRequest) error {
	file, err := os.OpenFile("/tmp/"+r.RequestId+".json", os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	var old UserRequest
	if err := json.NewDecoder(file).Decode(&old); err != nil {
		return err
	}
	old.UpdateTime = r.UpdateTime
	old.Status = r.Status
	old.Error = r.Error
	old.ResponseContentType = r.ResponseContentType
	old.ResponseSize = r.ResponseSize
	old.ResponseStatus = r.ResponseStatus

	file.Close()

	if file, err = os.Create("/tmp/" + r.RequestId + ".json"); err != nil {
		return err
	}
	return json.NewEncoder(file).Encode(old)

}

func (d DbSim) SetUser(ctx context.Context, u User) error {
	return errors.New("setuser dbsim not implemented")
}

var TestDb = DbSim{
	Users: []User{
		{
			Id:    "user_1",
			Token: "user_token_1",
		},
	},
}
