package user

import "context"

type IUserUsecase interface {
	List(ctx context.Context) map[string]interface{}
	Register(ctx context.Context, request *RegistrationRequest) (interface{}, error)
}
