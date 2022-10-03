package service

// IServiceManager exports an interface that is implemented by the engine to manage all services
type IServiceManager interface {
	SetService(name string, status bool)
	GetService(name string) bool
}
