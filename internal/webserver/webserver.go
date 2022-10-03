package webserver

import (
	"context"
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/database"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/service"
	"github.com/quantstop/quantstopterminal/internal/system/crypto"
	"github.com/quantstop/quantstopterminal/internal/webserver/router"
	"github.com/quantstop/quantstopterminal/internal/webserver/websocket"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type Webserver struct {
	*Config
	*service.Service
	*http.Server
	*router.Router
	*websocket.Hub
	*database.Database
	service.IServiceManager
}

func NewWebserver(config *Config, db *database.Database) (*Webserver, error) {

	if config == nil {
		return nil, fmt.Errorf("create webserver failed, nil config received")
	}

	if db == nil {
		return nil, fmt.Errorf("create webserver failed, nil database received")
	}

	mux, err := router.New(config.DevMode, db)
	if err != nil {
		return nil, err
	}

	ws := &Webserver{
		Config:   config,
		Service:  service.NewService("webserver", true),
		Router:   mux,
		Database: db,
	}

	ws.Server = &http.Server{
		Addr:    config.HttpListenAddr,
		Handler: mux,
	}

	ws.ConfigureRouter(config.DevMode)
	ws.PrintRoutes()

	hub, err := websocket.NewHub(db, ws.Service.Shutdown)
	if err != nil {
		return nil, err
	}
	ws.Hub = hub
	log.Debugln(log.Webserver, ws.GetName()+service.MsgServiceInitialized)
	return ws, nil
}

// Start is the main process for the service, run as a goroutine.
func (s *Webserver) Start(serviceWG *sync.WaitGroup) error {
	if err := s.Service.Start(serviceWG); err != nil {
		return err
	}

	// run websocket server
	go s.Hub.Run(serviceWG)

	// run http server
	go s.Run(serviceWG)

	return nil
}

func (s *Webserver) Run(serviceWG *sync.WaitGroup) {
	// if dev mode, run node server
	if s.DevMode {
		go s.StartNodeDevelopmentServer(serviceWG)
	}

	// start the server
	if s.TLS {

		// if using tls, get the directory of certificates, then check them
		targetDir := crypto.GetTLSDir(s.ConfigDir)
		if err := crypto.CheckCerts(targetDir); err != nil {
			log.Errorf(log.Webserver, "checkCerts failed. err: %s\n", err)
		}

		log.Debugln(log.Webserver, s.GetName()+service.MsgServiceStarted)
		log.Infof(log.Webserver, s.GetName()+" listening on https://%v.\n", s.HttpListenAddr)

		// ListenAndServeTLS is blocking until Server.Shutdown() is called
		err := s.Server.ListenAndServeTLS(filepath.Join(targetDir, "cert.pem"), filepath.Join(targetDir, "key.pem"))
		if err == http.ErrServerClosed {
			err = nil // expected error after calling Server.Shutdown()
			log.Infof(log.Webserver, s.GetName()+" stopped listening on https://%v.\n", s.HttpListenAddr)
		} else if err != nil {
			log.Errorf(log.Webserver, "unexpected error from ListenAndServe: %s", err)
		}
	} else {
		log.Debugln(log.Webserver, s.GetName()+service.MsgServiceStarted)
		log.Infof(log.Webserver, s.GetName()+" listening on http://%v.\n", s.HttpListenAddr)

		// ListenAndServe is blocking until Server.Shutdown() is called
		err := s.Server.ListenAndServe()
		if err == http.ErrServerClosed {
			err = nil // expected error after calling Server.Shutdown()
			log.Infof(log.Webserver, s.GetName()+" stopped listening on http://%v.\n", s.HttpListenAddr)
		} else if err != nil {
			log.Errorf(log.Webserver, "unexpected error from ListenAndServe: %s", err)
		}
	}
	// anything here won't run until the Service.Shutdown channel is closed
	<-s.Service.Shutdown
	serviceWG.Done()
	log.Debugln(log.Webserver, s.GetName()+service.MsgServiceShutdown)
}

// Stop terminates all processes belonging to the service, blocking until they are all terminated.
func (s *Webserver) Stop() error {
	if err := s.Service.Stop(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := s.Server.Shutdown(ctx)
	if err != nil {
		log.Errorf(log.Webserver, "Webserver could not shutdown %v\n", err)
		close(s.Service.Shutdown)
	} else {
		close(s.Service.Shutdown)
	}
	return nil
}

func (s *Webserver) StartNodeDevelopmentServer(serviceWG *sync.WaitGroup) {
	serviceWG.Add(1)
	log.Debugf(log.Webserver, "Starting node development server ...")

	var cmd *exec.Cmd
	var err error

	cmd = exec.Command("npm", "run", "dev")
	cmd.Dir = "./web"
	cmd.Stdout = os.Stdout

	if err = cmd.Start(); err != nil {
		log.Errorf(log.Webserver, "Error starting node development server %v.\n", err)
	}

	<-s.Service.Shutdown
	serviceWG.Done()
	log.Infoln(log.Webserver, "Shutting down node development server ...")
	err = cmd.Process.Kill()
	if err != nil {
		log.Errorf(log.Webserver, "Error unable to kill process node development server %v.\n", err)
	} else {
		log.Infoln(log.Webserver, "Shutting down node development server ... Success.")
	}

	return

}
