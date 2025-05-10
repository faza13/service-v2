package usecases

import (
	"service/app/repositories"
	"service/app/usecases/user"
)

type Usecase struct {
	UserUsecase *user.UserUsecase
}

func NewUsecase(repositories *repositories.Repositories) *Usecase {
	return &Usecase{
		UserUsecase: user.NewUserUsecase(repositories.UserDB),
	}
}
