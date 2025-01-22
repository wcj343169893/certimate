import { memo, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { Flex, Typography } from "antd";
import { produce } from "immer";

import type { WorkflowNodeConfigForUpload } from "@/domain/workflow";
import { WorkflowNodeType } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import UploadNodeConfigForm, { type UploadNodeConfigFormInstance } from "./UploadNodeConfigForm";
import SharedNode, { type SharedNodeProps } from "./_SharedNode";

export type UploadNodeProps = SharedNodeProps;

const UploadNode = ({ node, disabled }: UploadNodeProps) => {
  if (node.type !== WorkflowNodeType.Upload) {
    console.warn(`[certimate] current workflow node type is not: ${WorkflowNodeType.Upload}`);
  }

  const { t } = useTranslation();

  const { updateNode } = useWorkflowStore(useZustandShallowSelector(["updateNode"]));

  const formRef = useRef<UploadNodeConfigFormInstance>(null);
  const [formPending, setFormPending] = useState(false);

  const [drawerOpen, setDrawerOpen] = useState(false);
  const getFormValues = () => formRef.current!.getFieldsValue() as WorkflowNodeConfigForUpload;

  const wrappedEl = useMemo(() => {
    if (node.type !== WorkflowNodeType.Upload) {
      console.warn(`[certimate] current workflow node type is not: ${WorkflowNodeType.Upload}`);
    }

    if (!node.validated) {
      return <Typography.Link>{t("workflow_node.action.configure_node")}</Typography.Link>;
    }

    const config = (node.config as WorkflowNodeConfigForUpload) ?? {};
    return (
      <Flex className="size-full overflow-hidden" align="center" gap={8}>
        <Typography.Text className="truncate">{config.domains ?? ""}</Typography.Text>
      </Flex>
    );
  }, [node]);

  const handleDrawerConfirm = async () => {
    setFormPending(true);
    try {
      await formRef.current!.validateFields();
    } catch (err) {
      setFormPending(false);
      throw err;
    }

    try {
      const newValues = getFormValues();
      const newNode = produce(node, (draft) => {
        draft.config = {
          ...newValues,
        };
        draft.validated = true;
      });
      await updateNode(newNode);
    } finally {
      setFormPending(false);
    }
  };

  return (
    <>
      <SharedNode.Block node={node} disabled={disabled} onClick={() => setDrawerOpen(true)}>
        {wrappedEl}
      </SharedNode.Block>

      <SharedNode.ConfigDrawer
        node={node}
        open={drawerOpen}
        pending={formPending}
        onConfirm={handleDrawerConfirm}
        onOpenChange={(open) => setDrawerOpen(open)}
        getFormValues={() => formRef.current!.getFieldsValue()}
      >
        <UploadNodeConfigForm ref={formRef} disabled={disabled} initialValues={node.config} />
      </SharedNode.ConfigDrawer>
    </>
  );
};

export default memo(UploadNode);
