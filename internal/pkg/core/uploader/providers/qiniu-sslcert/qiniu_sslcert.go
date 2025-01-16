package qiniusslcert

import (
	"context"
	"errors"
	"fmt"
	"time"

	xerrors "github.com/pkg/errors"
	"github.com/qiniu/go-sdk/v7/auth"

	"certimate/internal/pkg/core/uploader"
	"certimate/internal/pkg/utils/x509"
	qiniuEx "certimate/internal/pkg/vendors/qiniu-sdk"
)

type QiniuSSLCertUploaderConfig struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type QiniuSSLCertUploader struct {
	config    *QiniuSSLCertUploaderConfig
	sdkClient *qiniuEx.Client
}

var _ uploader.Uploader = (*QiniuSSLCertUploader)(nil)

func New(config *QiniuSSLCertUploaderConfig) (*QiniuSSLCertUploader, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	client, err := createSdkClient(
		config.AccessKey,
		config.SecretKey,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &QiniuSSLCertUploader{
		config:    config,
		sdkClient: client,
	}, nil
}

func (u *QiniuSSLCertUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 解析证书内容
	certX509, err := x509.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, err
	}

	// 生成新证书名（需符合七牛云命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://developer.qiniu.com/fusion/8593/interface-related-certificate
	uploadSslCertResp, err := u.sdkClient.UploadSslCert(certName, certX509.Subject.CommonName, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.UploadSslCert'")
	}

	certId = uploadSslCertResp.CertID
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func createSdkClient(accessKey, secretKey string) (*qiniuEx.Client, error) {
	credential := auth.New(accessKey, secretKey)
	client := qiniuEx.NewClient(credential)
	return client, nil
}
