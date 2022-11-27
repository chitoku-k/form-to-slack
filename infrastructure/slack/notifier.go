package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chitoku-k/form-to-slack/service"
	"github.com/slack-go/slack"
)

type slackNotifier struct {
	webhookURL string
}

func NewSlackNotifier(webhookURL string) service.SlackNotifier {
	return &slackNotifier{
		webhookURL: webhookURL,
	}
}

func (s *slackNotifier) Do(ctx context.Context, message service.SlackMessage) error {
	timestamp := time.Now().Unix()
	author := message.Name
	if message.Email != "" {
		author += "（" + message.Email + "）"
	}

	err := slack.PostWebhookContext(ctx, s.webhookURL, &slack.WebhookMessage{
		Text: "New message has arrived:",
		Attachments: []slack.Attachment{
			{
				Title:      message.Subject,
				Text:       message.Body,
				Ts:         json.Number(fmt.Sprint(timestamp)),
				AuthorName: author,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to notify message: %w", err)
	}

	return nil
}
