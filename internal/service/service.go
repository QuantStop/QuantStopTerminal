package service

import (
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/log"
	"sync"
)

const (
	// MsgServiceInitializing message to return when subsystem is initializing.
	MsgServiceInitializing = " subsystem initializing..."

	// MsgServiceInitialized message to return when subsystem has initialized.
	MsgServiceInitialized = " subsystem initializing... Success."

	// MsgServiceStarting message to return when subsystem is starting up.
	MsgServiceStarting = " subsystem starting..."

	// MsgServiceStarted message to return when subsystem has started.
	MsgServiceStarted = " subsystem starting... Success."

	// MsgServiceShuttingDown message to return when a subsystem is shutting down.
	MsgServiceShuttingDown = " subsystem shutting down..."

	// MsgServiceShutdown message to return when a subsystem has shutdown.
	MsgServiceShutdown = " subsystem shutting down ... Success"
)

// IService exports an interface to the Service type.
type IService interface {

	// Start spawns all processes done by the service.
	Start(wg *sync.WaitGroup)

	// Stop terminates all processes belonging to the service, blocking until they are all terminated.
	Stop() error

	// IsRunning returns true if the service is currently running.
	IsRunning() bool

	// IsEnabled returns true if the service is allowed to run.
	IsEnabled() bool

	// GetName returns the name of the service.
	GetName() string
}

// Service is the base type for all services defined in the project.
type Service struct {
	name     string
	enabled  bool
	started  bool
	Shutdown chan struct{}
}

// NewService creates a pointer to a new Service struct with the provided values.
func NewService(name string, enabled bool) *Service {
	log.Debugln(log.Global, name+MsgServiceInitializing)
	return &Service{
		name:     name,
		enabled:  enabled,
		started:  false,
		Shutdown: make(chan struct{}),
	}
}

// Start is the main process for the service, run as a goroutine in the provided WaitGroup
func (service *Service) Start(wg *sync.WaitGroup) {
	if service == nil {
		log.Errorf(log.Global, "%s subsystem %w", service.name, ErrNilService)
	}
	if !service.enabled {
		log.Errorf(log.Global, "%s subsystem %w", service.name, ErrServiceNotEnabled)
	}
	if wg == nil {
		log.Errorf(log.Global, "%s subsystem %w", service.name, ErrServiceNilWaitGroup)
	}
	if service.started {
		log.Errorf(log.Global, "%s subsystem %w", service.name, ErrServiceAlreadyStarted)
	}
	service.started = true
	log.Debugln(log.Global, service.name+MsgServiceStarting)
}

// Stop The function to stop the service
func (service *Service) Stop() error {
	if service == nil {
		return fmt.Errorf("%s subsystem %w", service.name, ErrNilService)
	}
	if !service.started {
		return fmt.Errorf("%s subsystem %w", service.name, ErrServiceNotStarted)
	}
	service.started = false
	log.Debugln(log.Global, service.name+MsgServiceShuttingDown)
	return nil
}

// IsRunning checks whether the service is running
func (service *Service) IsRunning() bool {
	if service == nil {
		return false
	}
	return service.started
}

// IsEnabled checks whether the service is enabled or not
func (service *Service) IsEnabled() bool {
	return service.enabled
}

// GetName returns the subsystems name
func (service *Service) GetName() string {
	return service.name
}
