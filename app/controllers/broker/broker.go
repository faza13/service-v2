package broker

type BrokerHandler struct {
	UserHandler *UserHandler
}

func NewBroker() *BrokerHandler {
	return &BrokerHandler{
		UserHandler: NewUserHandler(),
	}
}
