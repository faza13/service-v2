package orm

import (
	"context"
	"gorm.io/gorm"
)

type IDatabase interface {
	DB(ctx context.Context) *gorm.DB
	WithTx(ctx context.Context) *gorm.DB
}

type ITransactor interface {
	WithTx(ctx context.Context, txFunc func(context.Context) error) error
}
