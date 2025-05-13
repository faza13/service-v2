package repositories

import (
	"service/app/repositories/user"
	loc "service/pkg/cache"
	"service/pkg/datastore/orm"
)

type Repositories struct {
	Transactor  orm.ITransactor
	Cache       cache.ICache
	UserDB      *user.UserDB
	UserElastic *user.UserElasticRepo
}

func NewRepositories(db orm.IDatabase, Cache cache.ICache) *Repositories {
	return &Repositories{
		Cache:       Cache,
		Transactor:  orm.NewTransactor(db),
		UserDB:      user.NewUserRepo(db),
		UserElastic: user.NewUserElasticRepo(),
	}
}
