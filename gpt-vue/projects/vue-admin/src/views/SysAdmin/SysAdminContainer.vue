<script lang="ts" setup>
import { Message } from "@arco-design/web-vue";
import { dateFormat } from "@gpt-vue/packages/utils";
import SearchTable from "@/components/SearchTable/SearchTable.vue";
import type { SearchTableColumns } from "@/components/SearchTable/type";
import ConfirmSwitch from "@/components/ConfirmSwitch.vue";
import usePopup from "@/composables/usePopup";
import SysAdminForm from "./SysAdminForm.vue";
import SysAdminResetPWD from "./SysAdminResetPWD.vue";
import { getList, save, remove, resetPass } from "./api";

const columns: SearchTableColumns[] = [
  {
    dataIndex: "username",
    title: "账号",
    search: {
      valueType: "input",
    },
  },
  {
    dataIndex: "last_login_ip",
    title: "最后登录IP",
  },
  {
    dataIndex: "last_login_at",
    title: "最后登录时间",
    render: ({ record }) => dateFormat(record.last_login_at),
  },
  {
    dataIndex: "created_at",
    title: "创建时间",
    render: ({ record }) => dateFormat(record.created_at),
  },
  {
    dataIndex: "updated_at",
    title: "更新时间",
    render: ({ record }) => dateFormat(record.updated_at),
  },
  {
    dataIndex: "status",
    title: "状态",
    slotName: "switch",
    width: 80,
  },
  {
    title: "操作",
    slotName: "actions",
    fixed: "right",
    width: 180,
  },
];

const openFormModal = usePopup(SysAdminForm, {
  nodeProps: ([_, record]) => ({ record }),
  popupProps: ([reload, record], exposed) => ({
    title: `${record?.id ? "编辑" : "新增"}系统管理员`,
    onBeforeOk: async (done) => {
      await exposed()?.handleSubmit(save, {
        id: record?.id,
      });
      await reload();
      done(true);
    },
  }),
});

const openResetPWDModal = usePopup(SysAdminResetPWD, {
  popupProps: ([reload, record], exposed) => ({
    title: `修改密码`,
    onBeforeOk: async (done) => {
      await exposed()?.handleSubmit(resetPass, {
        id: record?.id,
      });
      Message.success("修改成功");
      await reload();
      done(true);
    },
  }),
});

const handleRemove = async (id, reload) => {
  await remove({ id });
  Message.success("删除成功");
  await reload();
  return true;
};
</script>
<template>
  <SearchTable :request="getList" :columns="columns">
    <template #header-option="{ reload }">
      <a-button type="primary" @click="openFormModal(reload, {})">
        <template #icon> <icon-plus /> </template>
        新增
      </a-button>
    </template>
    <template #switch="{ record, column }">
      <ConfirmSwitch
        v-model="record[column.dataIndex]"
        :api="async () => save({ ...record, status: !record.status })"
      />
    </template>
    <template #actions="{ record, reload }">
      <a-link @click="openResetPWDModal(reload, record)">修改密码</a-link>
      <template v-if="record.id !== 1">
        <a-link @click="openFormModal(reload, record)">编辑</a-link>
        <a-popconfirm
          content="是否删除？"
          position="left"
          type="warning"
          :on-before-ok="() => handleRemove(record.id, reload)"
        >
          <a-link status="danger">删除</a-link>
        </a-popconfirm>
      </template>
    </template>
  </SearchTable>
</template>
