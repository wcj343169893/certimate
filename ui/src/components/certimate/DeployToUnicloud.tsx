import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { useDeployEditContext } from "./DeployEdit";

type DeployToUnicloudConfigParams = {
  spaceId: string;
};

const DeployToUnicloud = () => {
  const { t } = useTranslation();
  const { config, setConfig, errors, setErrors } = useDeployEditContext<DeployToUnicloudConfigParams>();

  useEffect(() => {
    if (!config.id) {
      setConfig({
        ...config,
        config: {
          spaceId: "",
        },
      });
    }
  }, []);

  useEffect(() => {
    setErrors({});
  }, []);

  return (
    <>
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
    </>
  );
};

export default DeployToUnicloud;
