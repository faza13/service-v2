package user

import "context"

type IUserUsecase interface {
	List(ctx context.Context) map[string]interface{}
}
