package aliyunnlb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	aliyunNlb "github.com/alibabacloud-go/nlb-20220430/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"certimate/internal/pkg/core/deployer"
	"certimate/internal/pkg/core/uploader"
	providerCas "certimate/internal/pkg/core/uploader/providers/aliyun-cas"
)

type AliyunNLBDeployerConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 阿里云地域。
	Region string `json:"region"`
	// 部署资源类型。
	ResourceType DeployResourceType `json:"resourceType"`
	// 负载均衡实例 ID。
	// 部署资源类型为 [DEPLOY_RESOURCE_LOADBALANCER] 时必填。
	LoadbalancerId string `json:"loadbalancerId,omitempty"`
	// 负载均衡监听 ID。
	// 部署资源类型为 [DEPLOY_RESOURCE_LISTENER] 时必填。
	ListenerId string `json:"listenerId,omitempty"`
}

type AliyunNLBDeployer struct {
	config      *AliyunNLBDeployerConfig
	logger      deployer.Logger
	sdkClient   *aliyunNlb.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*AliyunNLBDeployer)(nil)

func New(config *AliyunNLBDeployerConfig) (*AliyunNLBDeployer, error) {
	return NewWithLogger(config, deployer.NewNilLogger())
}

func NewWithLogger(config *AliyunNLBDeployerConfig, logger deployer.Logger) (*AliyunNLBDeployer, error) {
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

	aliyunCasRegion := config.Region
	if aliyunCasRegion != "" {
		// 阿里云 CAS 服务接入点是独立于 NLB 服务的
		// 国内版固定接入点：华东一杭州
		// 国际版固定接入点：亚太东南一新加坡
		if !strings.HasPrefix(aliyunCasRegion, "cn-") {
			aliyunCasRegion = "ap-southeast-1"
		} else {
			aliyunCasRegion = "cn-hangzhou"
		}
	}
	uploader, err := providerCas.New(&providerCas.AliyunCASUploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
		Region:          aliyunCasRegion,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &AliyunNLBDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *AliyunNLBDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 CAS
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

func (d *AliyunNLBDeployer) deployToLoadbalancer(ctx context.Context, cloudCertId string) error {
	if d.config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}

	listenerIds := make([]string, 0)

	// 查询负载均衡实例的详细信息
	// REF: https://help.aliyun.com/zh/slb/network-load-balancer/developer-reference/api-nlb-2022-04-30-getloadbalancerattribute
	getLoadBalancerAttributeReq := &aliyunNlb.GetLoadBalancerAttributeRequest{
		LoadBalancerId: tea.String(d.config.LoadbalancerId),
	}
	getLoadBalancerAttributeResp, err := d.sdkClient.GetLoadBalancerAttribute(getLoadBalancerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'nlb.GetLoadBalancerAttribute'")
	}

	d.logger.Logt("已查询到 NLB 负载均衡实例", getLoadBalancerAttributeResp)

	// 查询 TCPSSL 监听列表
	// REF: https://help.aliyun.com/zh/slb/network-load-balancer/developer-reference/api-nlb-2022-04-30-listlisteners
	listListenersPage := 1
	listListenersLimit := int32(100)
	var listListenersToken *string = nil
	for {
		listListenersReq := &aliyunNlb.ListListenersRequest{
			MaxResults:       tea.Int32(listListenersLimit),
			NextToken:        listListenersToken,
			LoadBalancerIds:  []*string{tea.String(d.config.LoadbalancerId)},
			ListenerProtocol: tea.String("TCPSSL"),
		}
		listListenersResp, err := d.sdkClient.ListListeners(listListenersReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'nlb.ListListeners'")
		}

		if listListenersResp.Body.Listeners != nil {
			for _, listener := range listListenersResp.Body.Listeners {
				listenerIds = append(listenerIds, *listener.ListenerId)
			}
		}

		if len(listListenersResp.Body.Listeners) == 0 || listListenersResp.Body.NextToken == nil {
			break
		} else {
			listListenersToken = listListenersResp.Body.NextToken
			listListenersPage += 1
		}
	}

	d.logger.Logt("已查询到 NLB 负载均衡实例下的全部 TCPSSL 监听", listenerIds)

	// 批量更新监听证书
	var errs []error
	for _, listenerId := range listenerIds {
		if err := d.updateListenerCertificate(ctx, listenerId, cloudCertId); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (d *AliyunNLBDeployer) deployToListener(ctx context.Context, cloudCertId string) error {
	if d.config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}

	// 更新监听
	if err := d.updateListenerCertificate(ctx, d.config.ListenerId, cloudCertId); err != nil {
		return err
	}

	return nil
}

func (d *AliyunNLBDeployer) updateListenerCertificate(ctx context.Context, cloudListenerId string, cloudCertId string) error {
	// 查询监听的属性
	// REF: https://help.aliyun.com/zh/slb/network-load-balancer/developer-reference/api-nlb-2022-04-30-getlistenerattribute
	getListenerAttributeReq := &aliyunNlb.GetListenerAttributeRequest{
		ListenerId: tea.String(cloudListenerId),
	}
	getListenerAttributeResp, err := d.sdkClient.GetListenerAttribute(getListenerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'nlb.GetListenerAttribute'")
	}

	d.logger.Logt("已查询到 NLB 监听配置", getListenerAttributeResp)

	// 修改监听的属性
	// REF: https://help.aliyun.com/zh/slb/network-load-balancer/developer-reference/api-nlb-2022-04-30-updatelistenerattribute
	updateListenerAttributeReq := &aliyunNlb.UpdateListenerAttributeRequest{
		ListenerId:     tea.String(cloudListenerId),
		CertificateIds: []*string{tea.String(cloudCertId)},
	}
	updateListenerAttributeResp, err := d.sdkClient.UpdateListenerAttribute(updateListenerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'nlb.UpdateListenerAttribute'")
	}

	d.logger.Logt("已更新 NLB 监听配置", updateListenerAttributeResp)

	return nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*aliyunNlb.Client, error) {
	if region == "" {
		region = "cn-hangzhou" // NLB 服务默认区域：华东一杭州
	}

	// 接入点一览 https://help.aliyun.com/zh/slb/network-load-balancer/developer-reference/api-nlb-2022-04-30-endpoint
	var endpoint string
	switch region {
	default:
		endpoint = fmt.Sprintf("nlb.%s.aliyuncs.com", region)
	}

	config := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}

	client, err := aliyunNlb.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
