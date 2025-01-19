import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { ClientResponseError } from "pocketbase";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { PbErrorData } from "@/domain/base";
import { Access, accessProvidersMap, accessTypeFormSchema, UnicloudConfig } from "@/domain/access";
import { save } from "@/repository/access";
import { useConfigContext } from "@/providers/config";

type AccessUnicloudFormProps = {
  op: "add" | "edit" | "copy";
  data?: Access;
  onAfterReq: () => void;
};

const AccessUnicloudForm = ({ data, op, onAfterReq }: AccessUnicloudFormProps) => {
  const { addAccess, updateAccess } = useConfigContext();
  const { t } = useTranslation();
  const formSchema = z.object({
    id: z.string().optional(),
    username: z
      .string()
      .min(1, "access.authorization.form.username.placeholder")
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    token: z.string().min(1, "access.authorization.form.token.placeholder"),
  });

  let config: UnicloudConfig = {
    username: "",
    token: "",
  };
  if (data) config = data.config as UnicloudConfig;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      id: data?.id,
      username: config.username,
      token: config.token,
    },
  });

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    const req: Access = {
      id: data.id as string,
      name: data.username,
      configType: "unicloud",
      usage: accessProvidersMap.get("unicloud")!.usage,
      config: {
        username: data.username,
        token: data.token,
      },
    };

    try {
      req.id = op == "copy" ? "" : req.id;
      const rs = await save(req);

      onAfterReq();

      req.id = rs.id;
      req.created = rs.created;
      req.updated = rs.updated;
      if (data.id && op == "edit") {
        updateAccess(req);
        return;
      }
      addAccess(req);
    } catch (e) {
      const err = e as ClientResponseError;

      Object.entries(err.response.data as PbErrorData).forEach(([key, value]) => {
        form.setError(key as keyof z.infer<typeof formSchema>, {
          type: "manual",
          message: value.message,
        });
      });
    }
  };

  return (
    <>
      <Form {...form}>
        <form
          onSubmit={(e) => {
            e.stopPropagation();
            form.handleSubmit(onSubmit)(e);
          }}
          className="space-y-8"
        >
          <FormField
            control={form.control}
            name="username"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("access.authorization.form.username.label")}</FormLabel>
                <FormControl>
                  <Input placeholder={t("access.authorization.form.username.placeholder")} {...field} />
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="token"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("access.authorization.form.token.label")}</FormLabel>
                <FormControl>
                  <Input type="text" placeholder={t("access.authorization.form.password.placeholder")} {...field} />
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          <div className="flex justify-end">
            <Button type="submit">{t("common.save")}</Button>
          </div>
        </form>
      </Form>
    </>
  );
};

export default AccessUnicloudForm;
