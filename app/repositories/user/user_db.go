package user

import (
	"context"
	"service/pkg/datastore/orm"
)

type UserDB struct {
	db orm.IDatabase
}

func NewUserRepo(db orm.IDatabase) *UserDB {
	return &UserDB{
		db: db,
	}
}

func (r *UserDB) List(ctx context.Context) {
}

func (r *UserDB) Create(ctx context.Context, data interface{}) {
}
