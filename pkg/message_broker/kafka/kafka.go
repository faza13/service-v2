package kafka

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"service/app/controllers/broker"
	"service/config"
	"strings"
	"time"
)

var logger = watermill.NewStdLogger(false, false)

type groupsHandlers func(
	router *message.Router,
	sub *kafka.Subscriber,
	pub *kafka.Publisher,
	brokerHandler *broker.BrokerHandler,
)

func NewSubscriber(ctx context.Context, cfg *config.Kafka) *kafka.Subscriber {
	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()

	kafkaHost := strings.Split(cfg.Host, ",")

	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               kafkaHost,
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaSubscriberConfig,
			ConsumerGroup:         "test_consumer_group",
		},
		watermill.NewStdLogger(false, false),
	)

	if err != nil {
		panic(err)
	}
	return subscriber
}

func NewProducer(ctx context.Context, cfg *config.Kafka) *kafka.Publisher {

	kafkaHost := strings.Split(cfg.Host, ",")

	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   kafkaHost,
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	return publisher
}

func NewMessageBroker(
	ctx context.Context,
	cfg *config.Kafka,
	sub *kafka.Subscriber,
	pub *kafka.Publisher,
	broker *broker.BrokerHandler,
	groupHandlers ...groupsHandlers,
) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	// SignalsHandler will gracefully shutdown Router when SIGTERM is received.
	// You can also close the router by just calling `r.Close()`.
	router.AddPlugin(plugin.SignalsHandler)

	router.AddMiddleware(
		// CorrelationID will copy the correlation id from the incoming message's metadata to the produced messages
		middleware.CorrelationID,

		// The handler function is retried if it returns an error.
		// After MaxRetries, the message is Nacked and it's up to the PubSub to resend it.
		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
			Logger:          logger,
		}.Middleware,

		// Recoverer handles panics from handlers.
		// In this case, it passes them as errors to the Retry middleware.
		middleware.Recoverer,
	)

	for i := range groupHandlers {
		groupHandlers[i](router, sub, pub, broker)
	}

	if err := router.Run(ctx); err != nil {
		panic(err)
	}

}
