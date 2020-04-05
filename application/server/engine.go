package server

import (
	"net/http"
	"strings"

	"github.com/chitoku-k/form-to-slack/infrastructure/config"
	"github.com/chitoku-k/form-to-slack/service"
	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type engine struct {
	Environment  config.Environment
	SlackService service.SlackService
}

type Engine interface {
	Start()
}

func NewEngine(
	environment config.Environment,
	slackService service.SlackService,
) Engine {
	return &engine{
		Environment:  environment,
		SlackService: slackService,
	}
}

func (e *engine) Start() {
	router := gin.Default()

	if len(e.Environment.AllowedOrigins) > 0 {
		router.Use(cors.New(cors.Config{
			AllowOrigins: strings.Split(e.Environment.AllowedOrigins, " "),
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
		result, err := recaptcha.Confirm(c.ClientIP(), c.PostForm("g-recaptcha-response"))
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

		err = e.SlackService.Send(message)
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

	router.Run(":" + e.Environment.Port)
}
