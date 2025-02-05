import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  CheckCircleOutlined as CheckCircleOutlinedIcon,
  ClockCircleOutlined as ClockCircleOutlinedIcon,
  CloseCircleOutlined as CloseCircleOutlinedIcon,
  DeleteOutlined as DeleteOutlinedIcon,
  PauseCircleOutlined as PauseCircleOutlinedIcon,
  PauseOutlined as PauseOutlinedIcon,
  SelectOutlined as SelectOutlinedIcon,
  SyncOutlined as SyncOutlinedIcon,
} from "@ant-design/icons";
import { useRequest } from "ahooks";
import { Button, Empty, Modal, Table, type TableProps, Tag, Tooltip, notification } from "antd";
import dayjs from "dayjs";
import { ClientResponseError } from "pocketbase";

import { cancelRun as cancelWorkflowRun } from "@/api/workflows";
import { WORKFLOW_TRIGGERS } from "@/domain/workflow";
import { WORKFLOW_RUN_STATUSES, type WorkflowRunModel } from "@/domain/workflowRun";
import { list as listWorkflowRuns, remove as removeWorkflowRun } from "@/repository/workflowRun";
import { getErrMsg } from "@/utils/error";
import WorkflowRunDetailDrawer from "./WorkflowRunDetailDrawer";

export type WorkflowRunsProps = {
  className?: string;
  style?: React.CSSProperties;
  workflowId: string;
};

const WorkflowRuns = ({ className, style, workflowId }: WorkflowRunsProps) => {
  const { t } = useTranslation();

  const [modalApi, ModelContextHolder] = Modal.useModal();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const tableColumns: TableProps<WorkflowRunModel>["columns"] = [
    {
      key: "$index",
      align: "center",
      fixed: "left",
      width: 50,
      render: (_, __, index) => (page - 1) * pageSize + index + 1,
    },
    {
      key: "id",
      title: t("workflow_run.props.id"),
      ellipsis: true,
      render: (_, record) => <span className="font-mono">{record.id}</span>,
    },
    {
      key: "status",
      title: t("workflow_run.props.status"),
      ellipsis: true,
      render: (_, record) => {
        if (record.status === WORKFLOW_RUN_STATUSES.PENDING) {
          return <Tag icon={<ClockCircleOutlinedIcon />}>{t("workflow_run.props.status.pending")}</Tag>;
        } else if (record.status === WORKFLOW_RUN_STATUSES.RUNNING) {
          return (
            <Tag icon={<SyncOutlinedIcon spin />} color="processing">
              {t("workflow_run.props.status.running")}
            </Tag>
          );
        } else if (record.status === WORKFLOW_RUN_STATUSES.SUCCEEDED) {
          return (
            <Tag icon={<CheckCircleOutlinedIcon />} color="success">
              {t("workflow_run.props.status.succeeded")}
            </Tag>
          );
        } else if (record.status === WORKFLOW_RUN_STATUSES.FAILED) {
          return (
            <Tag icon={<CloseCircleOutlinedIcon />} color="error">
              {t("workflow_run.props.status.failed")}
            </Tag>
          );
        } else if (record.status === WORKFLOW_RUN_STATUSES.CANCELED) {
          return (
            <Tag icon={<PauseCircleOutlinedIcon />} color="warning">
              {t("workflow_run.props.status.canceled")}
            </Tag>
          );
        }

        return <></>;
      },
    },
    {
      key: "trigger",
      title: t("workflow_run.props.trigger"),
      ellipsis: true,
      render: (_, record) => {
        if (record.trigger === WORKFLOW_TRIGGERS.AUTO) {
          return t("workflow_run.props.trigger.auto");
        } else if (record.trigger === WORKFLOW_TRIGGERS.MANUAL) {
          return t("workflow_run.props.trigger.manual");
        }

        return <></>;
      },
    },
    {
      key: "startedAt",
      title: t("workflow_run.props.started_at"),
      ellipsis: true,
      render: (_, record) => {
        if (record.startedAt) {
          return dayjs(record.startedAt).format("YYYY-MM-DD HH:mm:ss");
        }

        return <></>;
      },
    },
    {
      key: "endedAt",
      title: t("workflow_run.props.ended_at"),
      ellipsis: true,
      render: (_, record) => {
        if (record.endedAt) {
          return dayjs(record.endedAt).format("YYYY-MM-DD HH:mm:ss");
        }

        return <></>;
      },
    },
    {
      key: "$action",
      align: "end",
      fixed: "right",
      width: 120,
      render: (_, record) => {
        const allowCancel = record.status === WORKFLOW_RUN_STATUSES.PENDING || record.status === WORKFLOW_RUN_STATUSES.RUNNING;
        const aloowDelete =
          record.status === WORKFLOW_RUN_STATUSES.SUCCEEDED ||
          record.status === WORKFLOW_RUN_STATUSES.FAILED ||
          record.status === WORKFLOW_RUN_STATUSES.CANCELED;

        return (
          <Button.Group>
            <WorkflowRunDetailDrawer
              data={record}
              trigger={
                <Tooltip title={t("workflow_run.action.view")}>
                  <Button color="primary" icon={<SelectOutlinedIcon />} variant="text" />
                </Tooltip>
              }
            />

            <Tooltip title={t("workflow_run.action.cancel")}>
              <Button
                color="default"
                disabled={!allowCancel}
                icon={<PauseOutlinedIcon />}
                variant="text"
                onClick={() => {
                  handleCancelClick(record);
                }}
              />
            </Tooltip>

            <Tooltip title={t("workflow_run.action.delete")}>
              <Button
                color="danger"
                danger
                disabled={!aloowDelete}
                icon={<DeleteOutlinedIcon />}
                variant="text"
                onClick={() => {
                  handleDeleteClick(record);
                }}
              />
            </Tooltip>
          </Button.Group>
        );
      },
    },
  ];
  const [tableData, setTableData] = useState<WorkflowRunModel[]>([]);
  const [tableTotal, setTableTotal] = useState<number>(0);

  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);

  const {
    loading,
    error: loadedError,
    run: refreshData,
  } = useRequest(
    () => {
      return listWorkflowRuns({
        workflowId: workflowId,
        page: page,
        perPage: pageSize,
      });
    },
    {
      refreshDeps: [workflowId, page, pageSize],
      onSuccess: (res) => {
        setTableData(res.items);
        setTableTotal(res.totalItems);
      },
      onError: (err) => {
        if (err instanceof ClientResponseError && err.isAbort) {
          return;
        }

        console.error(err);
        notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });

        throw err;
      },
    }
  );

  const handleCancelClick = (workflowRun: WorkflowRunModel) => {
    modalApi.confirm({
      title: t("workflow_run.action.cancel"),
      content: t("workflow_run.action.cancel.confirm"),
      onOk: async () => {
        try {
          const resp = await cancelWorkflowRun(workflowId, workflowRun.id);
          if (resp) {
            refreshData();
          }
        } catch (err) {
          console.error(err);
          notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
        }
      },
    });
  };

  const handleDeleteClick = (workflowRun: WorkflowRunModel) => {
    modalApi.confirm({
      title: t("workflow_run.action.delete"),
      content: t("workflow_run.action.delete.confirm"),
      onOk: async () => {
        try {
          const resp = await removeWorkflowRun(workflowRun);
          if (resp) {
            setTableData((prev) => prev.filter((item) => item.id !== workflowRun.id));
            refreshData();
          }
        } catch (err) {
          console.error(err);
          notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
        }
      },
    });
  };

  return (
    <>
      {ModelContextHolder}
      {NotificationContextHolder}

      <div className={className} style={style}>
        <Table<WorkflowRunModel>
          columns={tableColumns}
          dataSource={tableData}
          loading={loading}
          locale={{
            emptyText: <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={loadedError ? getErrMsg(loadedError) : undefined} />,
          }}
          pagination={{
            current: page,
            pageSize: pageSize,
            total: tableTotal,
            showSizeChanger: true,
            onChange: (page: number, pageSize: number) => {
              setPage(page);
              setPageSize(pageSize);
            },
            onShowSizeChange: (page: number, pageSize: number) => {
              setPage(page);
              setPageSize(pageSize);
            },
          }}
          rowKey={(record: WorkflowRunModel) => record.id}
          scroll={{ x: "max(100%, 960px)" }}
        />
      </div>
    </>
  );
};

export default WorkflowRuns;
