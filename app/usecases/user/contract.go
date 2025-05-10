package user

import "context"

type IUserRepo interface {
	List(ctx context.Context)
	Create(ctx context.Context, data interface{})
}
