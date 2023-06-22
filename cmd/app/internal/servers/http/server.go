package serverHTTP

// import (
// 	"go-flow2sqlite/cmd/app/internal/config"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/gorilla/mux"

// 	"github.com/sirupsen/logrus"
// )

// type responseWriter struct {
// 	http.ResponseWriter
// 	code int
// }

// type Server struct {
// 	router            *mux.Router
// 	timerUpdatedevice *time.Timer
// 	getDevices        chan int
// }

// func NewServer(cfg *config.Config) *Server {
// 	return &Server{
// 		router: mux.NewRouter(),
// 	}
// }

// func (s *Server) configureRouter() {
// 	s.router.Use(s.logRequest)
// 	s.router.HandleFunc("/api/v1/updatedevices", s.handleUpdateDevices()).Methods("GET") // Update cashe of devices from Mikrotik
// 	s.router.HandleFunc("/api/v1/devices", s.handleGetDevices()).Methods("GET")          // Get all devices from lease or ARP
// }

// func (s *Server) logRequest(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var host string
// 		var hostArr []string
// 		var ok bool
// 		hostArr, ok = r.Header["X-Forwarded-For"]
// 		if !ok {
// 			hostPort := r.RemoteAddr
// 			hostArr = strings.Split(hostPort, ":")
// 		}
// 		if len(hostArr) > 0 {
// 			host = hostArr[0]
// 		}
// 		logger := logrus.WithFields(logrus.Fields{
// 			"remote_addr": host,
// 			// "request_id":  r.Context().Value(ctxKeyRequestID),
// 		})
// 		logger.Infof("started %s %s", r.Method, r.RequestURI)
// 		start := time.Now()
// 		rw := &responseWriter{w, http.StatusOK}
// 		next.ServeHTTP(rw, r)

// 		var level logrus.Level
// 		switch {
// 		case rw.code >= 500:
// 			level = logrus.ErrorLevel
// 		case rw.code >= 400:
// 			level = logrus.WarnLevel
// 		default:
// 			level = logrus.InfoLevel
// 		}
// 		logger.Logf(
// 			level,
// 			"completed with %d %s in %v",
// 			rw.code,
// 			http.StatusText(rw.code),
// 			time.Since(start),
// 		)
// 	})
// }
