package service

import "sync"

// IService exports an interface that is implemented by each Service.
// This is used by the implementation of the subsystem manager in the engine.
type IService interface {

	// Start spawns all processes done by the service.
	Start(wg *sync.WaitGroup) error

	// Run is the main thread of the process as a goroutine
	Run(wg *sync.WaitGroup)

	// Stop terminates all processes belonging to the service, blocking until they are all terminated.
	Stop() error

	// IsRunning returns true if the service is currently running.
	IsRunning() bool

	// IsEnabled returns true if the service is allowed to run.
	IsEnabled() bool

	// GetName returns the name of the service.
	GetName() string

	// Health returns nil if the service is healthy, otherwise returns the last error
	Health() error
}
