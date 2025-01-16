package bark

import (
	"context"
	"errors"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/bark"

	"certimate/internal/pkg/core/notifier"
)

type BarkNotifierConfig struct {
	ServerUrl string `json:"serverUrl"`
	DeviceKey string `json:"deviceKey"`
}

type BarkNotifier struct {
	config *BarkNotifierConfig
}

var _ notifier.Notifier = (*BarkNotifier)(nil)

func New(config *BarkNotifierConfig) (*BarkNotifier, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	return &BarkNotifier{
		config: config,
	}, nil
}

func (n *BarkNotifier) Notify(ctx context.Context, subject string, message string) (res *notifier.NotifyResult, err error) {
	var srv notify.Notifier
	if n.config.ServerUrl == "" {
		srv = bark.New(n.config.DeviceKey)
	} else {
		srv = bark.NewWithServers(n.config.DeviceKey, n.config.ServerUrl)
	}

	err = srv.Send(ctx, subject, message)
	if err != nil {
		return nil, err
	}

	return &notifier.NotifyResult{}, nil
}
