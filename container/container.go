package container

import (
	"context"
	"fmt"
	"os"
	"service/app/controllers/broker"
	"service/app/controllers/restapi"
	"service/app/middlewares"
	"service/app/repositories"
	"service/app/usecases"
	"service/config"
	"service/pkg/datastore/mariadb"
	"service/pkg/logger"
	"service/pkg/otel"
	"service/pkg/server"
	"service/routes/api"
)

func StartApp(ctx context.Context) {
	var log *logger.Logger
	cfg := config.NewConfig()
	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}
	traceIDFn := func(ctx context.Context) string {
		return otel.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, cfg.App.Name, traceIDFn, events)

	db := mariadb.New(cfg.Database)
	repo := repositories.NewRepositories(db)
	uc := usecases.NewUsecase(repo)
	rest := restapi.NewRestapi(uc)
	mid := middlewares.NewMiddlewares()

	traceProvider, teardown, err := otel.InitTracing(log, otel.Config{
		ServiceName: cfg.Otel.ServiceName,
		Host:        cfg.Otel.HostTempo,
		Probability: 0.05,
	})

	if err != nil {
		panic(fmt.Errorf("starting tracing: %w", err))
	}

	defer teardown(ctx)

	tracer := traceProvider.Tracer(cfg.Otel.ServiceName)
	tracer.Start(ctx, "main")

	// run message broker
	brokerHandler := broker.NewBroker()
	pub, sub := setupKafka(ctx, cfg.Kafka)
	defer pub.Close()

	runMessageBroker(ctx, cfg.Kafka, pub, sub, brokerHandler)

	// run http server
	server.RunHTTPServer(ctx, cfg.Router, tracer, rest, mid, api.NewUserApi, api.NewPermissionApi)
}
