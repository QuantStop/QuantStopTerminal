package engine

import (
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/service"
	"reflect"
	"sync"
)

// ServiceManager provides a useful pattern for managing services.
// It allows for ease of dependency management and ensures services
// dependent on others use the same references in memory.
type ServiceManager struct {

	// map of concrete implementations to service interface
	services map[reflect.Type]service.IService

	// keep an ordered slice of registered service types.
	serviceTypes []reflect.Type
}

// NewServiceManager creates a new instance of the ServiceManager struct
func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		services: make(map[reflect.Type]service.IService),
	}
}

// RegisterService appends a service to the service registry.
func (s *ServiceManager) RegisterService(service service.IService) error {
	kind := reflect.TypeOf(service)
	if _, exists := s.services[kind]; exists {
		return fmt.Errorf("service already exists: %v", kind)
	}
	s.services[kind] = service
	s.serviceTypes = append(s.serviceTypes, kind)
	return nil
}

// StartAll initialized each service in order of registration.
func (s *ServiceManager) StartAll(group *sync.WaitGroup) {
	log.Debugf(log.Global, "Starting %d services: %v", len(s.serviceTypes), s.serviceTypes)
	for _, kind := range s.serviceTypes {
		log.Debugf(log.Global, "Starting service type %v", kind)
		group.Add(1)
		go s.services[kind].Start(group)
	}
}

// StopAll ends every service in reverse order of registration, logging a panic if any of them fail to stop.
func (s *ServiceManager) StopAll() {
	for i := len(s.serviceTypes) - 1; i >= 0; i-- {
		kind := s.serviceTypes[i]
		if err := s.services[kind].Stop(); err != nil {
			log.Errorf(log.Global, "Could not stop the following service: %v, %v", kind, err)
		}
	}
}

// FetchService takes in a struct pointer and sets the value of that pointer
// to a service currently stored in the service registry. This ensures the input argument is
// set to the right pointer that refers to the originally registered service.
func (s *ServiceManager) FetchService(service interface{}) error {
	if reflect.TypeOf(service).Kind() != reflect.Ptr {
		return fmt.Errorf("input must be of pointer type, received value type instead: %T", service)
	}
	element := reflect.ValueOf(service).Elem()
	if running, ok := s.services[element.Type()]; ok {
		element.Set(reflect.ValueOf(running))
		return nil
	}
	return fmt.Errorf("unknown service: %T", service)
}
