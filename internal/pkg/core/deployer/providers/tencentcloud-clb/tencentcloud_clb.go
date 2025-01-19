package tencentcloudclb

import (
	"context"
	"errors"
	"fmt"

	xerrors "github.com/pkg/errors"
	tcClb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"certimate/internal/pkg/core/deployer"
	"certimate/internal/pkg/core/uploader"
	providerSsl "certimate/internal/pkg/core/uploader/providers/tencentcloud-ssl"
)

type TencentCloudCLBDeployerConfig struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secretId"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 腾讯云地域。
	Region string `json:"region"`
	// 部署资源类型。
	ResourceType DeployResourceType `json:"resourceType"`
	// 负载均衡器 ID。
	// 部署资源类型为 [DEPLOY_RESOURCE_SSLDEPLOY]、[DEPLOY_RESOURCE_LOADBALANCER]、[DEPLOY_RESOURCE_RULEDOMAIN] 时必填。
	LoadbalancerId string `json:"loadbalancerId,omitempty"`
	// 负载均衡监听 ID。
	// 部署资源类型为 [DEPLOY_RESOURCE_SSLDEPLOY]、[DEPLOY_RESOURCE_LOADBALANCER]、[DEPLOY_RESOURCE_LISTENER]、[DEPLOY_RESOURCE_RULEDOMAIN] 时必填。
	ListenerId string `json:"listenerId,omitempty"`
	// SNI 域名或七层转发规则域名（支持泛域名）。
	// 部署资源类型为 [DEPLOY_RESOURCE_SSLDEPLOY] 时选填；部署资源类型为 [DEPLOY_RESOURCE_RULEDOMAIN] 时必填。
	Domain string `json:"domain,omitempty"`
}

type TencentCloudCLBDeployer struct {
	config      *TencentCloudCLBDeployerConfig
	logger      deployer.Logger
	sdkClients  *wSdkClients
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*TencentCloudCLBDeployer)(nil)

type wSdkClients struct {
	ssl *tcSsl.Client
	clb *tcClb.Client
}

func New(config *TencentCloudCLBDeployerConfig) (*TencentCloudCLBDeployer, error) {
	return NewWithLogger(config, deployer.NewNilLogger())
}

func NewWithLogger(config *TencentCloudCLBDeployerConfig, logger deployer.Logger) (*TencentCloudCLBDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	clients, err := createSdkClients(config.SecretId, config.SecretKey, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk clients")
	}

	uploader, err := providerSsl.New(&providerSsl.TencentCloudSSLUploaderConfig{
		SecretId:  config.SecretId,
		SecretKey: config.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &TencentCloudCLBDeployer{
		logger:      logger,
		config:      config,
		sdkClients:  clients,
		sslUploader: uploader,
	}, nil
}

func (d *TencentCloudCLBDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	// 根据部署资源类型决定部署方式
	switch d.config.ResourceType {
	case DEPLOY_RESOURCE_USE_SSLDEPLOY:
		if err := d.deployToInstanceUseSsl(ctx, upres.CertId); err != nil {
			return nil, err
		}

	case DEPLOY_RESOURCE_LOADBALANCER:
		if err := d.deployToLoadbalancer(ctx, upres.CertId); err != nil {
			return nil, err
		}

	case DEPLOY_RESOURCE_LISTENER:
		if err := d.deployToListener(ctx, upres.CertId); err != nil {
			return nil, err
		}

	case DEPLOY_RESOURCE_RULEDOMAIN:
		if err := d.deployToRuleDomain(ctx, upres.CertId); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported resource type: %s", d.config.ResourceType)
	}

	return &deployer.DeployResult{}, nil
}

func (d *TencentCloudCLBDeployer) deployToInstanceUseSsl(ctx context.Context, cloudCertId string) error {
	if d.config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}
	if d.config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}

	// 证书部署到 CLB 实例
	// REF: https://cloud.tencent.com/document/product/400/91667
	deployCertificateInstanceReq := tcSsl.NewDeployCertificateInstanceRequest()
	deployCertificateInstanceReq.CertificateId = common.StringPtr(cloudCertId)
	deployCertificateInstanceReq.ResourceType = common.StringPtr("clb")
	deployCertificateInstanceReq.Status = common.Int64Ptr(1)
	if d.config.Domain == "" {
		// 未开启 SNI，只需指定到监听器
		deployCertificateInstanceReq.InstanceIdList = common.StringPtrs([]string{fmt.Sprintf("%s|%s", d.config.LoadbalancerId, d.config.ListenerId)})
	} else {
		// 开启 SNI，需指定到域名（支持泛域名）
		deployCertificateInstanceReq.InstanceIdList = common.StringPtrs([]string{fmt.Sprintf("%s|%s|%s", d.config.LoadbalancerId, d.config.ListenerId, d.config.Domain)})
	}
	deployCertificateInstanceResp, err := d.sdkClients.ssl.DeployCertificateInstance(deployCertificateInstanceReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'ssl.DeployCertificateInstance'")
	}

	d.logger.Logt("已部署证书到云资源实例", deployCertificateInstanceResp.Response)

	return nil
}

func (d *TencentCloudCLBDeployer) deployToLoadbalancer(ctx context.Context, cloudCertId string) error {
	if d.config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}

	listenerIds := make([]string, 0)

	// 查询监听器列表
	// REF: https://cloud.tencent.com/document/api/214/30686
	describeListenersReq := tcClb.NewDescribeListenersRequest()
	describeListenersReq.LoadBalancerId = common.StringPtr(d.config.LoadbalancerId)
	describeListenersResp, err := d.sdkClients.clb.DescribeListeners(describeListenersReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'clb.DescribeListeners'")
	} else {
		if describeListenersResp.Response.Listeners != nil {
			for _, listener := range describeListenersResp.Response.Listeners {
				if listener.Protocol == nil || (*listener.Protocol != "HTTPS" && *listener.Protocol != "TCP_SSL" && *listener.Protocol != "QUIC") {
					continue
				}

				listenerIds = append(listenerIds, *listener.ListenerId)
			}
		}
	}

	d.logger.Logt("已查询到负载均衡器下的监听器", listenerIds)

	// 批量更新监听器证书
	if len(listenerIds) > 0 {
		var errs []error

		for _, listenerId := range listenerIds {
			if err := d.modifyListenerCertificate(ctx, d.config.LoadbalancerId, listenerId, cloudCertId); err != nil {
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			return errors.Join(errs...)
		}
	}

	return nil
}

func (d *TencentCloudCLBDeployer) deployToListener(ctx context.Context, cloudCertId string) error {
	if d.config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}
	if d.config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}

	// 更新监听器证书
	if err := d.modifyListenerCertificate(ctx, d.config.LoadbalancerId, d.config.ListenerId, cloudCertId); err != nil {
		return err
	}

	return nil
}

func (d *TencentCloudCLBDeployer) deployToRuleDomain(ctx context.Context, cloudCertId string) error {
	if d.config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}
	if d.config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}
	if d.config.Domain == "" {
		return errors.New("config `domain` is required")
	}

	// 修改负载均衡七层监听器转发规则的域名级别属性
	// REF: https://cloud.tencent.com/document/api/214/38092
	modifyDomainAttributesReq := tcClb.NewModifyDomainAttributesRequest()
	modifyDomainAttributesReq.LoadBalancerId = common.StringPtr(d.config.LoadbalancerId)
	modifyDomainAttributesReq.ListenerId = common.StringPtr(d.config.ListenerId)
	modifyDomainAttributesReq.Domain = common.StringPtr(d.config.Domain)
	modifyDomainAttributesReq.Certificate = &tcClb.CertificateInput{
		SSLMode: common.StringPtr("UNIDIRECTIONAL"),
		CertId:  common.StringPtr(cloudCertId),
	}
	modifyDomainAttributesResp, err := d.sdkClients.clb.ModifyDomainAttributes(modifyDomainAttributesReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'clb.ModifyDomainAttributes'")
	}

	d.logger.Logt("已修改七层监听器转发规则的域名级别属性", modifyDomainAttributesResp.Response)

	return nil
}

func (d *TencentCloudCLBDeployer) modifyListenerCertificate(ctx context.Context, cloudLoadbalancerId, cloudListenerId, cloudCertId string) error {
	// 查询监听器列表
	// REF: https://cloud.tencent.com/document/api/214/30686
	describeListenersReq := tcClb.NewDescribeListenersRequest()
	describeListenersReq.LoadBalancerId = common.StringPtr(cloudLoadbalancerId)
	describeListenersReq.ListenerIds = common.StringPtrs([]string{cloudListenerId})
	describeListenersResp, err := d.sdkClients.clb.DescribeListeners(describeListenersReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'clb.DescribeListeners'")
	}
	if len(describeListenersResp.Response.Listeners) == 0 {
		return errors.New("listener not found")
	}

	d.logger.Logt("已查询到监听器属性", describeListenersResp.Response)

	// 修改监听器属性
	// REF: https://cloud.tencent.com/document/product/214/30681
	modifyListenerReq := tcClb.NewModifyListenerRequest()
	modifyListenerReq.LoadBalancerId = common.StringPtr(cloudLoadbalancerId)
	modifyListenerReq.ListenerId = common.StringPtr(cloudListenerId)
	modifyListenerReq.Certificate = &tcClb.CertificateInput{CertId: common.StringPtr(cloudCertId)}
	if describeListenersResp.Response.Listeners[0].Certificate != nil && describeListenersResp.Response.Listeners[0].Certificate.SSLMode != nil {
		modifyListenerReq.Certificate.SSLMode = describeListenersResp.Response.Listeners[0].Certificate.SSLMode
		modifyListenerReq.Certificate.CertCaId = describeListenersResp.Response.Listeners[0].Certificate.CertCaId
	} else {
		modifyListenerReq.Certificate.SSLMode = common.StringPtr("UNIDIRECTIONAL")
	}
	modifyListenerResp, err := d.sdkClients.clb.ModifyListener(modifyListenerReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'clb.ModifyListener'")
	}

	d.logger.Logt("已修改监听器属性", modifyListenerResp.Response)

	return nil
}

func createSdkClients(secretId, secretKey, region string) (*wSdkClients, error) {
	credential := common.NewCredential(secretId, secretKey)

	// 注意虽然官方文档中地域无需指定，但实际需要部署到 CLB 时必传
	sslClient, err := tcSsl.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	clbClient, err := tcClb.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return &wSdkClients{
		ssl: sslClient,
		clb: clbClient,
	}, nil
}
