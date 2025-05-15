package orm

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"service/config"
	"strconv"
	"time"
)

func newPostgres(dbCfg *config.Database) *Orm {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		dbCfg.Host, dbCfg.User, dbCfg.Password, dbCfg.Name, dbCfg.Port, "UTC",
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to connect database postgres: %s", err.Error()))
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("Error getting *sql.DB object:", err)

	}

	maxOpenConn, _ := strconv.Atoi(dbCfg.MaxOpenConn)
	maxIdleConn, _ := strconv.Atoi(dbCfg.MaxIdleConn)
	connMaxLifetime, _ := strconv.Atoi(dbCfg.MaxConnLifetime)

	// Configure connection pooling
	if maxIdleConn != 0 {
		// Maximum number of open connections
		sqlDB.SetMaxOpenConns(maxOpenConn)
	}

	if maxIdleConn != 0 {
		// Maximum number of open connections
		sqlDB.SetMaxIdleConns(maxIdleConn)
	}

	if connMaxLifetime != 0 {
		// Maximum amount of time a connection can be reused (0 means no limit)
		sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
	}
	return &Orm{db: db}
}
