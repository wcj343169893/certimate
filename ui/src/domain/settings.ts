export const SETTINGS_NAME_EMAILS = "emails" as const;
export const SETTINGS_NAME_NOTIFYTEMPLATES = "notifyTemplates" as const;
export const SETTINGS_NAME_NOTIFYCHANNELS = "notifyChannels" as const;
export const SETTINGS_NAME_SSLPROVIDER = "sslProvider" as const;
export const SETTINGS_NAMES = Object.freeze({
  EMAILS: SETTINGS_NAME_EMAILS,
  NOTIFY_TEMPLATES: SETTINGS_NAME_NOTIFYTEMPLATES,
  NOTIFY_CHANNELS: SETTINGS_NAME_NOTIFYCHANNELS,
  SSL_PROVIDER: SETTINGS_NAME_SSLPROVIDER,
} as const);

export type SettingsNames = (typeof SETTINGS_NAMES)[keyof typeof SETTINGS_NAMES];

export interface SettingsModel<T extends NonNullable<unknown> = NonNullable<unknown>> extends BaseModel {
  name: string;
  content: T;
}

// #region Settings: Emails
export type EmailsSettingsContent = {
  emails: string[];
};
// #endregion

// #region Settings: NotifyTemplates
export type NotifyTemplatesSettingsContent = {
  notifyTemplates: NotifyTemplate[];
};

export type NotifyTemplate = {
  subject: string;
  message: string;
};

export const defaultNotifyTemplate: NotifyTemplate = {
  subject: "您有 ${COUNT} 张证书即将过期",
  message: "有 ${COUNT} 张证书即将过期，域名分别为 ${DOMAINS}，请保持关注！",
};
// #endregion

// #region Settings: NotifyChannels
export const NOTIFY_CHANNELS = Object.freeze({
  BARK: "bark",
  DINGTALK: "dingtalk",
  EMAIL: "email",
  LARK: "lark",
  SERVERCHAN: "serverchan",
  TELEGRAM: "telegram",
  WEBHOOK: "webhook",
  WECOM: "wecom",
} as const);

export type NotifyChannels = (typeof NOTIFY_CHANNELS)[keyof typeof NOTIFY_CHANNELS];

export type NotifyChannelsSettingsContent = {
  /*
    注意：如果追加新的类型，请保持以 ASCII 排序。
    NOTICE: If you add new type, please keep ASCII order.
  */
  [key: string]: ({ enabled?: boolean } & Record<string, unknown>) | undefined;
  [NOTIFY_CHANNELS.BARK]?: BarkNotifyChannelConfig;
  [NOTIFY_CHANNELS.DINGTALK]?: DingTalkNotifyChannelConfig;
  [NOTIFY_CHANNELS.EMAIL]?: EmailNotifyChannelConfig;
  [NOTIFY_CHANNELS.LARK]?: LarkNotifyChannelConfig;
  [NOTIFY_CHANNELS.SERVERCHAN]?: ServerChanNotifyChannelConfig;
  [NOTIFY_CHANNELS.TELEGRAM]?: TelegramNotifyChannelConfig;
  [NOTIFY_CHANNELS.WEBHOOK]?: WebhookNotifyChannelConfig;
  [NOTIFY_CHANNELS.WECOM]?: WeComNotifyChannelConfig;
};

export type BarkNotifyChannelConfig = {
  deviceKey: string;
  serverUrl: string;
  enabled?: boolean;
};

export type EmailNotifyChannelConfig = {
  smtpHost: string;
  smtpPort: number;
  smtpTLS: boolean;
  username: string;
  password: string;
  senderAddress: string;
  receiverAddress: string;
  enabled?: boolean;
};

export type DingTalkNotifyChannelConfig = {
  accessToken: string;
  secret: string;
  enabled?: boolean;
};

export type LarkNotifyChannelConfig = {
  webhookUrl: string;
  enabled?: boolean;
};

export type ServerChanNotifyChannelConfig = {
  url: string;
  enabled?: boolean;
};

export type TelegramNotifyChannelConfig = {
  apiToken: string;
  chatId: string;
  enabled?: boolean;
};

export type WebhookNotifyChannelConfig = {
  url: string;
  enabled?: boolean;
};

export type WeComNotifyChannelConfig = {
  webhookUrl: string;
  enabled?: boolean;
};

export type NotifyChannel = {
  type: string;
  name: string;
};

export const notifyChannelsMap: Map<NotifyChannel["type"], NotifyChannel> = new Map(
  [
    [NOTIFY_CHANNELS.EMAIL, "common.notifier.email"],
    [NOTIFY_CHANNELS.DINGTALK, "common.notifier.dingtalk"],
    [NOTIFY_CHANNELS.LARK, "common.notifier.lark"],
    [NOTIFY_CHANNELS.WECOM, "common.notifier.wecom"],
    [NOTIFY_CHANNELS.TELEGRAM, "common.notifier.telegram"],
    [NOTIFY_CHANNELS.SERVERCHAN, "common.notifier.serverchan"],
    [NOTIFY_CHANNELS.BARK, "common.notifier.bark"],
    [NOTIFY_CHANNELS.WEBHOOK, "common.notifier.webhook"],
  ].map(([type, name]) => [type, { type, name }])
);
// #endregion

// #region Settings: SSLProvider
export const SSLPROVIDER_LETSENCRYPT = "letsencrypt" as const;
export const SSLPROVIDER_ZEROSSL = "zerossl" as const;
export const SSLPROVIDER_GOOGLETRUSTSERVICES = "gts" as const;
export const SSLPROVIDERS = Object.freeze({
  LETS_ENCRYPT: SSLPROVIDER_LETSENCRYPT,
  ZERO_SSL: SSLPROVIDER_ZEROSSL,
  GOOGLE_TRUST_SERVICES: SSLPROVIDER_GOOGLETRUSTSERVICES,
} as const);

export type SSLProviders = (typeof SSLPROVIDERS)[keyof typeof SSLPROVIDERS];

export type SSLProviderSettingsContent = {
  provider: (typeof SSLPROVIDERS)[keyof typeof SSLPROVIDERS];
  config: {
    [key: string]: Record<string, unknown> | undefined;
    letsencrypt?: SSLProviderLetsEncryptConfig;
    zerossl?: SSLProviderZeroSSLConfig;
    gts?: SSLProviderGoogleTrustServicesConfig;
  };
};

export type SSLProviderLetsEncryptConfig = NonNullable<unknown>;

export type SSLProviderZeroSSLConfig = {
  eabKid: string;
  eabHmacKey: string;
};

export type SSLProviderGoogleTrustServicesConfig = {
  eabKid: string;
  eabHmacKey: string;
};
// #endregion
