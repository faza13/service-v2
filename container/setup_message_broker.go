package container

import (
	"context"
	kafkasdk "github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"service/config"
	"service/pkg/message_broker/kafka"
)

func setupKafka(ctx context.Context, cfg *config.Kafka) (*kafkasdk.Subscriber, *kafkasdk.Publisher) {
	sub := kafka.NewSubscriber(ctx, cfg)
	pub := kafka.NewProducer(ctx, cfg)

	return sub, pub
}
