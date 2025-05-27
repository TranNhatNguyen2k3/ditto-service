package graceful

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"ditto/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go/log"
)

const (
	TimeOutDefault  = 10 * time.Second
	DefaultWaitTime = 10 * time.Second
)

type service struct {
	currentStatus int
	waitTime      time.Duration
	timeout       time.Duration
	server        http.Server
}

type Service interface {
	Register(g *gin.Engine)
	StartServer(handler http.Handler, port string)
	Close(logger logger.Logger)
}

func NewService(opts ...Option) Service {
	o := &opt{waitTime: DefaultWaitTime, stopTimeout: TimeOutDefault}
	for _, opt := range opts {
		opt.apply(o)
	}
	return &service{
		currentStatus: http.StatusOK,
		waitTime:      o.waitTime,
		timeout:       o.stopTimeout,
	}
}

func (s *service) Register(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "GREEN")
	})
}

func (s *service) StartServer(handler http.Handler, port string) {
	s.server = http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: TimeOutDefault,
	}
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error(fmt.Errorf("failed to listen and serve from server: %v", err))
	}
}

func (s *service) stopServer(logger logger.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Error(fmt.Errorf("server shutdown error: %v", err))
		return
	}
	logger.Info("stop server success")
}

func (s *service) Close(logger logger.Logger) {
	logger.Info("set ping status to 503")
	s.currentStatus = http.StatusServiceUnavailable
	time.Sleep(s.waitTime)
	s.stopServer(logger)
	logger.Info("server exited...")
}

func (s *service) SignalStop(logger logger.Logger) {
	logger.Info("set ping status to 503")
	s.currentStatus = http.StatusServiceUnavailable
	time.Sleep(s.waitTime)
}
