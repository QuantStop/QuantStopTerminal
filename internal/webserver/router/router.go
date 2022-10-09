package router

import (
	"context"
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/webserver/middleware"
	"github.com/quantstop/quantstopterminal/internal/webserver/utils"
	"github.com/quantstop/quantstopterminal/web"
	"net/http"
	"regexp"
)

// Simple RegExp based Http Router
// Inspiration and core design from https://benhoyt.com/writings/go-routing/

type Router struct {
	devMode            bool
	routes             []Route
	middlewareHandlers []http.HandlerFunc
	MethodNotAllowed   http.HandlerFunc
	NotFound           http.HandlerFunc
	WebsocketHandler   http.HandlerFunc
	// Function to handle panics recovered from http handlers.
	// Used to keep the server from crashing because of un-recovered panics.
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})
}

// ctxKey context key for request context
type ctxKey struct{}

// GetField is a helper function for handlers to get parameters from url
func GetField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}

// New returns a pointer to a new Router
func New(devMode bool) (*Router, error) {
	return &Router{
		devMode: devMode,
	}, nil
}

// PrintRoutes to display in the console/log
func (r *Router) PrintRoutes() {
	for _, route := range r.routes {
		path := route.regex.String()
		switch route.method {
		case http.MethodGet:
			//log.Println("    " + route.method + "     ->  " + path)
			log.Debugf(log.Webserver, "    %s     ->  %s", route.method, path)
		case http.MethodPut:
			//log.Println("    " + route.method + "     ->  " + path)
			log.Debugf(log.Webserver, "    %s     ->  %s", route.method, path)
		case http.MethodPost:
			//log.Println("    " + route.method + "    ->  " + path)
			log.Debugf(log.Webserver, "    %s    ->  %s", route.method, path)
		case http.MethodDelete:
			//log.Println("    " + route.method + "  ->  " + path)
			log.Debugf(log.Webserver, "    %s  ->  %s", route.method, path)
		}
	}
}

// GET is a shortcut for creating a new Route
func (r *Router) GET(pattern string, handler AuthHandler, authType AuthType) Route {
	return r.Handle(http.MethodGet, pattern, handler, authType)
}

// PUT is a shortcut for creating a new Route
func (r *Router) PUT(pattern string, handler AuthHandler, authType AuthType) Route {
	return r.Handle(http.MethodPut, pattern, handler, authType)
}

// POST is a shortcut for creating a new Route
func (r *Router) POST(pattern string, handler AuthHandler, authType AuthType) Route {
	return r.Handle(http.MethodPost, pattern, handler, authType)
}

// DELETE is a shortcut for creating a new Route
func (r *Router) DELETE(pattern string, handler AuthHandler, authType AuthType) Route {
	return r.Handle(http.MethodDelete, pattern, handler, authType)
}

// Handle validates, creates, and appends a Route to the routes array
func (r *Router) Handle(httpMethod, pattern string, handler AuthHandler, authType AuthType) Route {

	// validate method, fatal un-recoverable if not valid
	if matches, err := regexp.MatchString("^[A-Z]+$", httpMethod); !matches || err != nil {
		log.Error(log.Webserver, "http method "+httpMethod+" is not valid")
	}

	// create the Route
	route := newRoute(httpMethod, pattern, r.wrap(handler, authType))

	// add it the routes array
	r.routes = append(r.routes, route)

	// return the Route
	return route
}

// RegisterMiddleware takes a handle function and adds it to the array of middleware handlers
func (r *Router) RegisterMiddleware(handler http.HandlerFunc) {
	r.middlewareHandlers = append(r.middlewareHandlers, handler)
}

// wrap does all the middleware together
func (r *Router) wrap(h AuthHandler, authType AuthType) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		//log.Println("Middleware chain | 0 | wrap")

		// 1: role based authentication middleware
		authHandler := AuthRoute(h, response, request, authType)

		// Handlers are executed in reverse order from where chain is built starting here

		// 4: csrf prevention middleware
		csrfHandler := middleware.Csrf(authHandler)

		// 3: cors middleware
		corsHandler := middleware.Cors(csrfHandler)

		// 2: log middleware
		logHandler := middleware.HttpRequestLogger(corsHandler)
		logHandler(response, request)

	}
}

// recover is the deferred function that calls the supplied PanicHandler on a panic condition
func (r *Router) recover(w http.ResponseWriter, req *http.Request) {
	var empty interface{}
	if err := recover(); err != empty {
		r.PanicHandler(w, req, err)
	}
}

// ServeHTTP implements the http.handler interface
func (r *Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	var head string

	// shift head and tail to get below "api/" part of the path
	head, _ = utils.ShiftPath(request.URL.Path)

	// looking for filesystem (frontend)
	if head != "api" && !r.devMode {
		serveFileContents(response, request)
		return
	}

	// defer panic
	if r.PanicHandler != nil {
		defer r.recover(response, request)
	}

	// handle websocket route
	if request.URL.Path == "/api/ws" {
		r.WebsocketHandler(response, request)
		return
	}

	// allow holds requests with invalid methods
	var allow []string

	// loop through all routes
	for _, route := range r.routes {

		// look for matches in the request path
		matches := route.regex.FindStringSubmatch(request.URL.Path)
		if len(matches) > 0 {

			// match found but request method doesn't match
			if request.Method != route.method && request.Method != http.MethodOptions {
				// add it to the array defined earlier
				allow = append(allow, route.method)
				continue
			}

			// match of request path and method found! execute the handler with context
			ctx := context.WithValue(request.Context(), ctxKey{}, matches[1:])
			route.handler(response, request.WithContext(ctx))
			return
		}
	}

	// return method not allowed for requests to path with invalid method
	if len(allow) > 0 {
		r.MethodNotAllowed(response, request)
		return
	}

	// no path was found at all ...
	r.NotFound(response, request)
}

func serveFileContents(w http.ResponseWriter, r *http.Request) {

	wfs, err := web.GetFileSystem()
	if err != nil {
		log.Errorf(log.Webserver, "error getting file system: %v", err)
	}
	httpFS := http.FS(wfs)

	path := r.URL.Path
	if path == "/" {
		path = "index.html"
	}
	log.Debugf(log.Webserver, "requesting file system at: %v", r.URL.Path)

	// Restrict only to instances where the browser is looking for html, css, or javascript files
	/*if !strings.Contains(r.Header.Get("Accept"), "text/html") ||
		!strings.Contains(r.Header.Get("Accept"), "text/css") ||
		!strings.Contains(r.Header.Get("Accept"), "text/javascript") {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 not found")
		return
	}*/

	// Open the file and return its contents using http.ServeContent
	index, err := httpFS.Open(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s not found", path)
		return
	}

	fi, err := index.Stat()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s not found", path)
		return
	}

	// todo: idk is this needed?
	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeContent(w, r, fi.Name(), fi.ModTime(), index)
}
