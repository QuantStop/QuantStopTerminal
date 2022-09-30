package webserver

import (
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/database"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/service"

	"sync"
	"time"
)

type Webserver struct {
	*service.Service
	*Config
	*database.Database
}

func NewWebserver(config *Config, db *database.Database) (*Webserver, error) {

	if config == nil {
		return nil, fmt.Errorf("create webserver failed, nil config received")
	}

	if db == nil {
		return nil, fmt.Errorf("create webserver failed, nil database received")
	}

	return &Webserver{
		Service:  service.NewService("webserver", true),
		Config:   config,
		Database: db,
	}, nil
}

// Start spawns the main process done by the service.
func (ws *Webserver) Start(group *sync.WaitGroup) {
	ws.Service.Start(group)
	t := time.NewTicker(time.Second * 5)

	defer func() {
		t.Stop()
		group.Done()
		log.Debugln(log.Webserver, ws.GetName()+" shutdown completed.")
	}()

	// This lets the goroutine wait on communication from the channel
	// Docs: https://tour.golang.org/concurrency/5
	for {
		select {
		case <-ws.Shutdown: // if channel message is shutdown finish loop
			return
		case <-t.C: // on channel tick check the connection
			/*err := s.CheckConnection()
			if err != nil {
				log.Error("database connection error: %v", err)
			}*/
			log.Debugln(log.Webserver, ws.GetName()+" tick.")
		}
	}
}

// Stop terminates all processes belonging to the service, blocking until they are all terminated.
func (ws *Webserver) Stop() error {
	if err := ws.Service.Stop(); err != nil {
		return err
	}
	log.Debugln(log.Webserver, ws.GetName()+service.MsgServiceShutdown)
	close(ws.Shutdown)
	return nil
}
