package deployer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	xerrors "github.com/pkg/errors"

	"certimate/internal/domain"
	"certimate/internal/pkg/utils/maps"
	xhttp "certimate/internal/utils/http"
)

// 定义访问url
const (
	Url = "https://unicloud-api.dcloud.net.cn/unicloud/api"
)

type UnicloudDeployer struct {
	option *DeployerOption
	infos  []string
}

func NewUnicloudDeployer(option *DeployerOption) (Deployer, error) {
	return &UnicloudDeployer{
		option: option,
		infos:  make([]string, 0),
	}, nil
}

func (d *UnicloudDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *UnicloudDeployer) GetInfos() []string {
	return d.infos
}

type unicloudData struct {
	SpaceId  string `json:"spaceId"`
	Provider string `json:"provider"`
	Domain   string `json:"domain"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
}

func (d *UnicloudDeployer) Deploy(ctx context.Context) error {
	access := &domain.UnicloudAccess{}
	d.infos = append(d.infos, toStr("Unicloud Access", d.option.Access))
	if err := json.Unmarshal([]byte(d.option.Access), access); err != nil {
		return xerrors.Wrap(err, "failed to get access")
	}

	data := &unicloudData{
		SpaceId:  maps.GetValueAsString(d.option.DeployConfig.Config, "spaceId"),
		Provider: maps.GetValueAsString(d.option.DeployConfig.Config, "provider"),
		Domain:   d.option.Domain,
		Cert:     url.QueryEscape(d.option.Certificate.Certificate),
		Key:      url.QueryEscape(d.option.Certificate.PrivateKey),
	}
	// 打印请求信息
	// d.infos = append(d.infos, toStr("Unicloud Request Data", data))
	// header 设置token，从配置中读取
	url := Url + "/host/create-domain-with-cert"
	d.infos = append(d.infos, toStr("Unicloud Request", url))
	body, _ := json.Marshal(data)
	resp, err := xhttp.Req(url, http.MethodPost, bytes.NewReader(body), map[string]string{
		"Content-Type":    "application/json",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
		"Token":           access.Token,
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"Accept":          "*/*",
		"Connection":      "keep-alive",
	})
	if err != nil {
		return xerrors.Wrap(err, "failed to send unicloud request")
	}

	d.infos = append(d.infos, toStr("Unicloud Response", string(resp)))

	return nil
}
