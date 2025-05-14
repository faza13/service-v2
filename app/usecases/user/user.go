package user

import "service/pkg/datastore/orm"

type UserUsecase struct {
	transactor orm.ITransactor
	userRepo   IUserRepo
}

func NewUserUsecase(transaaction orm.ITransactor, userRepo IUserRepo) *UserUsecase {
	return &UserUsecase{
		transactor: transaaction,
		userRepo:   userRepo,
	}
}
