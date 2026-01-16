package deployer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"certimate/internal/domain"
	"certimate/internal/pkg/utils/maps"

	xerrors "github.com/pkg/errors"

	xhttp "certimate/internal/utils/http"
)

// 定义访问url
const (
	Url = "https://ulogin.cqsort.com"
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
	Username string `json:"username"`
	Password string `json:"password"`
	Id       string `json:"id"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
}

func (d *UnicloudDeployer) Deploy(ctx context.Context) error {
	access := &domain.UnicloudAccess{}
	d.infos = append(d.infos, toStr("Unicloud Access", d.option.Access))
	if err := json.Unmarshal([]byte(d.option.Access), access); err != nil {
		return xerrors.Wrap(err, "failed to get access")
	}

	// 打印日志
	d.infos = append(d.infos, toStr("Unicloud Access", access.Username))

	data := &unicloudData{
		Username: access.Username,
		Password: access.Password,
		Id:       maps.GetValueAsString(d.option.DeployConfig.Config, "spaceId"),
		Cert:     d.option.Certificate.Certificate,
		Key:      d.option.Certificate.PrivateKey,
	}
	
	// 打印请求信息
	url := Url + "/cert"
	d.infos = append(d.infos, toStr("Unicloud Request", url))
	body, _ := json.Marshal(data)
	resp, err := xhttp.Req(url, http.MethodPost, bytes.NewReader(body), map[string]string{
		"Content-Type":    "application/json",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
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
