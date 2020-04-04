package main

import (
	"github.com/chitoku-k/form-to-slack/application/server"
	"github.com/chitoku-k/form-to-slack/infrastructure/config"
	"github.com/chitoku-k/form-to-slack/infrastructure/slack"
	"github.com/chitoku-k/form-to-slack/service"
	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/pkg/errors"
)

func main() {
	env, err := config.Get()
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize config"))
	}

	recaptcha.Init(env.ReCaptchaSecret)
	slackNotifier := slack.NewSlackNotifier(env)
	slackService := service.NewSlackService(slackNotifier)
	engine := server.NewEngine(env, slackService)
	engine.Start()
}
