package deployer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

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
type unicloudToken struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type TokenResponse struct {
	Token string `json:"token"`
}
type HTTPError interface {
	error
	StatusCode() int
}

func (d *UnicloudDeployer) Deploy(ctx context.Context) error {
	access := &domain.UnicloudAccess{}
	d.infos = append(d.infos, toStr("Unicloud Access", d.option.Access))
	if err := json.Unmarshal([]byte(d.option.Access), access); err != nil {
		return xerrors.Wrap(err, "failed to get access")
	}

	// 打印日志
	d.infos = append(d.infos, toStr("Unicloud Access", access.Username))
	d.CheckToken(access)

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

// 检查token是否有效
func (d *UnicloudDeployer) CheckToken(access *domain.UnicloudAccess) error {
	// 如果 access.Token 为空，则获取新的 token
	// 打印日志
	d.infos = append(d.infos, toStr("开始检查Token", access.Token))
	if access.Token == "" {
		token, err := d.GetToken(access)
		if err != nil {
			return xerrors.Wrap(err, "failed to get new token")
		}
		access.Token = token
	}

	// 使用当前的 token 发送请求检查其是否有效
	url := Url + "/user/info"
	resp, err := xhttp.Req(url, http.MethodGet, nil, map[string]string{
		"Content-Type":    "application/json",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
		"Token":           access.Token,
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"Accept":          "*/*",
		"Connection":      "keep-alive",
	})
	if err != nil {
		// 打印日志
		d.infos = append(d.infos, toStr("CheckToken Error", err.Error()))
		// 检查是否是 HTTP 错误并且状态码为 401
		// 重新获取新的 token
		token, err := d.GetToken(access)
		if err != nil {
			return xerrors.Wrap(err, "failed to get new token")
		}
		access.Token = token
		// 重新发送请求检查新的 token 是否有效
		resp, err = xhttp.Req(url, http.MethodGet, nil, map[string]string{
			"Content-Type":    "application/json",
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
			"Token":           access.Token,
			"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
			"Accept":          "*/*",
			"Connection":      "keep-alive",
		})
		if err != nil {
			return xerrors.Wrap(err, "failed to send unicloud request with new token")
		}
	}
	// 处理响应
	d.infos = append(d.infos, toStr("CheckToken Response", string(resp)))
	d.infos = append(d.infos, toStr("token 检查通过 ", access.Token))

	return nil
}

// 如果token过期，重新获取token，请求地址：https://ulogin.cqsort.com/get_token?id=xxx&username=xxx&password=xxx
func (d *UnicloudDeployer) GetToken(access *domain.UnicloudAccess) (string, error) {
	url := "https://ulogin.cqsort.com/get_token"
	// data := &unicloudToken{
	// 	Id:       maps.GetValueAsString(d.option.DeployConfig.Config, "spaceId"),
	// 	Username: access.Username,
	// 	Password: access.Password,
	// }
	// body, _ := json.Marshal(data)
	data := fmt.Sprintf("?id=%s&username=%s&password=%s", maps.GetValueAsString(d.option.DeployConfig.Config, "spaceId"), access.Username, access.Password)
	// 打印请求信息
	d.infos = append(d.infos, toStr("Unicloud GetToken Request Data", data))
	resp, err := xhttp.Req(url+data, http.MethodGet, nil, map[string]string{
		"Content-Type":    "application/x-www-form-urlencoded",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"Accept":          "*/*",
		"Connection":      "keep-alive",
	}, xhttp.WithTimeout(2*time.Minute))
	if err != nil {
		// 打印日志
		d.infos = append(d.infos, toStr("GetToken Error", err.Error()))
		return "", xerrors.Wrap(err, "failed to send unicloud request")
	}
	// 定义一个 TokenResponse 结构体实例
	var tokenResp TokenResponse
	// 解析 JSON 数据
	err = json.Unmarshal(resp, &tokenResp)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return "", xerrors.Wrap(err, "failed to parse JSON")
	}

	// 将解析出的 token 值赋值给 access.Token
	// access.Token = tokenResp.Token
	d.infos = append(d.infos, toStr("Unicloud Response", string(resp)))
	return tokenResp.Token, nil
}
