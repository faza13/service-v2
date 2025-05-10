package broker

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Updated(msg *message.Message) error {
	fmt.Printf(
		"\n> Received message: %s\n> %s\n> metadata: %v\n\n",
		msg.UUID, string(msg.Payload), msg.Metadata,
	)
	return nil
}
