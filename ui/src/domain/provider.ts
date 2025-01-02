/*
  注意：如果追加新的常量值，请保持以 ASCII 排序。
  NOTICE: If you add new constant, please keep ASCII order.
 */
export const ACCESS_PROVIDERS = Object.freeze({
  ACMEHTTPREQ: "acmehttpreq",
  ALIYUN: "aliyun",
  AWS: "aws",
  BAIDUCLOUD: "baiducloud",
  BYTEPLUS: "byteplus",
  CLOUDFLARE: "cloudflare",
  DOGECLOUD: "dogecloud",
  GODADDY: "godaddy",
  HUAWEICLOUD: "huaweicloud",
  KUBERNETES: "k8s",
  LOCAL: "local",
  NAMEDOTCOM: "namedotcom",
  NAMESILO: "namesilo",
  POWERDNS: "powerdns",
  QINIU: "qiniu",
  SSH: "ssh",
  TENCENTCLOUD: "tencentcloud",
  VOLCENGINE: "volcengine",
  WEBHOOK: "webhook",
} as const);

export type AccessProviderType = (typeof ACCESS_PROVIDERS)[keyof typeof ACCESS_PROVIDERS];

export const ACCESS_USAGES = Object.freeze({
  ALL: "all",
  APPLY: "apply",
  DEPLOY: "deploy",
} as const);

export type AccessUsageType = (typeof ACCESS_USAGES)[keyof typeof ACCESS_USAGES];

export type AccessProvider = {
  type: AccessProviderType;
  name: string;
  icon: string;
  usage: AccessUsageType;
};

export const accessProvidersMap: Map<AccessProvider["type"] | string, AccessProvider> = new Map(
  /*
   注意：此处的顺序决定显示在前端的顺序。
   NOTICE: The following order determines the order displayed at the frontend.
  */
  [
    [ACCESS_PROVIDERS.LOCAL, "common.provider.local", "/imgs/providers/local.svg", "deploy"],
    [ACCESS_PROVIDERS.SSH, "common.provider.ssh", "/imgs/providers/ssh.svg", "deploy"],
    [ACCESS_PROVIDERS.WEBHOOK, "common.provider.webhook", "/imgs/providers/webhook.svg", "deploy"],
    [ACCESS_PROVIDERS.KUBERNETES, "common.provider.kubernetes", "/imgs/providers/kubernetes.svg", "deploy"],
    [ACCESS_PROVIDERS.ALIYUN, "common.provider.aliyun", "/imgs/providers/aliyun.svg", "all"],
    [ACCESS_PROVIDERS.TENCENTCLOUD, "common.provider.tencentcloud", "/imgs/providers/tencentcloud.svg", "all"],
    [ACCESS_PROVIDERS.HUAWEICLOUD, "common.provider.huaweicloud", "/imgs/providers/huaweicloud.svg", "all"],
    [ACCESS_PROVIDERS.BAIDUCLOUD, "common.provider.baiducloud", "/imgs/providers/baiducloud.svg", "all"],
    [ACCESS_PROVIDERS.QINIU, "common.provider.qiniu", "/imgs/providers/qiniu.svg", "deploy"],
    [ACCESS_PROVIDERS.DOGECLOUD, "common.provider.dogecloud", "/imgs/providers/dogecloud.svg", "deploy"],
    [ACCESS_PROVIDERS.VOLCENGINE, "common.provider.volcengine", "/imgs/providers/volcengine.svg", "all"],
    [ACCESS_PROVIDERS.BYTEPLUS, "common.provider.byteplus", "/imgs/providers/byteplus.svg", "all"],
    [ACCESS_PROVIDERS.AWS, "common.provider.aws", "/imgs/providers/aws.svg", "apply"],
    [ACCESS_PROVIDERS.CLOUDFLARE, "common.provider.cloudflare", "/imgs/providers/cloudflare.svg", "apply"],
    [ACCESS_PROVIDERS.NAMEDOTCOM, "common.provider.namedotcom", "/imgs/providers/namedotcom.svg", "apply"],
    [ACCESS_PROVIDERS.NAMESILO, "common.provider.namesilo", "/imgs/providers/namesilo.svg", "apply"],
    [ACCESS_PROVIDERS.GODADDY, "common.provider.godaddy", "/imgs/providers/godaddy.svg", "apply"],
    [ACCESS_PROVIDERS.POWERDNS, "common.provider.powerdns", "/imgs/providers/powerdns.svg", "apply"],
    [ACCESS_PROVIDERS.ACMEHTTPREQ, "common.provider.acmehttpreq", "/imgs/providers/acmehttpreq.svg", "apply"],
  ].map(([type, name, icon, usage]) => [
    type,
    {
      type: type as AccessProviderType,
      name: name,
      icon: icon,
      usage: usage as AccessUsageType,
    },
  ])
);

/*
  注意：如果追加新的常量值，请保持以 ASCII 排序。
  NOTICE: If you add new constant, please keep ASCII order.
 */
export const DEPLOY_PROVIDERS = Object.freeze({
  ALIYUN_ALB: `${ACCESS_PROVIDERS.ALIYUN}-alb`,
  ALIYUN_CDN: `${ACCESS_PROVIDERS.ALIYUN}-cdn`,
  ALIYUN_CLB: `${ACCESS_PROVIDERS.ALIYUN}-clb`,
  ALIYUN_DCDN: `${ACCESS_PROVIDERS.ALIYUN}-dcdn`,
  ALIYUN_NLB: `${ACCESS_PROVIDERS.ALIYUN}-nlb`,
  ALIYUN_OSS: `${ACCESS_PROVIDERS.ALIYUN}-oss`,
  BAIDUCLOUD_CDN: `${ACCESS_PROVIDERS.BAIDUCLOUD}-cdn`,
  BYTEPLUS_CDN: `${ACCESS_PROVIDERS.BYTEPLUS}-cdn`,
  DOGECLOUD_CDN: `${ACCESS_PROVIDERS.DOGECLOUD}-cdn`,
  HUAWEICLOUD_CDN: `${ACCESS_PROVIDERS.HUAWEICLOUD}-cdn`,
  HUAWEICLOUD_ELB: `${ACCESS_PROVIDERS.HUAWEICLOUD}-elb`,
  KUBERNETES_SECRET: `${ACCESS_PROVIDERS.KUBERNETES}-secret`,
  LOCAL: `${ACCESS_PROVIDERS.LOCAL}`,
  QINIU_CDN: `${ACCESS_PROVIDERS.QINIU}-cdn`,
  SSH: `${ACCESS_PROVIDERS.SSH}`,
  TENCENTCLOUD_CDN: `${ACCESS_PROVIDERS.TENCENTCLOUD}-cdn`,
  TENCENTCLOUD_CLB: `${ACCESS_PROVIDERS.TENCENTCLOUD}-clb`,
  TENCENTCLOUD_COS: `${ACCESS_PROVIDERS.TENCENTCLOUD}-cos`,
  TENCENTCLOUD_ECDN: `${ACCESS_PROVIDERS.TENCENTCLOUD}-ecdn`,
  TENCENTCLOUD_EO: `${ACCESS_PROVIDERS.TENCENTCLOUD}-eo`,
  VOLCENGINE_CDN: `${ACCESS_PROVIDERS.VOLCENGINE}-cdn`,
  VOLCENGINE_LIVE: `${ACCESS_PROVIDERS.VOLCENGINE}-live`,
  WEBHOOK: `${ACCESS_PROVIDERS.WEBHOOK}`,
} as const);

export type DeployProviderType = (typeof DEPLOY_PROVIDERS)[keyof typeof DEPLOY_PROVIDERS];

export type DeployProvider = {
  type: DeployProviderType;
  name: string;
  icon: string;
  provider: AccessProviderType;
};

export const deployProvidersMap: Map<DeployProvider["type"] | string, DeployProvider> = new Map(
  /*
   注意：此处的顺序决定显示在前端的顺序。
   NOTICE: The following order determines the order displayed at the frontend.
  */
  [
    [DEPLOY_PROVIDERS.LOCAL, "common.provider.local"],
    [DEPLOY_PROVIDERS.SSH, "common.provider.ssh"],
    [DEPLOY_PROVIDERS.WEBHOOK, "common.provider.webhook"],
    [DEPLOY_PROVIDERS.KUBERNETES_SECRET, "common.provider.kubernetes.secret"],
    [DEPLOY_PROVIDERS.ALIYUN_OSS, "common.provider.aliyun.oss"],
    [DEPLOY_PROVIDERS.ALIYUN_CDN, "common.provider.aliyun.cdn"],
    [DEPLOY_PROVIDERS.ALIYUN_DCDN, "common.provider.aliyun.dcdn"],
    [DEPLOY_PROVIDERS.ALIYUN_CLB, "common.provider.aliyun.clb"],
    [DEPLOY_PROVIDERS.ALIYUN_ALB, "common.provider.aliyun.alb"],
    [DEPLOY_PROVIDERS.ALIYUN_NLB, "common.provider.aliyun.nlb"],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_CDN, "common.provider.tencentcloud.cdn"],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_ECDN, "common.provider.tencentcloud.ecdn"],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_CLB, "common.provider.tencentcloud.clb"],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_COS, "common.provider.tencentcloud.cos"],
    [DEPLOY_PROVIDERS.TENCENTCLOUD_EO, "common.provider.tencentcloud.eo"],
    [DEPLOY_PROVIDERS.HUAWEICLOUD_CDN, "common.provider.huaweicloud.cdn"],
    [DEPLOY_PROVIDERS.HUAWEICLOUD_ELB, "common.provider.huaweicloud.elb"],
    [DEPLOY_PROVIDERS.BAIDUCLOUD_CDN, "common.provider.baiducloud.cdn"],
    [DEPLOY_PROVIDERS.VOLCENGINE_CDN, "common.provider.volcengine.cdn"],
    [DEPLOY_PROVIDERS.VOLCENGINE_LIVE, "common.provider.volcengine.live"],
    [DEPLOY_PROVIDERS.QINIU_CDN, "common.provider.qiniu.cdn"],
    [DEPLOY_PROVIDERS.DOGECLOUD_CDN, "common.provider.dogecloud.cdn"],
    [DEPLOY_PROVIDERS.BYTEPLUS_CDN, "common.provider.byteplus.cdn"],
  ].map(([type, name]) => [
    type,
    {
      type: type as DeployProviderType,
      name: name,
      icon: accessProvidersMap.get(type.split("-")[0])!.icon,
      provider: type.split("-")[0] as AccessProviderType,
    },
  ])
);
