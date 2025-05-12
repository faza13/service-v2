package orm

import (
	"context"
	"gorm.io/gorm"
	"log"
	"service/config"
)

type Orm struct {
	db *gorm.DB
}

func NewProvider(cfg *config.Database) *Orm {
	//if cfg.Driver == "mariadb" {
	//	return newMysql(cfg)
	//}
	//
	//if cfg.Driver == "postgres" {
	//	return newPostgres(cfg)
	//}

	log.Fatal("not support database driver")
	return nil
}

func (d Orm) DB(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return d.db
	}

	return d.db.WithContext(ctx)
}

func (d Orm) WithTx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(txCtxKey{}).(*gorm.DB)
	if ok {
		return tx
	}

	return nil
}
