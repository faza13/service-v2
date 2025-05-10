package user

import (
	"context"
)

type UserElasticRepo struct {
}

func NewUserElasticRepo() *UserElasticRepo {
	return &UserElasticRepo{}
}

func (r *UserElasticRepo) List(ctx context.Context) {

}
