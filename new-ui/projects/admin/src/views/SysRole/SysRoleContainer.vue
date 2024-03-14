<script lang="ts" setup>
import { Message } from "@arco-design/web-vue";
import SearchTable from "@/components/SearchTable/SearchTable.vue";
import type { SearchTableColumns } from "@/components/SearchTable/type";
import usePopup from "@/composables/usePopup";
import SysRoleForm from "./SysRoleForm.vue";
import { getList, save, remove } from "./api";

const columns: SearchTableColumns[] = [
  {
    title: "角色名称",
    dataIndex: "name",
    search: {
      valueType: "input",
    },
  },
  {
    title: "角色描述",
    dataIndex: "description",
  },
  {
    title: "创建时间",
    dataIndex: "created_at",
  },
  {
    title: "操作",
    slotName: "actions",
    width: 120,
  },
];

const openFormModal = usePopup(SysRoleForm, {
  nodeProps: ([record]) => ({ record }),
  popupProps: ([record, reload], exposed) => ({
    title: `${record?.id ? "编辑" : "新增"}角色`,
    width: 800,
    onBeforeOk: async (done) => {
      await exposed()?.handleSubmit(save, {
        id: record?.id,
      });
      Message.success("操作成功");
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
      <a-button type="primary" @click="openFormModal({}, reload)">
        <template #icon>
          <icon-plus />
        </template>
        新增
      </a-button>
    </template>
    <template #actions="{ record, reload }">
      <a-link @click="openFormModal(record, reload)">编辑</a-link>
      <a-popconfirm
        content="是否删除？"
        position="left"
        type="warning"
        :on-before-ok="() => handleRemove(record.id, reload)"
      >
        <a-link status="danger">删除</a-link>
      </a-popconfirm>
    </template>
  </SearchTable>
</template>
./SysRoleForm.vue
