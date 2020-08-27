package main

import (
	"fmt"

	"github.com/chitoku-k/form-to-slack/application/server"
	"github.com/chitoku-k/form-to-slack/infrastructure/config"
	"github.com/chitoku-k/form-to-slack/infrastructure/slack"
	"github.com/chitoku-k/form-to-slack/service"
	"github.com/dpapathanasiou/go-recaptcha"
)

func main() {
	env, err := config.Get()
	if err != nil {
		panic(fmt.Errorf("failed to initialize config: %w", err))
	}

	recaptcha.Init(env.ReCaptchaSecret)
	slackNotifier := slack.NewSlackNotifier(env)
	slackService := service.NewSlackService(slackNotifier)
	engine := server.NewEngine(env, slackService)
	engine.Start()
}
