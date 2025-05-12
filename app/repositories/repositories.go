package repositories

import (
	"service/app/repositories/user"
	"service/pkg/datastore/orm"
)

type Repositories struct {
	Transactor  orm.ITransactor
	UserDB      *user.UserDB
	UserElastic *user.UserElasticRepo
}

func NewRepositories(db orm.IDatabase) *Repositories {
	return &Repositories{
		Transactor:  orm.NewTransactor(db),
		UserDB:      user.NewUserRepo(db),
		UserElastic: user.NewUserElasticRepo(),
	}
}
