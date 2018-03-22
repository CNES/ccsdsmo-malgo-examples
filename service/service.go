package service

type Service interface {
	CreateService() Service

	StartConsumer() error

	StartProvider() error
}
