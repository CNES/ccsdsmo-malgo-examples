package service

type Service interface {
	CreateService() Service

	Start() error
}
