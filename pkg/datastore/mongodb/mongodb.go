package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"service/config"
	"strconv"
	"strings"
	"time"
)

func NewMongodb(ctx context.Context, cfg *config.Config) *mongo.Client {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoOptions := options.Client().ApplyURI(cfg.Mongodb.Url)
	maxPoolConn, _ := strconv.ParseUint(cfg.Mongodb.MaxPoolConn, 10, 64)
	if maxPoolConn != 0 {
		mongoOptions.SetMaxPoolSize(maxPoolConn)
	}
	maxIdleConn, _ := strconv.ParseUint(cfg.Mongodb.MaxIdleConn, 10, 64)
	if maxIdleConn != 0 {
		duration := (time.Duration(maxPoolConn) * time.Hour)
		mongoOptions.SetMaxConnIdleTime(duration)
	}

	if cfg.Mongodb.Compression != "" {
		compression := strings.Split(",", cfg.Mongodb.Compression)
		mongoOptions.Compressors = compression
	}

	//.SetMaxPoolSize(50).SetCompressors([]string{"snappy"})
	if cfg.App.Environment == "development" || cfg.App.Environment == "staging" {
		cmdMonitor := &event.CommandMonitor{
			Started: func(_ context.Context, evt *event.CommandStartedEvent) {
				log.Print(evt.Command)
			},
		}
		mongoOptions.SetMonitor(cmdMonitor)
	}

	client, err := mongo.Connect(ctx, mongoOptions)

	if err != nil {
		log.Fatalf("failed to connect to mongodb: %s", err.Error())
	}

	return client
}
