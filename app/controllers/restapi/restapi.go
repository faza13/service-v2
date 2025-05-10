package restapi

import (
	"service/app/controllers/restapi/user"
	"service/app/usecases"
)

type Restapi struct {
	UserHandler *user.UserHandler
}

func NewRestapi(usecase *usecases.Usecase) *Restapi {
	return &Restapi{
		UserHandler: user.NewUserHandler(usecase.UserUsecase),
	}
}
