package container

import (
	"context"
	kafkasdk "github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"service/app/controllers/broker"
	"service/config"
	"service/pkg/message_broker/kafka"
	routerBroker "service/routes/broker"
)

func setupKafka(ctx context.Context, cfg *config.Kafka) (*kafkasdk.Subscriber, *kafkasdk.Publisher) {
	sub := kafka.NewSubscriber(ctx, cfg)
	pub := kafka.NewProducer(ctx, cfg)

	return sub, pub
}

func runMessageBroker(
	ctx context.Context,
	cfg *config.Kafka,
	sub *kafkasdk.Subscriber,
	pub *kafkasdk.Publisher,
	handler *broker.BrokerHandler,
) {

	go func(
		ctx context.Context,
		cfg *config.Kafka,
		sub *kafkasdk.Subscriber,
		pub *kafkasdk.Publisher,
		handler *broker.BrokerHandler) {

		kafka.NewMessageBroker(
			ctx,
			cfg,
			sub,
			pub,
			handler,
			routerBroker.NewUserBroker)

	}(ctx, cfg, sub, pub, handler)
}
