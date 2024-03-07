<script lang="ts" setup>
import SearchTable from "@/components/SearchTable/SearchTable.vue";
import type { SearchTableColumns } from "@/components/SearchTable/type";
import { getList, save as saveApi, deletApi, resetPassword } from "./api";
import UserForm from "./UserForm.vue";
import { Message } from "@arco-design/web-vue";
import { dateFormat } from "@gpt-vue/packages/utils";
import UserPassword from "./UserPassword.vue";
import useCustomFormPopup from "@/composables/useCustomFormPopup";
const columns: SearchTableColumns[] = [
  {
    title: "账号",
    dataIndex: "username",
    search: {
      valueType: "input",
    },
  },
  {
    title: "剩余对话次数",
    dataIndex: "calls",
  },
  {
    title: "剩余绘图次数",
    dataIndex: "img_calls",
  },
  {
    title: "累计消耗tokens",
    dataIndex: "total_tokens",
  },
  {
    title: "状态",
    dataIndex: "status",
    render: ({ record }) => {
      return record.status ? "正常" : "停用";
    },
  },
  {
    title: "过期时间",
    dataIndex: "expired_time",
    render: ({ record }) => {
      return dateFormat(record.expired_time);
    },
  },
  {
    title: "注册时间",
    dataIndex: "created_at",
    render: ({ record }) => {
      return dateFormat(record.created_at);
    },
  },
  {
    title: "操作",
    slotName: "actions",
  },
];

//弹窗
const editModal = useCustomFormPopup(UserForm, saveApi);
const password = useCustomFormPopup(UserPassword, resetPassword);

const handleDelete = async ({ id }: { id: string }, reload) => {
  const res = await deletApi(id);
  if (res.code === 0) {
    Message.success("操作成功");
    reload();
  }
};
</script>
<template>
  <SearchTable :request="getList" :columns="columns">
    <template #actions="{ record, reload }">
      <a-link @click="editModal({ record, reload })">编辑</a-link>
      <a-popconfirm content="确定删除？" @ok="handleDelete(record, reload)">
        <a-link>删除</a-link>
      </a-popconfirm>
      <a-link @click="password({ record, reload })">重置密码</a-link>
    </template>
    <template #search-extra="{ reload }">
      <a-button @click="editModal({ reload })" status="success" size="small"
        ><icon-plus />新增用户</a-button
      >
    </template>
  </SearchTable>
</template>
