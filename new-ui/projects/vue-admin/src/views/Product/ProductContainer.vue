<script lang="ts" setup>
import { getList, save, deleting, setStatus } from "./api";
import { ref } from "vue";
import ProductForm from "./ProductForm.vue";
import useCustomFormPopup from "@/composables/useCustomFormPopup";
import { Message } from "@arco-design/web-vue";
import SimpleTable from "@/components/SimpleTable/SimpleTable.vue";
import { dateFormat } from "@gpt-vue/packages/utils";
// table 配置
const columns = [
  {
    title: "产品名称",
    dataIndex: "name",
  },
  {
    title: "产品价格",
    dataIndex: "price",
  },
  {
    title: "优惠金额",
    dataIndex: "discount",
  },
  {
    title: "有效期（天）",
    dataIndex: "days",
  },
  {
    title: "对话次数",
    dataIndex: "calls",
  },
  {
    title: "绘图次数",
    dataIndex: "img_calls",
  },
  {
    title: "销量",
    dataIndex: "sales",
  },
  {
    title: "启用状态",
    dataIndex: "enabled",
    slotName: "status",
    align: "center",
    width: 100
  },
  {
    title: "更新时间",
    dataIndex: "updated_at",
    width: 180,
    render: ({ record }) => {
      return dateFormat(record.updated_at);
    },
  },
  {
    title: "操作",
    slotName: "action",
    width: 120,
    fixed: "right"
  },
];

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

//  新增编辑
const popup = useCustomFormPopup(ProductForm, save, {
  popupProps: (arg) => ({ title: arg[0].record ? "编辑产品" : "新增产品" }),
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
    enabled: value,
  }).then(({ code }) => {
    if (code === 0) {
      Message.success("操作成功");
      reload();
    }
  });
};
</script>
<template>
  <SimpleTable :columns="columns" :request="getList">
    <template #action="{ record, reload }">
      <a-link @click="popup({ record, reload })">编辑</a-link>
      <a-popconfirm content="确定删除？" @ok="handleDelete(record, reload)">
        <a-link status="danger">删除</a-link>
      </a-popconfirm>
    </template>
    <template #header="{ reload }">
      <a-button @click="popup({ reload })" size="small" type="primary"
        ><template #icon> <icon-plus /> </template>新增
      </a-button>
    </template>
    <template #status="{ record, reload }">
      <a-switch
        v-model="record.enabled"
        @change="
          (value) => {
            handleStatusChange({ value, record, reload });
          }
        "
      />
    </template>
  </SimpleTable>
</template>
