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
          <Label>{t("Space ID")}</Label>
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
          <Label>{t("Provider")}</Label>
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
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem value="aliyun">阿里云</SelectItem>
                <SelectItem value="tencent">腾讯云</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>
      </div>
    </>
  );
};

export default DeployToUnicloud;
