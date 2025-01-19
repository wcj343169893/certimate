import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useDeployEditContext } from "./DeployEdit";

type DeployToUnicloudConfigParams = {
  spaceId: string;
  provider: string;
  domain: string;
};

const DeployToUnicloud = () => {
  const { t } = useTranslation();
  const { config, setConfig, setErrors } = useDeployEditContext<DeployToUnicloudConfigParams>();

  useEffect(() => {
    if (!config.id) {
      setConfig({
        ...config,
        config: {
          spaceId: "",
          provider: "",
          domain: "",
        },
      });
    }
  }, []);

  useEffect(() => {
    setErrors({});
  }, []);

  return (
    <>
      <div className="flex flex-col space-y-8">
        <div>
          <Label>{t("domain.deployment.form.unicloud.space_id")}</Label>
          <Input
            value={config?.config?.spaceId}
            onChange={(e) => {
              const nv = produce(config, (draft) => {
                draft.config.spaceId = e.target.value;
              });
              setConfig(nv);
            }}
          />
        </div>

        <div>
          <Label>{t("domain.deployment.form.unicloud.provider")}</Label>
          <Select
            value={config?.config?.provider}
            onValueChange={(value) => {
              const nv = produce(config, (draft) => {
                draft.config.provider = value;
              });
              setConfig(nv);
            }}
          >
            <SelectTrigger>
              <SelectValue placeholder={t("domain.deployment.form.unicloud.provider.placeholder")} />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem value="aliyun">{t("domain.deployment.form.unicloud.provider.aliyun")}</SelectItem>
                <SelectItem value="tencent">{t("domain.deployment.form.unicloud.provider.tencent")}</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>
      </div>
    </>
  );
};

export default DeployToUnicloud;
