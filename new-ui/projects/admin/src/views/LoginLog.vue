<script lang="ts" setup>
import { dateFormat } from "@chatgpt-plus/packages/utils";
import SearchTable from "@/components/SearchTable/SearchTable.vue";
import type { SearchTableColumns } from "@/components/SearchTable/type";
import { loginLog } from "@/http/login";

const columns: SearchTableColumns[] = [
  {
    dataIndex: "username",
    title: "用户名",
  },
  {
    dataIndex: "login_ip",
    title: "登录IP",
  },
  {
    dataIndex: "login_address",
    title: "登录地址",
  },
  {
    dataIndex: "created_at",
    title: "登陆时间",
    render: ({ record }) => dateFormat(record.created_at),
  },
];
</script>
<template>
  <SearchTable :request="loginLog" :columns="columns" />
</template>
