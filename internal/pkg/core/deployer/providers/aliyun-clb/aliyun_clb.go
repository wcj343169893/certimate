package aliyunclb

import (
	"context"
	"errors"
	"fmt"

	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	aliyunSlb "github.com/alibabacloud-go/slb-20140515/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"certimate/internal/pkg/core/deployer"
	"certimate/internal/pkg/core/uploader"
	providerSlb "certimate/internal/pkg/core/uploader/providers/aliyun-slb"
)

type AliyunCLBDeployerConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 阿里云地域。
	Region string `json:"region"`
	// 部署资源类型。
	ResourceType DeployResourceType `json:"resourceType"`
	// 负载均衡实例 ID。
	// 部署资源类型为 [DEPLOY_RESOURCE_LOADBALANCER]、[DEPLOY_RESOURCE_LISTENER] 时必填。
	LoadbalancerId string `json:"loadbalancerId,omitempty"`
	// 负载均衡监听端口。
	// 部署资源类型为 [DEPLOY_RESOURCE_LISTENER] 时必填。
	ListenerPort int32 `json:"listenerPort,omitempty"`
}

type AliyunCLBDeployer struct {
	config      *AliyunCLBDeployerConfig
	logger      deployer.Logger
	sdkClient   *aliyunSlb.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*AliyunCLBDeployer)(nil)

func New(config *AliyunCLBDeployerConfig) (*AliyunCLBDeployer, error) {
	return NewWithLogger(config, deployer.NewNilLogger())
}

func NewWithLogger(config *AliyunCLBDeployerConfig, logger deployer.Logger) (*AliyunCLBDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := providerSlb.New(&providerSlb.AliyunSLBUploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
		Region:          config.Region,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &AliyunCLBDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *AliyunCLBDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 SLB
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	// 根据部署资源类型决定部署方式
	switch d.config.ResourceType {
	case DEPLOY_RESOURCE_LOADBALANCER:
		if err := d.deployToLoadbalancer(ctx, upres.CertId); err != nil {
			return nil, err
		}

	case DEPLOY_RESOURCE_LISTENER:
		if err := d.deployToListener(ctx, upres.CertId); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported resource type: %s", d.config.ResourceType)
	}

	return &deployer.DeployResult{}, nil
}

func (d *AliyunCLBDeployer) deployToLoadbalancer(ctx context.Context, cloudCertId string) error {
	if d.config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}

	listenerPorts := make([]int32, 0)

	// 查询负载均衡实例的详细信息
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-describeloadbalancerattribute
	describeLoadBalancerAttributeReq := &aliyunSlb.DescribeLoadBalancerAttributeRequest{
		RegionId:       tea.String(d.config.Region),
		LoadBalancerId: tea.String(d.config.LoadbalancerId),
	}
	describeLoadBalancerAttributeResp, err := d.sdkClient.DescribeLoadBalancerAttribute(describeLoadBalancerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'slb.DescribeLoadBalancerAttribute'")
	}

	d.logger.Logt("已查询到 CLB 负载均衡实例", describeLoadBalancerAttributeResp)

	// 查询 HTTPS 监听列表
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-describeloadbalancerlisteners
	listListenersPage := 1
	listListenersLimit := int32(100)
	var listListenersToken *string = nil
	for {
		describeLoadBalancerListenersReq := &aliyunSlb.DescribeLoadBalancerListenersRequest{
			RegionId:         tea.String(d.config.Region),
			MaxResults:       tea.Int32(listListenersLimit),
			NextToken:        listListenersToken,
			LoadBalancerId:   []*string{tea.String(d.config.LoadbalancerId)},
			ListenerProtocol: tea.String("https"),
		}
		describeLoadBalancerListenersResp, err := d.sdkClient.DescribeLoadBalancerListeners(describeLoadBalancerListenersReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'slb.DescribeLoadBalancerListeners'")
		}

		if describeLoadBalancerListenersResp.Body.Listeners != nil {
			for _, listener := range describeLoadBalancerListenersResp.Body.Listeners {
				listenerPorts = append(listenerPorts, *listener.ListenerPort)
			}
		}

		if len(describeLoadBalancerListenersResp.Body.Listeners) == 0 || describeLoadBalancerListenersResp.Body.NextToken == nil {
			break
		} else {
			listListenersToken = describeLoadBalancerListenersResp.Body.NextToken
			listListenersPage += 1
		}
	}

	d.logger.Logt("已查询到 CLB 负载均衡实例下的全部 HTTPS 监听", listenerPorts)

	// 批量更新监听证书
	var errs []error
	for _, listenerPort := range listenerPorts {
		if err := d.updateListenerCertificate(ctx, d.config.LoadbalancerId, listenerPort, cloudCertId); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (d *AliyunCLBDeployer) deployToListener(ctx context.Context, cloudCertId string) error {
	if d.config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}
	if d.config.ListenerPort == 0 {
		return errors.New("config `listenerPort` is required")
	}

	// 更新监听
	if err := d.updateListenerCertificate(ctx, d.config.LoadbalancerId, d.config.ListenerPort, cloudCertId); err != nil {
		return err
	}

	return nil
}

func (d *AliyunCLBDeployer) updateListenerCertificate(ctx context.Context, cloudLoadbalancerId string, cloudListenerPort int32, cloudCertId string) error {
	// 查询监听配置
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-describeloadbalancerhttpslistenerattribute
	describeLoadBalancerHTTPSListenerAttributeReq := &aliyunSlb.DescribeLoadBalancerHTTPSListenerAttributeRequest{
		LoadBalancerId: tea.String(cloudLoadbalancerId),
		ListenerPort:   tea.Int32(cloudListenerPort),
	}
	describeLoadBalancerHTTPSListenerAttributeResp, err := d.sdkClient.DescribeLoadBalancerHTTPSListenerAttribute(describeLoadBalancerHTTPSListenerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'slb.DescribeLoadBalancerHTTPSListenerAttribute'")
	}

	d.logger.Logt("已查询到 CLB HTTPS 监听配置", describeLoadBalancerHTTPSListenerAttributeResp)

	// 查询扩展域名
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-describedomainextensions
	describeDomainExtensionsReq := &aliyunSlb.DescribeDomainExtensionsRequest{
		RegionId:       tea.String(d.config.Region),
		LoadBalancerId: tea.String(cloudLoadbalancerId),
		ListenerPort:   tea.Int32(cloudListenerPort),
	}
	describeDomainExtensionsResp, err := d.sdkClient.DescribeDomainExtensions(describeDomainExtensionsReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'slb.DescribeDomainExtensions'")
	}

	d.logger.Logt("已查询到 CLB 扩展域名", describeDomainExtensionsResp)

	// 遍历修改扩展域名
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-setdomainextensionattribute
	//
	// 这里仅修改跟被替换证书一致的扩展域名
	if describeDomainExtensionsResp.Body.DomainExtensions != nil && describeDomainExtensionsResp.Body.DomainExtensions.DomainExtension != nil {
		for _, domainExtension := range describeDomainExtensionsResp.Body.DomainExtensions.DomainExtension {
			if *domainExtension.ServerCertificateId != *describeLoadBalancerHTTPSListenerAttributeResp.Body.ServerCertificateId {
				continue
			}

			setDomainExtensionAttributeReq := &aliyunSlb.SetDomainExtensionAttributeRequest{
				RegionId:            tea.String(d.config.Region),
				DomainExtensionId:   tea.String(*domainExtension.DomainExtensionId),
				ServerCertificateId: tea.String(cloudCertId),
			}
			_, err := d.sdkClient.SetDomainExtensionAttribute(setDomainExtensionAttributeReq)
			if err != nil {
				return xerrors.Wrap(err, "failed to execute sdk request 'slb.SetDomainExtensionAttribute'")
			}
		}
	}

	// 修改监听配置
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-setloadbalancerhttpslistenerattribute
	//
	// 注意修改监听配置要放在修改扩展域名之后
	setLoadBalancerHTTPSListenerAttributeReq := &aliyunSlb.SetLoadBalancerHTTPSListenerAttributeRequest{
		RegionId:            tea.String(d.config.Region),
		LoadBalancerId:      tea.String(cloudLoadbalancerId),
		ListenerPort:        tea.Int32(cloudListenerPort),
		ServerCertificateId: tea.String(cloudCertId),
	}
	setLoadBalancerHTTPSListenerAttributeResp, err := d.sdkClient.SetLoadBalancerHTTPSListenerAttribute(setLoadBalancerHTTPSListenerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'slb.SetLoadBalancerHTTPSListenerAttribute'")
	}

	d.logger.Logt("已更新 CLB HTTPS 监听配置", setLoadBalancerHTTPSListenerAttributeResp)

	return nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*aliyunSlb.Client, error) {
	if region == "" {
		region = "cn-hangzhou" // CLB(SLB) 服务默认区域：华东一杭州
	}

	// 接入点一览 https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-endpoint
	var endpoint string
	switch region {
	case
		"cn-hangzhou",
		"cn-hangzhou-finance",
		"cn-shanghai-finance-1",
		"cn-shenzhen-finance-1":
		endpoint = "slb.aliyuncs.com"
	default:
		endpoint = fmt.Sprintf("slb.%s.aliyuncs.com", region)
	}

	config := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}

	client, err := aliyunSlb.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
