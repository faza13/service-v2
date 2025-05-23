package container

import (
	"context"
	"fmt"
	"service/app/controllers/broker"
	"service/app/controllers/grpc"
	"service/app/controllers/restapi"
	"service/app/middlewares"
	"service/app/repositories"
	"service/app/usecases"
	"service/config"
	"service/pkg/cache"
	"service/pkg/datastore/elastic"
	"service/pkg/datastore/orm"
	"service/pkg/message_broker/kafka"
	"service/pkg/otel"
	"service/pkg/server"
	"service/pkg/setting"
	"service/routes/api"
	brokerRouter "service/routes/broker"
)

func StartApp(ctx context.Context) {
	cfg := config.NewConfig()
	setting.NewSetting(&cfg)

	db := orm.NewProvider(&cfg.Database)
	cache := cache.NewCache(ctx, &cfg)
	esClient := elastic.NewElasticClient(ctx, &cfg)
	//mongoClient := mongodb.NewMongodb(ctx, &cfg)

	repo := repositories.NewRepositories(db, cache, esClient, nil)
	//repo := repositories.NewRepositories(db, cache, esClient, mongoClient)
	uc := usecases.NewUsecase(repo)
	rest := restapi.NewRestapi(uc)
	mid := middlewares.NewMiddlewares()

	traceProvider, teardown, err := otel.InitTracing(otel.Config{
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
	pub, sub := setupKafka(ctx, &cfg.Kafka)
	defer pub.Close()

	//run message broker
	go kafka.NewMessageBroker(
		ctx,
		&cfg.Kafka,
		pub,
		sub,
		brokerHandler,
		brokerRouter.NewUserBroker)

	grpxController := grpc.NewGrpc(ctx, uc)

	go server.RunGrpcServer(ctx, &cfg, grpxController)

	// run http server
	server.RunHTTPServer(ctx, &cfg, tracer, rest, mid, api.NewUserApi, api.NewPermissionApi)
}
