package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/chitoku-k/form-to-slack/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type engine struct {
	Port            string
	CertFile        string
	KeyFile         string
	AllowedOrigins  string
	ReCaptchaSecret string
	SlackService    service.SlackService
}

type Engine interface {
	Start(ctx context.Context) error
}

func NewEngine(
	port string,
	certFile string,
	keyFile string,
	allowedOrigins string,
	reCaptchaSecret string,
	slackService service.SlackService,
) Engine {
	return &engine{
		Port:            port,
		CertFile:        certFile,
		KeyFile:         keyFile,
		AllowedOrigins:  allowedOrigins,
		ReCaptchaSecret: reCaptchaSecret,
		SlackService:    slackService,
	}
}

func (e *engine) Start(ctx context.Context) error {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/healthz"},
	}))

	if len(e.AllowedOrigins) > 0 {
		router.Use(cors.New(cors.Config{
			AllowOrigins: strings.Split(e.AllowedOrigins, " "),
		}))
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
		})
	})

	router.Any("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	router.POST("/", func(c *gin.Context) {
		result, err := e.verifyReCaptcha(c, c.PostForm("g-recaptcha-response"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
			})
			return
		}
		if !result {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
			})
			return
		}

		var message service.SlackMessage
		err = c.ShouldBind(&message)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}

		err = e.SlackService.Send(ctx, message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
		})
	})

	server := http.Server{
		Addr:    net.JoinHostPort("", e.Port),
		Handler: router,
	}

	var eg errgroup.Group
	eg.Go(func() error {
		<-ctx.Done()
		return server.Shutdown(context.Background())
	})

	var err error
	if e.CertFile != "" && e.KeyFile != "" {
		server.TLSConfig = &tls.Config{
			GetCertificate: e.getCertificate,
		}
		err = server.ListenAndServeTLS("", "")
	} else {
		err = server.ListenAndServe()
	}

	if err == http.ErrServerClosed {
		return eg.Wait()
	}

	return err
}

func (e *engine) getCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(e.CertFile, e.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get certificate: %w", err)
	}

	return &cert, nil
}
