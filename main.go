package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/chitoku-k/form-to-slack/application/server"
	"github.com/chitoku-k/form-to-slack/infrastructure/config"
	"github.com/chitoku-k/form-to-slack/infrastructure/slack"
	"github.com/chitoku-k/form-to-slack/service"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	env, err := config.Get()
	if err != nil {
		logrus.Fatalf("Failed to initialize config: %v", err)
	}

	slackNotifier := slack.NewSlackNotifier(env.SlackWebhookURL)
	slackService := service.NewSlackService(slackNotifier)
	engine := server.NewEngine(env.Port, env.TLSCert, env.TLSKey, env.AllowedOrigins, env.ReCaptchaSecret, slackService)
	err = engine.Start(ctx)
	if err != nil {
		logrus.Fatalf("Failed to start web server: %v", err)
	}
}
