package middleware

import (
	"github.com/quantstop/quantstopterminal/internal/log"
	"net/http"
	"strings"
	"time"
)

// HTTPReqInfo describes info about HTTP request
type HTTPReqInfo struct {
	// GET etc.
	method  string
	uri     string
	referer string
	ipaddr  string
	// response code, like 200, 404
	//code int
	// number of bytes of the response sent
	//size int64
	// how long did it take to
	//duration time.Duration
	userAgent string
}

func HttpRequestLogger(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		reqInfo := &HTTPReqInfo{
			method:    request.Method,
			uri:       request.URL.String(),
			referer:   request.Header.Get("Referer"),
			userAgent: request.Header.Get("User-Agent"),
		}

		reqInfo.ipaddr = getRemoteAddressFromHttpRequest(request)

		//todo: log response code, size, duration
		logHttpRequest(reqInfo)
		handlerFunc(response, request)
	}
}

// getIpFromRemoteAddress returns ip address from http.Request.RemoteAddress
func getIpFromRemoteAddress(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// GetRemoteAddressFromHttpRequest returns ip address of the client making the request
func getRemoteAddressFromHttpRequest(request *http.Request) string {
	header := request.Header
	xRealIP := header.Get("X-Real-Ip")
	xForwardedFor := header.Get("X-Forwarded-For")
	if xRealIP == "" && xForwardedFor == "" {
		return getIpFromRemoteAddress(request.RemoteAddr)
	}
	if xForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(xForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// todo: should return first non-local address
		return parts[0]
	}
	return xRealIP
}

// logHttpRequest logs information about http requests
func logHttpRequest(ri *HTTPReqInfo) {
	log.Infoln(log.Webserver, ri.ipaddr+" "+time.Now().String()+" "+ri.method+" "+ri.uri+" "+ri.referer+" "+ri.userAgent)
}
