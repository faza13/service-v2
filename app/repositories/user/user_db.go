package user

import (
	"context"
	"service/pkg/datastore/mariadb"
)

type UserDB struct {
	db mariadb.IDatabase
}

func NewUserRepo(db mariadb.IDatabase) *UserDB {
	return &UserDB{
		db: db,
	}
}

func (r *UserDB) List(ctx context.Context) {
}

func (r *UserDB) Create(ctx context.Context, data interface{}) {
}
