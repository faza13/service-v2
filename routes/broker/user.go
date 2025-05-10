package broker

import (
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"service/app/controllers/broker"
)

func NewUserBroker(
	router *message.Router,
	subs *kafka.Subscriber,
	_ *kafka.Publisher,
	handler *broker.BrokerHandler,
) {
	router.AddNoPublisherHandler(
		"user.updated",
		"user.updated",
		subs,
		handler.UserHandler.Updated)
}
