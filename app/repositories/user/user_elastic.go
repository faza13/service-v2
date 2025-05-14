package user

import (
	"context"
	"service/pkg/datastore/elastic"
)

type UserElasticRepo struct {
	elastic *elastic.Elastic
}

func NewUserElasticRepo(elastic *elastic.Elastic) *UserElasticRepo {
	return &UserElasticRepo{
		elastic: elastic,
	}
}

func (r *UserElasticRepo) List(ctx context.Context) {

}
