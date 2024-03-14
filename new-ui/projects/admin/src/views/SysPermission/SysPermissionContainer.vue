<script lang="ts" setup>
import type { TableColumnData } from "@arco-design/web-vue";
import { Message } from "@arco-design/web-vue";
import useRequest from "@/composables/useRequest";
import usePopup from "@/composables/usePopup";
import SysPermissionForm from "./SysPermissionForm.vue";
import { getList, save, remove } from "./api";

const columns: TableColumnData[] = [
  {
    dataIndex: "name",
    title: "权限名称",
  },
  {
    dataIndex: "slug",
    title: "权限标识",
  },
  {
    dataIndex: "sort",
    title: "排序",
  },
  {
    title: "操作",
    slotName: "actions",
  },
];

const [reload, data, loading] = useRequest(getList);

const openFormModal = usePopup(SysPermissionForm, {
  nodeProps: ([_, record]) => ({ record, options: data.value }),
  popupProps: ([reload, record], exposed) => ({
    title: `${record?.id ? "编辑" : "新增"}权限`,
    onBeforeOk: async (done) => {
      await exposed()?.handleSubmit(save, {
        id: record?.id,
      });
      await reload();
      done(true);
    },
  }),
});

const handleRemove = async (id) => {
  await remove({ id });
  Message.success("删除成功");
  await reload();
  return true;
};

reload();
</script>
<template>
  <div style="padding-bottom: 10px; text-align: right">
    <a-button type="primary" :loading="loading" @click="openFormModal(reload, {})">
      <template #icon>
        <icon-plus />
      </template>
      新增
    </a-button>
  </div>
  <a-table
    :key="String(loading)"
    :data="data"
    :loading="loading"
    :columns="columns"
    :pagination="false"
    default-expand-all-rows
  >
    <template #actions="{ record, reload }">
      <a-link @click="openFormModal(reload, record)">编辑</a-link>
      <a-popconfirm
        content="是否删除？"
        position="left"
        type="warning"
        :on-before-ok="() => handleRemove(record.id)"
      >
        <a-link status="danger">删除</a-link>
      </a-popconfirm>
    </template>
  </a-table>
</template>
