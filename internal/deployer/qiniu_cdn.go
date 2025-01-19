package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	xerrors "github.com/pkg/errors"
	"github.com/qiniu/go-sdk/v7/auth"

	"certimate/internal/domain"
	"certimate/internal/pkg/core/uploader"
	uploaderQiniu "certimate/internal/pkg/core/uploader/providers/qiniu-sslcert"
	qiniuEx "certimate/internal/pkg/vendors/qiniu-sdk"
)

type QiniuCDNDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClient   *qiniuEx.Client
	sslUploader uploader.Uploader
}

func NewQiniuCDNDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.QiniuAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}

	client, err := (&QiniuCDNDeployer{}).createSdkClient(
		access.AccessKey,
		access.SecretKey,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploaderQiniu.New(&uploaderQiniu.QiniuSSLCertUploaderConfig{
		AccessKey: access.AccessKey,
		SecretKey: access.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &QiniuCDNDeployer{
		option:      option,
		infos:       make([]string, 0),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *QiniuCDNDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *QiniuCDNDeployer) GetInfos() []string {
	return d.infos
}

func (d *QiniuCDNDeployer) Deploy(ctx context.Context) error {
	// 上传证书
	upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	// 在七牛 CDN 中泛域名表示为 .example.com，需去除前缀星号
	domain := d.option.DeployConfig.GetConfigAsString("domain")
	if strings.HasPrefix(domain, "*") {
		domain = strings.TrimPrefix(domain, "*")
	}

	// 获取域名信息
	// REF: https://developer.qiniu.com/fusion/4246/the-domain-name
	getDomainInfoResp, err := d.sdkClient.GetDomainInfo(domain)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'cdn.GetDomainInfo'")
	}

	d.infos = append(d.infos, toStr("已获取域名信息", getDomainInfoResp))

	// 判断域名是否已启用 HTTPS。如果已启用，修改域名证书；否则，启用 HTTPS
	// REF: https://developer.qiniu.com/fusion/4246/the-domain-name
	if getDomainInfoResp.Https != nil && getDomainInfoResp.Https.CertID != "" {
		modifyDomainHttpsConfResp, err := d.sdkClient.ModifyDomainHttpsConf(domain, upres.CertId, getDomainInfoResp.Https.ForceHttps, getDomainInfoResp.Https.Http2Enable)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'cdn.ModifyDomainHttpsConf'")
		}

		d.infos = append(d.infos, toStr("已修改域名证书", modifyDomainHttpsConfResp))
	} else {
		enableDomainHttpsResp, err := d.sdkClient.EnableDomainHttps(domain, upres.CertId, true, true)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'cdn.EnableDomainHttps'")
		}

		d.infos = append(d.infos, toStr("已将域名升级为 HTTPS", enableDomainHttpsResp))
	}

	return nil
}

func (u *QiniuCDNDeployer) createSdkClient(accessKey, secretKey string) (*qiniuEx.Client, error) {
	credential := auth.New(accessKey, secretKey)
	client := qiniuEx.NewClient(credential)
	return client, nil
}
