package service

import (
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/log"
	"sync"
)

const (
	// MsgServiceInitializing message to return when service is initializing.
	MsgServiceInitializing = " service initializing..."

	// MsgServiceInitialized message to return when service has initialized.
	MsgServiceInitialized = " service initializing... Success."

	// MsgServiceStarting message to return when service is starting up.
	MsgServiceStarting = " service starting..."

	// MsgServiceStarted message to return when service has started.
	MsgServiceStarted = " service starting... Success."

	// MsgServiceShuttingDown message to return when a service is shutting down.
	MsgServiceShuttingDown = " service shutting down..."

	// MsgServiceShutdown message to return when a service has shutdown.
	MsgServiceShutdown = " service shutting down ... Success"
)

// Service is the base type for all services defined in the project.
type Service struct {
	name     string
	enabled  bool
	started  bool
	lastErr  error
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
		log.Errorf(log.Global, "%s service %w", service.name, ErrNilService)
	}
	if !service.enabled {
		log.Errorf(log.Global, "%s service %w", service.name, ErrServiceNotEnabled)
	}
	if wg == nil {
		log.Errorf(log.Global, "%s service %w", service.name, ErrServiceNilWaitGroup)
	}
	if service.started {
		log.Errorf(log.Global, "%s service %w", service.name, ErrServiceAlreadyStarted)
	}
	service.started = true
	log.Debugln(log.Global, service.name+MsgServiceStarting)
}

// Stop The function to stop the service
func (service *Service) Stop() error {
	if service == nil {
		return fmt.Errorf("%s service %w", service.name, ErrNilService)
	}
	if !service.started {
		return fmt.Errorf("%s service %w", service.name, ErrServiceAlreadyStopped)
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

func (service *Service) Health() error {
	return service.lastErr
}
