<script lang="ts" setup>
import { ref, watchEffect } from "vue";
import useSubmit from "@/composables/useSubmit";
import FunctionsFormTable from "./FunctionsFormTable.vue";
import translateTableData from "./translateTableData";

const props = defineProps({
  record: Object,
});

const tableData = ref([]);
const { formRef, formData, handleSubmit, submitting } = useSubmit({
  name: "",
  label: "",
  description: "",
  action: "",
  token: "",
  parameters: {},
  enabled: false,
});

const rules = {
  name: [{ required: true, message: "请输入函数名称" }],
  label: [{ required: true, message: "请输入函数标签" }],
  description: [{ required: true, message: "请输入函数功能描述" }],
};

watchEffect(() => {
  Object.assign(formData, props.record ?? {});
  tableData.value = translateTableData.get(formData.parameters);
});

defineExpose({
  handleSubmit,
  parameters: () => translateTableData.set(tableData.value),
});
</script>
<template>
  <a-spin :loading="submitting" style="width: 100%">
    <a-form ref="formRef" :model="formData" auto-label-width :rules="rules">
      <a-form-item field="name" label="函数名称">
        <a-input v-model="formData.name" placeholder="函数名称最好为英文" />
      </a-form-item>
      <a-form-item field="label" label="函数标签">
        <a-input v-model="formData.label" placeholder="函数的中文名称" />
      </a-form-item>
      <a-form-item field="description" label="功能描述">
        <a-input v-model="formData.description" placeholder="函数的中文名称" />
      </a-form-item>

      <a-form-item field="parameters" label="函数参数">
        <FunctionsFormTable v-model="tableData" />
      </a-form-item>

      <a-form-item field="action" label="API 地址">
        <a-input v-model="formData.action" placeholder="该函数实现的API地址，可以是第三方服务API" />
      </a-form-item>
      <a-form-item field="token" label="API Token">
        <a-input-search v-model="formData.token" placeholder="API授权Token">
          <template #append>
            <a-tooltip
              content="只有本地服务才可以使用自动生成Token第三方服务请填写第三方服务API Token"
            >
              <a-button>生成Token</a-button>
            </a-tooltip>
          </template>
        </a-input-search>
      </a-form-item>
      <a-form-item field="enabled" label="启用状态">
        <a-switch v-model="formData.enabled" />
      </a-form-item>
    </a-form>
  </a-spin>
</template>
