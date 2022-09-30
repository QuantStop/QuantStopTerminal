package database

import (
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/service"

	"sync"
	"time"
)

type Database struct {
	*service.Service
	*Config
}

func NewDatabase(config *Config) (*Database, error) {

	if config == nil {
		return nil, fmt.Errorf("create database failed, nil config received")
	}

	return &Database{
		Service: service.NewService("database", true),
		Config:  config,
	}, nil
}

// Start spawns the main process done by the service.
func (db *Database) Start(group *sync.WaitGroup) {
	db.Service.Start(group)
	t := time.NewTicker(time.Second * 5)

	defer func() {
		t.Stop()
		group.Done()
		log.Debugln(log.Database, db.GetName()+" shutdown completed.")
	}()

	// This lets the goroutine wait on communication from the channel
	// Docs: https://tour.golang.org/concurrency/5
	for {
		select {
		case <-db.Shutdown: // if channel message is shutdown finish loop
			return
		case <-t.C: // on channel tick check the connection
			/*err := s.CheckConnection()
			if err != nil {
				log.Error("database connection error: %v", err)
			}*/
			log.Debugln(log.Database, db.GetName()+" tick.")
		}
	}
}

// Stop terminates all processes belonging to the service, blocking until they are all terminated.
func (db *Database) Stop() error {
	if err := db.Service.Stop(); err != nil {
		return err
	}
	log.Debugln(log.Database, db.GetName()+service.MsgServiceShutdown)
	close(db.Shutdown)
	return nil
}
