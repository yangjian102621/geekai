<script lang="ts" setup>
import { getList, save, deleting, setStatus } from "./api";
import { reactive, ref } from "vue";
import RoleForm from "./RoleForm.vue";
import useCustomFormPopup from "@/composables/useCustomFormPopup";
import { Message } from "@arco-design/web-vue";
import SimpleTable from "@/components/SimpleTable/SimpleTable.vue";
// table 配置
const columns = [
  {
    title: "角色名称",
    dataIndex: "name",
  },
  {
    title: "角色标识",
    dataIndex: "key",
  },
  {
    title: "启用状态",
    dataIndex: "enable",
    slotName: "status",
  },
  {
    title: "角色图标",
    dataIndex: "icon",
    slotName: "icon",
  },
  {
    title: "打招呼信息",
    dataIndex: "hello_msg",
  },
  {
    title: "操作",
    slotName: "action",
  },
];

const expandable = reactive({
  title: "",
  width: 80,
});

// 数据
const tableData = ref([]);
const getData = () => {
  getList().then(({ code, data }) => {
    if (code === 0) {
      tableData.value = data;
    }
  });
};
getData();

//展开行table
const expandColumns = [
  {
    dataIndex: "role",
    title: "对话角色",
    width: 120,
  },
  {
    dataIndex: "content",
    title: "对话内容",
  },
];

//  新增编辑
const popup = useCustomFormPopup(RoleForm, save, {
  popupProps: (arg) => ({ title: arg[0].record ? "编辑角色" : "新增角色" }),
});

// 删除
const handleDelete = ({ id }, reload) => {
  deleting(id).then(({ code }) => {
    if (code === 0) {
      Message.success("操作成功");
      reload();
    }
  });
};

// 状态
const handleStatusChange = ({ value, record, reload }) => {
  setStatus({
    id: record.id,
    value,
    filed: "enable",
  }).then(({ code }) => {
    if (code === 0) {
      Message.success("操作成功");
      reload();
    }
  });
};
</script>
<template>
  <SimpleTable :columns="columns" :request="getList" :expandable="expandable">
    <template #expand-row="{ record }">
      <a-table :columns="expandColumns" :data="record.context || []" :pagination="false"></a-table>
    </template>
    <template #action="{ record, reload }">
      <a-link @click="popup({ record, reload })">编辑</a-link>
      <a-popconfirm content="确定删除？" @ok="handleDelete(record, reload)">
        <a-link status="danger">删除</a-link>
      </a-popconfirm>
    </template>
    <template #header="{ reload }">
      <a-button @click="popup({ reload })" size="small"><icon-plus />新增</a-button>
    </template>
    <template #status="{ record, reload }">
      <a-switch
        v-model="record.enable"
        @change="
          (value) => {
            handleStatusChange({ value, record, reload });
          }
        "
      />
    </template>

    <template #icon="{ record }">
      <a-avatar>
        <img alt="avatar" :src="record.icon" />
      </a-avatar>
    </template>
  </SimpleTable>
</template>
