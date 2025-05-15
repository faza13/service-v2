package orm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"service/config"
	"strconv"
	"time"
)

func newMysql(dbCfg *config.Database) *Orm {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// Access the raw *sql.DB object
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
	// Maximum number of idle connections

	if err != nil {
		log.Fatal(fmt.Sprintf("failed to connect database mysql: %s", err.Error()))
	}
	return &Orm{db: db}
}
