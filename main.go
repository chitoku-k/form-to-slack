package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/chitoku-k/form-to-slack/application/server"
	"github.com/chitoku-k/form-to-slack/infrastructure/config"
	"github.com/chitoku-k/form-to-slack/infrastructure/slack"
	"github.com/chitoku-k/form-to-slack/service"
	"github.com/spf13/pflag"
)

var (
	signals = []os.Signal{os.Interrupt}
	name    = "form-to-slack"
	version = "v0.0.0-dev"

	flagversion = pflag.BoolP("version", "V", false, "show version")
)

func main() {
	pflag.Parse()
	if *flagversion {
		fmt.Println(name, version)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()

	env, err := config.Get()
	if err != nil {
		slog.Error("Failed to initialize config", slog.Any("err", err))
		os.Exit(1)
	}

	slackNotifier := slack.NewSlackNotifier(env.SlackWebhookURL)
	slackService := service.NewSlackService(slackNotifier)
	engine := server.NewEngine(env.Port, env.TLSCert, env.TLSKey, env.AllowedOrigins, env.ReCaptchaURL, env.ReCaptchaSecret, slackService)
	err = engine.Start(ctx)
	if err != nil {
		slog.Error("Failed to start web server", slog.Any("err", err))
		os.Exit(1)
	}
}
