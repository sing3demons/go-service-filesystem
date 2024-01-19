package logger

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type contextKey string
const Key = "logger"

func Middleware(logger ILogger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return LogMiddleware(next, logger)
	}
}

func LogMiddleware(next http.Handler, logger ILogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := LogParentID(r, logger)
		r = setLoggerContext(r, l)
		next.ServeHTTP(w, r)
	})
}

func LogParentID(r *http.Request, logger ILogger) ILogger {
	xParent := r.Header.Get("X-Parent-ID")
	if xParent == "" {
		xParent = uuid.NewString()
	}
	xSpan := uuid.NewString()

	return logger.With(zap.Any("headers", GetHeaders(r)), zap.String("parent-id", xParent), zap.String("span-id", xSpan))
}

func setLoggerContext(r *http.Request, val ILogger) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), Key, val))
}

func getMACAndIP() MacIP {
	interfaces, _ := net.Interfaces()
	macAddr := MacIP{}
	for _, iface := range interfaces {

		if iface.Name != "" {
			macAddr.InterfaceName = iface.Name
		}

		if iface.HardwareAddr != nil {
			macAddr.HardwareAddr = iface.HardwareAddr.String()
		}

		var ips []string
		addrs, _ := iface.Addrs()

		for _, addr := range addrs {
			ips = append(ips, addr.String())
		}

		if len(ips) > 0 {
			macAddr.IPs = ips
		}
	}

	return macAddr
}

type MacIP struct {
	InterfaceName string   `json:"interface_name"`
	HardwareAddr  string   `json:"hardware_addr"`
	IPs           []string `json:"ips"`
}

func GetHeaders(r *http.Request) map[string]interface{} {
	// Request user agent
	userAgent := r.UserAgent()
	platform := strings.Split(r.Header.Get("sec-ch-ua"), ",")
	mobile := r.Header.Get("sec-ch-ua-mobile")
	operatingSystem := r.Header.Get("sec-ch-ua-platform")
	clientIP := r.RemoteAddr
	reqId := r.Header.Get("X-Session-Id")
	if reqId == "" {
		reqId = uuid.NewString()
	}

	macIp := getMACAndIP()
	return map[string]interface{}{
		"user_agent": userAgent,
		"Platform":   platform,
		"Mobile":     mobile,
		"OS":         operatingSystem,
		"client_ip":  clientIP,
		"request_id": reqId,
		"remote_ip":  r.RemoteAddr,
		"mac_ip":     macIp,
	}
}
