package repositories

import (
	"service/app/repositories/user"
	"service/pkg/datastore/mariadb"
)

type Repositories struct {
	Transactor  mariadb.ITransactor
	UserDB      *user.UserDB
	UserElastic *user.UserElasticRepo
}

func NewRepositories(db mariadb.IDatabase) *Repositories {
	return &Repositories{
		Transactor:  mariadb.NewTransactor(db),
		UserDB:      user.NewUserRepo(db),
		UserElastic: user.NewUserElasticRepo(),
	}
}
