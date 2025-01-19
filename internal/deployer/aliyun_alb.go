package deployer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	aliyunAlb "github.com/alibabacloud-go/alb-20200616/v2/client"
	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"certimate/internal/domain"
	"certimate/internal/pkg/core/uploader"
	uploaderAliyunCas "certimate/internal/pkg/core/uploader/providers/aliyun-cas"
)

type AliyunALBDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClient   *aliyunAlb.Client
	sslUploader uploader.Uploader
}

func NewAliyunALBDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.AliyunAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}

	client, err := (&AliyunALBDeployer{}).createSdkClient(
		access.AccessKeyId,
		access.AccessKeySecret,
		option.DeployConfig.GetConfigAsString("region"),
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	aliCasRegion := option.DeployConfig.GetConfigAsString("region")
	if aliCasRegion != "" {
		// 阿里云 CAS 服务接入点是独立于 ALB 服务的
		// 国内版接入点：华东一杭州
		// 国际版接入点：亚太东南一新加坡
		if !strings.HasPrefix(aliCasRegion, "cn-") {
			aliCasRegion = "ap-southeast-1"
		} else {
			aliCasRegion = "cn-hangzhou"
		}
	}
	uploader, err := uploaderAliyunCas.New(&uploaderAliyunCas.AliyunCASUploaderConfig{
		AccessKeyId:     access.AccessKeyId,
		AccessKeySecret: access.AccessKeySecret,
		Region:          aliCasRegion,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &AliyunALBDeployer{
		option:      option,
		infos:       make([]string, 0),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *AliyunALBDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *AliyunALBDeployer) GetInfos() []string {
	return d.infos
}

func (d *AliyunALBDeployer) Deploy(ctx context.Context) error {
	switch d.option.DeployConfig.GetConfigAsString("resourceType") {
	case "loadbalancer":
		if err := d.deployToLoadbalancer(ctx); err != nil {
			return err
		}
	case "listener":
		if err := d.deployToListener(ctx); err != nil {
			return err
		}
	default:
		return errors.New("unsupported resource type")
	}

	return nil
}

func (d *AliyunALBDeployer) createSdkClient(accessKeyId, accessKeySecret, region string) (*aliyunAlb.Client, error) {
	if region == "" {
		region = "cn-hangzhou" // ALB 服务默认区域：华东一杭州
	}

	aConfig := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}

	var endpoint string
	switch region {
	case "cn-hangzhou-finance":
		endpoint = "alb.cn-hangzhou.aliyuncs.com"
	default:
		endpoint = fmt.Sprintf("alb.%s.aliyuncs.com", region)
	}
	aConfig.Endpoint = tea.String(endpoint)

	client, err := aliyunAlb.NewClient(aConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (d *AliyunALBDeployer) deployToLoadbalancer(ctx context.Context) error {
	aliLoadbalancerId := d.option.DeployConfig.GetConfigAsString("loadbalancerId")
	if aliLoadbalancerId == "" {
		return errors.New("`loadbalancerId` is required")
	}

	aliListenerIds := make([]string, 0)

	// 查询负载均衡实例的详细信息
	// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-getloadbalancerattribute
	getLoadBalancerAttributeReq := &aliyunAlb.GetLoadBalancerAttributeRequest{
		LoadBalancerId: tea.String(aliLoadbalancerId),
	}
	getLoadBalancerAttributeResp, err := d.sdkClient.GetLoadBalancerAttribute(getLoadBalancerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'alb.GetLoadBalancerAttribute'")
	}

	d.infos = append(d.infos, toStr("已查询到 ALB 负载均衡实例", getLoadBalancerAttributeResp))

	// 查询 HTTPS 监听列表
	// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-listlisteners
	listListenersPage := 1
	listListenersLimit := int32(100)
	var listListenersToken *string = nil
	for {
		listListenersReq := &aliyunAlb.ListListenersRequest{
			MaxResults:       tea.Int32(listListenersLimit),
			NextToken:        listListenersToken,
			LoadBalancerIds:  []*string{tea.String(aliLoadbalancerId)},
			ListenerProtocol: tea.String("HTTPS"),
		}
		listListenersResp, err := d.sdkClient.ListListeners(listListenersReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'alb.ListListeners'")
		}

		if listListenersResp.Body.Listeners != nil {
			for _, listener := range listListenersResp.Body.Listeners {
				aliListenerIds = append(aliListenerIds, *listener.ListenerId)
			}
		}

		if listListenersResp.Body.NextToken == nil {
			break
		} else {
			listListenersToken = listListenersResp.Body.NextToken
			listListenersPage += 1
		}
	}

	d.infos = append(d.infos, toStr("已查询到 ALB 负载均衡实例下的全部 HTTPS 监听", aliListenerIds))

	// 查询 QUIC 监听列表
	// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-listlisteners
	listListenersPage = 1
	listListenersToken = nil
	for {
		listListenersReq := &aliyunAlb.ListListenersRequest{
			MaxResults:       tea.Int32(listListenersLimit),
			NextToken:        listListenersToken,
			LoadBalancerIds:  []*string{tea.String(aliLoadbalancerId)},
			ListenerProtocol: tea.String("QUIC"),
		}
		listListenersResp, err := d.sdkClient.ListListeners(listListenersReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'alb.ListListeners'")
		}

		if listListenersResp.Body.Listeners != nil {
			for _, listener := range listListenersResp.Body.Listeners {
				aliListenerIds = append(aliListenerIds, *listener.ListenerId)
			}
		}

		if listListenersResp.Body.NextToken == nil {
			break
		} else {
			listListenersToken = listListenersResp.Body.NextToken
			listListenersPage += 1
		}
	}

	d.infos = append(d.infos, toStr("已查询到 ALB 负载均衡实例下的全部 QUIC 监听", aliListenerIds))

	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	// 批量更新监听证书
	var errs []error
	for _, aliListenerId := range aliListenerIds {
		if err := d.updateListenerCertificate(ctx, aliListenerId, upres.CertId); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (d *AliyunALBDeployer) deployToListener(ctx context.Context) error {
	aliListenerId := d.option.DeployConfig.GetConfigAsString("listenerId")
	if aliListenerId == "" {
		return errors.New("`listenerId` is required")
	}

	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	// 更新监听
	if err := d.updateListenerCertificate(ctx, aliListenerId, upres.CertId); err != nil {
		return err
	}

	return nil
}

func (d *AliyunALBDeployer) updateListenerCertificate(ctx context.Context, aliListenerId string, aliCertId string) error {
	// 查询监听的属性
	// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-getlistenerattribute
	getListenerAttributeReq := &aliyunAlb.GetListenerAttributeRequest{
		ListenerId: tea.String(aliListenerId),
	}
	getListenerAttributeResp, err := d.sdkClient.GetListenerAttribute(getListenerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'alb.GetListenerAttribute'")
	}

	d.infos = append(d.infos, toStr("已查询到 ALB 监听配置", getListenerAttributeResp))

	// 修改监听的属性
	// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-updatelistenerattribute
	updateListenerAttributeReq := &aliyunAlb.UpdateListenerAttributeRequest{
		ListenerId: tea.String(aliListenerId),
		Certificates: []*aliyunAlb.UpdateListenerAttributeRequestCertificates{{
			CertificateId: tea.String(aliCertId),
		}},
	}
	updateListenerAttributeResp, err := d.sdkClient.UpdateListenerAttribute(updateListenerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'alb.UpdateListenerAttribute'")
	}

	d.infos = append(d.infos, toStr("已更新 ALB 监听配置", updateListenerAttributeResp))

	return nil
}
