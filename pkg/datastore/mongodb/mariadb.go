package mongodb

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"service/config"
)

type Orm struct {
	db *gorm.DB
}

func New(dbCfg *config.Database) *Orm {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Orm{db: db}
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
