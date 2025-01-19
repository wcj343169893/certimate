package deployer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pocketbase/pocketbase/models"

	"certimate/internal/applicant"
	"certimate/internal/domain"
	"certimate/internal/utils/app"
)

const (
	targetAliyunOSS      = "aliyun-oss"
	targetAliyunCDN      = "aliyun-cdn"
	targetAliyunDCDN     = "aliyun-dcdn"
	targetAliyunCLB      = "aliyun-clb"
	targetAliyunALB      = "aliyun-alb"
	targetAliyunNLB      = "aliyun-nlb"
	targetTencentCDN     = "tencent-cdn"
	targetTencentECDN    = "tencent-ecdn"
	targetTencentCLB     = "tencent-clb"
	targetTencentCOS     = "tencent-cos"
	targetTencentTEO     = "tencent-teo"
	targetHuaweiCloudCDN = "huaweicloud-cdn"
	targetHuaweiCloudELB = "huaweicloud-elb"
	targetBaiduCloudCDN  = "baiducloud-cdn"
	targetVolcEngineLive = "volcengine-live"
	targetVolcEngineCDN  = "volcengine-cdn"
	targetBytePlusCDN    = "byteplus-cdn"
	targetQiniuCdn       = "qiniu-cdn"
	targetDogeCloudCdn   = "dogecloud-cdn"
	targetLocal          = "local"
	targetSSH            = "ssh"
	targetWebhook        = "webhook"
	targetK8sSecret      = "k8s-secret"
)

type DeployerOption struct {
	DomainId     string                `json:"domainId"`
	Domain       string                `json:"domain"`
	Access       string                `json:"access"`
	AccessRecord *models.Record        `json:"-"`
	DeployConfig domain.DeployConfig   `json:"deployConfig"`
	Certificate  applicant.Certificate `json:"certificate"`
	Variables    map[string]string     `json:"variables"`
}

type Deployer interface {
	Deploy(ctx context.Context) error
	GetInfos() []string
	GetID() string
}

func Gets(record *models.Record, cert *applicant.Certificate) ([]Deployer, error) {
	rs := make([]Deployer, 0)
	if record.GetString("deployConfig") == "" {
		return rs, nil
	}

	deployConfigs := make([]domain.DeployConfig, 0)

	err := record.UnmarshalJSONField("deployConfig", &deployConfigs)
	if err != nil {
		return nil, fmt.Errorf("解析部署配置失败: %w", err)
	}

	if len(deployConfigs) == 0 {
		return rs, nil
	}

	for _, deployConfig := range deployConfigs {
		deployer, err := getWithDeployConfig(record, cert, deployConfig)
		if err != nil {
			return nil, err
		}

		rs = append(rs, deployer)
	}

	return rs, nil
}

func getWithDeployConfig(record *models.Record, cert *applicant.Certificate, deployConfig domain.DeployConfig) (Deployer, error) {
	access, err := app.GetApp().Dao().FindRecordById("access", deployConfig.Access)
	if err != nil {
		return nil, fmt.Errorf("access record not found: %w", err)
	}

	option := &DeployerOption{
		DomainId:     record.Id,
		Domain:       record.GetString("domain"),
		Access:       access.GetString("config"),
		AccessRecord: access,
		DeployConfig: deployConfig,
	}
	if cert != nil {
		option.Certificate = *cert
	} else {
		option.Certificate = applicant.Certificate{
			Certificate: record.GetString("certificate"),
			PrivateKey:  record.GetString("privateKey"),
		}
	}

	switch deployConfig.Type {
	case targetAliyunOSS:
		return NewAliyunOSSDeployer(option)
	case targetAliyunCDN:
		return NewAliyunCDNDeployer(option)
	case targetAliyunDCDN:
		return NewAliyunDCDNDeployer(option)
	case targetAliyunCLB:
		return NewAliyunCLBDeployer(option)
	case targetAliyunALB:
		return NewAliyunALBDeployer(option)
	case targetAliyunNLB:
		return NewAliyunNLBDeployer(option)
	case targetTencentCDN:
		return NewTencentCDNDeployer(option)
	case targetTencentECDN:
		return NewTencentECDNDeployer(option)
	case targetTencentCLB:
		return NewTencentCLBDeployer(option)
	case targetTencentCOS:
		return NewTencentCOSDeployer(option)
	case targetTencentTEO:
		return NewTencentTEODeployer(option)
	case targetHuaweiCloudCDN:
		return NewHuaweiCloudCDNDeployer(option)
	case targetHuaweiCloudELB:
		return NewHuaweiCloudELBDeployer(option)
	case targetBaiduCloudCDN:
		return NewBaiduCloudCDNDeployer(option)
	case targetQiniuCdn:
		return NewQiniuCDNDeployer(option)
	case targetDogeCloudCdn:
		return NewDogeCloudCDNDeployer(option)
	case targetLocal:
		return NewLocalDeployer(option)
	case targetSSH:
		return NewSSHDeployer(option)
	case targetWebhook:
		return NewWebhookDeployer(option)
	case targetK8sSecret:
		return NewK8sSecretDeployer(option)
	case targetVolcEngineLive:
		return NewVolcengineLiveDeployer(option)
	case targetVolcEngineCDN:
		return NewVolcengineCDNDeployer(option)
	case targetBytePlusCDN:
		return NewByteplusCDNDeployer(option)
	}
	return nil, errors.New("unsupported deploy target")
}

func toStr(tag string, data any) string {
	if data == nil {
		return tag
	}
	byts, _ := json.Marshal(data)
	return tag + "：" + string(byts)
}
