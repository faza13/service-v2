package mariadb

import (
	"context"
	"gorm.io/gorm"
)

type transactor struct {
	db *gorm.DB
}

var _ ITransactor = &transactor{}

func NewTransactor(db IDatabase) *transactor {
	return &transactor{db: db.DB(nil)}
}

func (d *transactor) WithTx(ctx context.Context, txFunc func(context.Context) error) error {
	tx := d.db.WithContext(ctx).Begin()

	txCtx := txToContext(ctx, tx)
	if err := txFunc(txCtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func txToContext(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txCtxKey{}, tx)
}
