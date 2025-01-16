package volcenginecdn

import (
	"context"
	"errors"
	"fmt"
	"strings"

	xerrors "github.com/pkg/errors"
	veCdn "github.com/volcengine/volc-sdk-golang/service/cdn"

	"certimate/internal/pkg/core/deployer"
	"certimate/internal/pkg/core/uploader"
	providerCdn "certimate/internal/pkg/core/uploader/providers/volcengine-cdn"
)

type VolcEngineCDNDeployerConfig struct {
	// 火山引擎 AccessKey。
	AccessKey string `json:"accessKey"`
	// 火山引擎 SecretKey。
	SecretKey string `json:"secretKey"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type VolcEngineCDNDeployer struct {
	config      *VolcEngineCDNDeployerConfig
	logger      deployer.Logger
	sdkClient   *veCdn.CDN
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*VolcEngineCDNDeployer)(nil)

func New(config *VolcEngineCDNDeployerConfig) (*VolcEngineCDNDeployer, error) {
	return NewWithLogger(config, deployer.NewNilLogger())
}

func NewWithLogger(config *VolcEngineCDNDeployerConfig, logger deployer.Logger) (*VolcEngineCDNDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client := veCdn.NewInstance()
	client.Client.SetAccessKey(config.AccessKey)
	client.Client.SetSecretKey(config.SecretKey)

	uploader, err := providerCdn.New(&providerCdn.VolcEngineCDNUploaderConfig{
		AccessKeyId:     config.AccessKey,
		AccessKeySecret: config.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &VolcEngineCDNDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *VolcEngineCDNDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 CDN
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	domains := make([]string, 0)
	if strings.HasPrefix(d.config.Domain, "*.") {
		// 获取指定证书可关联的域名
		// REF: https://www.volcengine.com/docs/6454/125711
		describeCertConfigReq := &veCdn.DescribeCertConfigRequest{
			CertId: upres.CertId,
		}
		describeCertConfigResp, err := d.sdkClient.DescribeCertConfig(describeCertConfigReq)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.DescribeCertConfig'")
		}

		if describeCertConfigResp.Result.CertNotConfig != nil {
			for i := range describeCertConfigResp.Result.CertNotConfig {
				domains = append(domains, describeCertConfigResp.Result.CertNotConfig[i].Domain)
			}
		}

		if describeCertConfigResp.Result.OtherCertConfig != nil {
			for i := range describeCertConfigResp.Result.OtherCertConfig {
				domains = append(domains, describeCertConfigResp.Result.OtherCertConfig[i].Domain)
			}
		}

		if len(domains) == 0 {
			if len(describeCertConfigResp.Result.SpecifiedCertConfig) > 0 {
				// 所有可关联的域名都配置了该证书，跳过部署
			} else {
				return nil, xerrors.New("domain not found")
			}
		}
	} else {
		domains = append(domains, d.config.Domain)
	}

	if len(domains) > 0 {
		var errs []error

		for _, domain := range domains {
			// 关联证书与加速域名
			// REF: https://www.volcengine.com/docs/6454/125712
			batchDeployCertReq := &veCdn.BatchDeployCertRequest{
				CertId: upres.CertId,
				Domain: domain,
			}
			batchDeployCertResp, err := d.sdkClient.BatchDeployCert(batchDeployCertReq)
			if err != nil {
				errs = append(errs, err)
			} else {
				d.logger.Logt(fmt.Sprintf("已关联证书到域名 %s", domain), batchDeployCertResp)
			}
		}

		if len(errs) > 0 {
			return nil, errors.Join(errs...)
		}
	}

	return &deployer.DeployResult{}, nil
}
