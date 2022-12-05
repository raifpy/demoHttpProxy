package pkg

import "context"

type UserDatabase interface {
	UpdateUser(context.Context, User) error
}

type RequestDatabase interface {
	UpdateRequest(context.Context, UserRequest) error
}

type Database interface {
	GetToken(ctx context.Context, token string) (User, error)
	SetRequest(context.Context, UserRequest) error
	UpdateRequest(context.Context, UserRequest) error
	ListenUserRequests(context.Context, User) chan UserRequest
	UserDatabase
}
