package unicloud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	xerrors "github.com/pkg/errors"

	"certimate/internal/pkg/core/deployer"
	xhttp "certimate/internal/utils/http"
)

// 定义访问url
const (
	Url = "https://unicloud.dcloud.net.cn"
)

type UnicloudDeployerConfig struct {
	SpaceId  string `json:"spaceId"`
	Domain   string `json:"domain"`
	Provider string `json:"provider"`
	Token    string `json:"token"`
}

type UnicloudDeployer struct {
	config *UnicloudDeployerConfig
	logger deployer.Logger
}

var _ deployer.Deployer = (*UnicloudDeployer)(nil)

func New(config *UnicloudDeployerConfig) (*UnicloudDeployer, error) {
	return NewWithLogger(config, deployer.NewNilLogger())
}

func NewWithLogger(config *UnicloudDeployerConfig, logger deployer.Logger) (*UnicloudDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	return &UnicloudDeployer{
		config: config,
		logger: logger,
	}, nil
}

type unicloudData struct {
	SpaceId  string `json:"spaceId"`
	Provider string `json:"provider"`
	Domain   string `json:"domain"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
}

func (d *UnicloudDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	data := &unicloudData{
		SpaceId:  d.config.SpaceId,
		Provider: d.config.Provider,
		Domain:   d.config.Domain,
		Cert:     url.QueryEscape(certPem),
		Key:      url.QueryEscape(privkeyPem),
	}
	// header 设置token，从配置中读取

	body, _ := json.Marshal(data)
	resp, err := xhttp.Req(Url, http.MethodPost, bytes.NewReader(body), map[string]string{
		"Content-Type":    "application/json",
		"token":           d.config.Token,
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"Accept":          "*/*",
		"Connection":      "keep-alive",
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to send webhook request")
	}

	d.logger.Logt("Unicloud Response", string(resp))

	return &deployer.DeployResult{
		DeploymentData: map[string]any{
			"responseText": string(resp),
		},
	}, nil
}
