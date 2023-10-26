package api

type EcnSubscribeService interface {
	// SubscribeAllServices Subscribe all services changes in Nacos:
	SubscribeAllServices(SubscribeCallback func())

	// SubscribeService Subscribe one service changes in Nacos:
	SubscribeService(ServiceName string, SubscribeCallback func())
}
