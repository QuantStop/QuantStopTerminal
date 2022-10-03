package engine

import (
	"errors"
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/database"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/service"
	"github.com/quantstop/quantstopterminal/internal/webserver"
	"reflect"
	"sync"
)

var (
	errNilServiceManager = errors.New("error: ServiceManager instance is nil")
)

// ServiceManager provides a useful pattern for managing services.
// It allows for ease of dependency management and ensures services
// dependent on others use the same references in memory.
type ServiceManager struct {

	// map of concrete implementations to service interface
	services map[reflect.Type]service.IService

	// keep an ordered slice of registered service types.
	serviceTypes []reflect.Type

	// service wait group
	ServiceWG *sync.WaitGroup
}

// NewServiceManager creates a new instance of the ServiceManager struct, then creates and registers all services.
func NewServiceManager(bot *Engine) (*ServiceManager, error) {

	sm := &ServiceManager{
		services:  make(map[reflect.Type]service.IService),
		ServiceWG: &sync.WaitGroup{},
	}

	if err := sm.registerDatabaseService(bot); err != nil {
		return nil, err
	}
	if err := sm.registerWebserverService(bot); err != nil {
		return nil, err
	}

	return sm, nil
}

func (s *ServiceManager) registerDatabaseService(bot *Engine) error {

	// Create Database
	db, err := database.NewDatabase(bot.Config.DatabaseConfig)
	if err != nil {
		return err
	}

	// Register Database
	if err = s.registerService(db); err != nil {
		return err
	}

	return nil
}

func (s *ServiceManager) registerWebserverService(bot *Engine) error {

	// Fetch dependencies
	var db *database.Database
	if err := s.FetchService(&db); err != nil {
		log.Error(log.Global, err)
	}

	// Create Webserver
	ws, err := webserver.NewWebserver(bot.Config.WebserverConfig, db)
	if err != nil {
		return err
	}

	// Register Webserver
	if err = s.registerService(ws); err != nil {
		return err
	}

	return nil
}

// RegisterService appends a service to the service registry.
func (s *ServiceManager) registerService(service service.IService) error {
	kind := reflect.TypeOf(service)
	if _, exists := s.services[kind]; exists {
		return fmt.Errorf("service already exists: %v", kind)
	}
	s.services[kind] = service
	s.serviceTypes = append(s.serviceTypes, kind)
	return nil
}

// SetService allows starting/stopping individual services.
// Not all services are allowed as some are core features that rely on other services as well.
// If the service is not found, or the service is not allowed, an error will be returned.
func (s *ServiceManager) SetService(name string, status bool) error {
	if s == nil {
		return errNilServiceManager
	}

	//engineMutex.Lock()
	//defer engineMutex.Unlock()

	//var err error
	/*switch strings.ToLower(name) {

	case NTPSubsystemName:
		if status {
			return bot.NTPCheckerSubsystem.Start()
		} else {
			return bot.NTPCheckerSubsystem.Stop()
		}

	case TraderSubsystemName:
		if status {
			return bot.TraderSubsystem.Start()
		} else {
			return bot.TraderSubsystem.Stop()
		}

	case InternetCheckerName:
		if status {
			return bot.InternetSubsystem.Start()
		} else {
			return bot.InternetSubsystem.Stop()
		}

	}
	return fmt.Errorf("%s: %w", name, service.ErrServiceNotFound)*/
	return nil
}

func (s *ServiceManager) GetService(name string) bool {
	return false
}

// Start takes in a subsystem name string, and wait group, attempts to find the service by its name, and then start it.
func (s *ServiceManager) Start(name string) error {

	found := false

	for _, kind := range s.serviceTypes {
		if s.services[kind].GetName() == name {
			found = true
			if !s.services[kind].IsEnabled() {
				return service.ErrSubsystemNotEnabled
			}
			if s.services[kind].IsRunning() {
				return service.ErrServiceAlreadyStarted
			}
			log.Debugf(log.Global, "Starting service type %v", kind)
			s.ServiceWG.Add(1)
			go s.services[kind].Start(s.ServiceWG)
			return nil
		}
	}

	if !found {
		return fmt.Errorf("unknown service: %s", name)
	}

	return nil
}

// StartAll initialized each service in order of registration.
func (s *ServiceManager) StartAll() error {
	if s == nil {
		return errNilServiceManager
	}
	log.Debugf(log.Global, "Found %d services: %v", len(s.serviceTypes), s.serviceTypes)
	for _, kind := range s.serviceTypes {
		if !s.services[kind].IsEnabled() {
			log.Debugf(log.Global, "Service %v disabled", kind)
			continue
		}
		if s.services[kind].IsRunning() {
			log.Debugf(log.Global, "Service %v already started", kind)
			continue
		}
		log.Debugf(log.Global, "Starting service type %v", kind)
		s.ServiceWG.Add(1)
		go s.services[kind].Start(s.ServiceWG)
	}
	return nil
}

// StopAll ends every service in reverse order of registration, logging a panic if any of them fail to stop.
func (s *ServiceManager) StopAll() {
	for i := len(s.serviceTypes) - 1; i >= 0; i-- {
		kind := s.serviceTypes[i]
		if s.services[kind].IsEnabled() && s.services[kind].IsRunning() {
			if err := s.services[kind].Stop(); err != nil {
				log.Errorf(log.Global, "Could not stop the following service: %v, %v", kind, err)
			}
		}
	}
}

// Stop takes in a subsystem name string and attempts to find the service by its name, and then Stop
func (s *ServiceManager) Stop(name string) error {

	found := false

	for _, kind := range s.serviceTypes {
		if s.services[kind].GetName() == name {
			found = true
			if !s.services[kind].IsEnabled() {
				return service.ErrSubsystemNotEnabled
			}
			if !s.services[kind].IsRunning() {
				return service.ErrServiceAlreadyStopped
			}
			log.Debugf(log.Global, "Stopping service type %v", kind)
			if err := s.services[kind].Stop(); err != nil {
				return err
			}
			return nil
		}
	}

	if !found {
		return fmt.Errorf("unknown service: %s", name)
	}

	return nil
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
