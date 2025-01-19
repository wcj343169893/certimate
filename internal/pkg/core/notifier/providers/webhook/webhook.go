package webhook

import (
	"context"
	"errors"

	"github.com/nikoksr/notify/service/http"

	"certimate/internal/pkg/core/notifier"
)

type WebhookNotifierConfig struct {
	Url string `json:"url"`
}

type WebhookNotifier struct {
	config *WebhookNotifierConfig
}

var _ notifier.Notifier = (*WebhookNotifier)(nil)

func New(config *WebhookNotifierConfig) (*WebhookNotifier, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	return &WebhookNotifier{
		config: config,
	}, nil
}

func (n *WebhookNotifier) Notify(ctx context.Context, subject string, message string) (res *notifier.NotifyResult, err error) {
	srv := http.New()

	srv.AddReceiversURLs(n.config.Url)

	err = srv.Send(ctx, subject, message)
	if err != nil {
		return nil, err
	}

	return &notifier.NotifyResult{}, nil
}
