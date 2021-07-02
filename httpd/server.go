package httpd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/booking/api"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/sirupsen/logrus"

	"github.com/booking/config"
)

const (
	// WriteTimeout defines max write timeout
	WriteTimeout = 15 * time.Second
	// ReadTimeout defines max read timeout
	ReadTimeout = 15 * time.Second
)

// Server defines a HTTP server
type Server struct {
	config    *config.Config
	container *restful.Container
	httpd     *http.Server
	logger    *logrus.Entry

	stop chan interface{}
	done chan interface{}
}

// NewServer return a configured HTTP server
func NewServer(c *config.Config, l *logrus.Entry) *Server {
	return &Server{
		config:    c,
		container: restful.NewContainer(),
		httpd: &http.Server{
			Addr:         fmt.Sprintf(":%v", c.ListenPort),
			WriteTimeout: WriteTimeout,
			ReadTimeout:  ReadTimeout,
		},
		logger: l,
		stop:   make(chan interface{}),
		done:   make(chan interface{}),
	}
}

// Add adds restful.WebService to restful.Container
func (s *Server) Add(svc *restful.WebService) {
	s.container.Add(svc)
}

// Start configures APIDocs endpoints and starts HTTP server in background
func (s *Server) Start(parentCtx context.Context) {
	s.logger.WithField("address", s.httpd.Addr).
		Info("starting httpd")

	ctx, cancel := context.WithCancel(parentCtx)

	s.httpd.BaseContext = func(l net.Listener) context.Context {
		return ctx
	}

	s.httpd.Handler = s.container

	rsc := restfulspec.Config{
		WebServices:                   s.container.RegisteredWebServices(),
		APIPath:                       "/api/swagger.json",
		PostBuildSwaggerObjectHandler: api.EnrichSwaggerObject}
	s.container.Add(restfulspec.NewOpenAPIService(rsc))

	s.container.ServeMux.Handle("/api/apidocs/", http.StripPrefix("/api/apidocs/", http.FileServer(http.Dir(s.config.SwaggerDistPath))))

	// Optionally, you may need to enable CORS for the UI to work.
	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      s.container,
	}
	s.container.Filter(cors.Filter)

	go func() {
		defer func() {
			cancel()
			close(s.done)
		}()

		s.logger.Info(fmt.Sprintf("Get the API using http://%v:%v/api/swagger.json", s.config.Hostname, s.config.ListenPort))
		s.logger.Info(fmt.Sprintf("Open Swagger UI using http://%v:%v/api/apidocs/?url=http://%v:%v/api/swagger.json", s.config.Hostname, s.config.ListenPort, s.config.Hostname, s.config.ListenPort))

		if err := s.httpd.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.WithError(err).
				Info("httpd server error")
			return
		}

		<-s.stop
	}()
}

// Stop gracefully shutdown webserver and closes stop channel
func (s *Server) Stop() {
	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.logger.WithError(s.httpd.Shutdown(stopCtx)).
		Info("terminating httpd")
	close(s.stop)
}

// Shutdown returns done channel that signals webserver has shutdown
func (s *Server) Shutdown() <-chan interface{} {
	return s.done
}
