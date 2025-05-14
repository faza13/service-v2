package repositories

import (
	"github.com/elastic/go-elasticsearch/v9"
	"service/app/repositories/user"
	"service/pkg/cache"
	pkg_elastic "service/pkg/datastore/elastic"
	"service/pkg/datastore/orm"
)

type Repositories struct {
	Transactor  orm.ITransactor
	Cache       cache.ICache
	UserDB      *user.UserDB
	UserElastic *user.UserElasticRepo
}

func NewRepositories(db orm.IDatabase, Cache cache.ICache, elastic *elasticsearch.Client) *Repositories {
	userIdxElastic := pkg_elastic.NewElastic(elastic, "sekolahmu_user")
	return &Repositories{
		Cache:       Cache,
		Transactor:  orm.NewTransactor(db),
		UserDB:      user.NewUserRepo(db),
		UserElastic: user.NewUserElasticRepo(userIdxElastic),
	}
}
