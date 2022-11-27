package service

import (
	"context"
)

type SlackMessage struct {
	Name    string `form:"name" binding:"required"`
	Email   string `form:"email"`
	Subject string `form:"subject" binding:"required"`
	Body    string `form:"body" binding:"required"`
}

type SlackNotifier interface {
	Do(ctx context.Context, message SlackMessage) error
}

type slackService struct {
	Notifier SlackNotifier
}

type SlackService interface {
	Send(ctx context.Context, message SlackMessage) error
}

func NewSlackService(notifier SlackNotifier) SlackService {
	return &slackService{
		Notifier: notifier,
	}
}

func (s *slackService) Send(ctx context.Context, message SlackMessage) error {
	return s.Notifier.Do(ctx, message)
}
