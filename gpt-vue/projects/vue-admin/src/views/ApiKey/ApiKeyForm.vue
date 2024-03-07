<template>
  <a-form ref="formRef" :model="form" :style="{ width: '600px' }" @submit="handleSubmit">
    <a-form-item
      field="platform"
      label="所属平台"
      :rules="[{ required: true, message: '请输入所属平台' }]"
      :validate-trigger="['change', 'input']"
    >
      <a-input v-model="form.platform" placeholder="请输入所属平台" />
    </a-form-item>
    <a-form-item
      field="name"
      label="名称"
      :rules="[{ required: true, message: '请输入名称' }]"
      :validate-trigger="['change', 'input']"
      showable
    >
      <a-input v-model="form.name" placeholder="请输入名称" />
    </a-form-item>
    <a-form-item
      field="type"
      label="用途"
      :rules="[{ required: true, message: '请输入用途' }]"
      :validate-trigger="['change', 'input']"
    >
      <a-select v-model="form.type" placeholder="请输入用途" :options="typeOPtions"> </a-select>
    </a-form-item>
    <a-form-item
      field="value"
      label="API KEY"
      :rules="[{ required: true, message: '请输入API KEY' }]"
      :validate-trigger="['change', 'input']"
    >
      <a-input v-model="form.value" placeholder="请输入API KEY" />
    </a-form-item>
    <a-form-item
      field="api_url"
      label="API URL"
      :rules="[{ required: true, message: '请输入API URL' }]"
      :validate-trigger="['change', 'input']"
    >
      <a-input v-model="form.api_url" placeholder="请输入API URL" />
    </a-form-item>

    <a-form-item field="use_proxy" label="使用代理">
      <a-switch v-model="form.use_proxy" />
    </a-form-item>
    <a-form-item field="enable" label="启用状态">
      <a-switch v-model="form.enable" />
    </a-form-item>
  </a-form>
</template>

<script setup>
import { ref, defineExpose, defineProps } from "vue";
const props = defineProps({
  data: {},
});

const formRef = ref();
const form = ref({});
if (props.data?.id) {
  form.value = Object.assign({}, props.data);
}

defineExpose({
  formRef,
  form,
});

const typeOPtions = [
  {
    label: "聊天",
    value: "chart",
  },
  {
    label: "绘图",
    value: "img",
  },
];
</script>
<style lang="less" scoped>
.content-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 350px;
}
.content-cell {
  display: flex;
  align-items: center;
  width: 350px;
  svg {
    margin-left: 10px;
  }
}
</style>
