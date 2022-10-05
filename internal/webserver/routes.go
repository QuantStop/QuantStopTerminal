package webserver

import (
	"github.com/quantstop/quantstopterminal/internal/database/models"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/webserver/errors"
	"github.com/quantstop/quantstopterminal/internal/webserver/middleware"
	"github.com/quantstop/quantstopterminal/internal/webserver/router"
	"github.com/quantstop/quantstopterminal/internal/webserver/websocket"
	"github.com/quantstop/quantstopterminal/internal/webserver/write"
	"net/http"
	"runtime/debug"
)

func (s *Webserver) ConfigureRoutes(isDev bool) {

	log.Debugln(log.Webserver, "Setting up middleware ... ")
	//s.mux.RegisterMiddleware()

	log.Debugln(log.Webserver, "Setting up error handlers ... ")
	s.MethodNotAllowed = write.Error(errors.BadRequestMethod)
	s.NotFound = write.Error(errors.RouteNotFound)
	s.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		log.Errorf(log.Webserver, "Panic on %s", r.URL.Path)
		debug.PrintStack()
		write.Error(errors.InternalError)(w, r)
	}

	log.Debugln(log.Webserver, "Setting up route handlers ... ")

	// Public routes
	/*s.GET("/api/test", handlers.Test, router.Public)
	s.GET("/api/version", handlers.GetVersion, router.Public)
	s.GET("/api/sub-status", handlers.GetSubsystemStatus, router.Public)
	s.GET("/api/uptime", handlers.GetUptime, router.Public)*/

	// Session routes
	s.POST("/api/session", s.Login, router.Public)
	s.DELETE("/api/session", s.Logout, router.User)
	s.GET("/api/refresh-token", s.RefreshToken, router.User)

	// User routes
	/*s.POST("/api/signup", handlers.Signup, router.Public)
	s.GET("/api/user", handlers.Whoami, router.User)
	s.mux.POST("/reset", handlers.CreateReset, router.User)
	s.mux.GET("/reset/([0-9]+)", handlers.DoReset, router.User)*/

	/* Exchange routes */
	/*s.GET("/api/exchanges", handlers.GetExchanges, router.User)
	s.GET("/api/exchanges/([^/]+)/products", handlers.GetProducts, router.User)
	s.GET("/api/exchanges/([^/]+)/products/([^/]+)/candles", handlers.GetCandles, router.User)*/

	// Admin routes
	/*s.GET("/api/get-users", handlers.GetAllUsers, router.Admin)
	s.POST("/api/set-subsystem", handlers.SetSubsystem, router.Admin)
	s.POST("/api/set-sysconfig", handlers.SetSystemConfig, router.Admin)
	s.POST("/api/exchange", handlers.SetExchange, router.Admin)*/

	log.Debugln(log.Webserver, "Setting up websocket handler ... ")
	s.WebsocketHandler = func(writer http.ResponseWriter, request *http.Request) {
		authHandler := router.AuthRoute(
			func(user *models.User, w http.ResponseWriter, r *http.Request) http.HandlerFunc {
				return func(writer http.ResponseWriter, request *http.Request) {
					websocket.ServeWs(s.Hub, writer, request)
				}
			},
			writer,
			request,
			router.User,
		)
		logHandler := middleware.HttpRequestLogger(authHandler)
		logHandler(writer, request)

	}

}
