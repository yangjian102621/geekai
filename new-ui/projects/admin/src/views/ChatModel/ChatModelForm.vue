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
      field="value"
      label="模型值"
      :rules="[{ required: true, message: '请输入名称' }]"
      :validate-trigger="['change', 'input']"
      showable
    >
      <a-input v-model="form.value" placeholder="请输入名称" />
    </a-form-item>
    <a-form-item
      field="weight"
      label="对话权重"
      :rules="[{ required: true, message: '请输入对话权重' }]"
      :validate-trigger="['change', 'input']"
      showable
    >
      <a-input-number v-model="form.weight" placeholder="请输入对话权重" />
      <a-tooltip content="对话权重，每次对话扣减多少次对话额度" position="right">
        <icon-info-circle-fill />
      </a-tooltip>
    </a-form-item>

    <a-form-item field="open" label="开放状态代理">
      <a-switch v-model="form.open" />
    </a-form-item>
    <a-form-item field="enabled" label="启用状态">
      <a-switch v-model="form.enabled" />
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
