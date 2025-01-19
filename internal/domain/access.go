package domain

type AliyunAccess struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

type ByteplusAccess struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type TencentAccess struct {
	SecretId  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
}

type HuaweiCloudAccess struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	Region          string `json:"region"`
}

type BaiduCloudAccess struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type AwsAccess struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	Region          string `json:"region"`
	HostedZoneId    string `json:"hostedZoneId"`
}

type CloudflareAccess struct {
	DnsApiToken string `json:"dnsApiToken"`
}

type QiniuAccess struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type DogeCloudAccess struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type NameSiloAccess struct {
	ApiKey string `json:"apiKey"`
}

type GodaddyAccess struct {
	ApiKey    string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
}

type PdnsAccess struct {
	ApiUrl string `json:"apiUrl"`
	ApiKey string `json:"apiKey"`
}

type VolcEngineAccess struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`

	// Deprecated: Use [AccessKey] and [SecretKey] instead in the future
	AccessKeyId string `json:"accessKeyId"`
	// Deprecated: Use [AccessKey] and [SecretKey] instead in the future
	SecretAccessKey string `json:"secretAccessKey"`
}

type HttpreqAccess struct {
	Endpoint string `json:"endpoint"`
	Mode     string `json:"mode"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LocalAccess struct{}

type SSHAccess struct {
	Host          string `json:"host"`
	Port          string `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Key           string `json:"key"`
	KeyPassphrase string `json:"keyPassphrase"`
}

type WebhookAccess struct {
	Url string `json:"url"`
}

type UnicloudAccess struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
}

type KubernetesAccess struct {
	KubeConfig string `json:"kubeConfig"`
}
