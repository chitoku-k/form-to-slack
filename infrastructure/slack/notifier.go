package slack

import (
	"errors"
	"strings"
	"time"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/chitoku-k/form-to-slack/infrastructure/config"
	"github.com/chitoku-k/form-to-slack/service"
)

type slackNotifier struct {
	Environment config.Environment
}

func NewSlackNotifier(environment config.Environment) service.SlackNotifier {
	return &slackNotifier{
		Environment: environment,
	}
}

func (s *slackNotifier) Do(message service.SlackMessage) error {
	timestamp := time.Now().Unix()
	author := message.Name
	if message.Email != "" {
		author += "（" + message.Email + "）"
	}

	errs := slack.Send(s.Environment.SlackWebhookURL, "", slack.Payload{
		Text: "New message has arrived:",
		Attachments: []slack.Attachment{
			{
				Title:      &message.Subject,
				Text:       &message.Body,
				Timestamp:  &timestamp,
				AuthorName: &author,
			},
		},
	})

	if len(errs) > 0 {
		var msgs []string
		for _, err := range errs {
			msgs = append(msgs, err.Error())
		}
		return errors.New("failed to notify message: " + strings.Join(msgs, ", "))
	}

	return nil
}
