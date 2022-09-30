package service

import "errors"

var (
	// ErrServiceAlreadyStarted message to return when a service is called to start but is already running.
	ErrServiceAlreadyStarted = errors.New("service already started")

	// ErrServiceNotStarted message to return when service not started.
	ErrServiceNotStarted = errors.New("service not started")

	// ErrNilService is returned when service functions are called but the service is not instantiated.
	ErrNilService = errors.New("service not setup")

	// ErrServiceNotEnabled is returned when a service is called to start but is not enabled.
	ErrServiceNotEnabled = errors.New("service not enabled")

	// ErrServiceNotInitialized is returned when a service is called to start but has not been initialized.
	ErrServiceNotInitialized = errors.New("service not initialized")

	// ErrServiceNotFound is returned when a service can not be found.
	ErrServiceNotFound = errors.New("service not found")

	// ErrServiceNilWaitGroup is returned when a service has nil wait group.
	ErrServiceNilWaitGroup = errors.New("service nil wait group received")
)
