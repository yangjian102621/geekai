<template>
  <a-form ref="formRef" :model="form" :style="{ width: '600px' }" @submit="handleSubmit">
    <a-form-item
      field="name"
      label="产品名称"
      :rules="[{ required: true, message: '请输入产品名称' }]"
      :validate-trigger="['change', 'input']"
    >
      <a-input v-model="form.name" placeholder="请输入产品名称" />
    </a-form-item>
    <a-form-item
      field="price"
      label="产品价格"
      :rules="[{ required: true, message: '请输入产品价格' }]"
      :validate-trigger="['change', 'input']"
      showable
    >
      <a-input-number v-model="form.price" placeholder="请输入产品价格" />
    </a-form-item>
    <a-form-item
      field="discount"
      label="优惠金额"
      :rules="[{ required: true, message: '请输入优惠金额' }]"
      :validate-trigger="['change', 'input']"
      showable
    >
      <a-input-number v-model="form.discount" placeholder="请输入优惠金额" />
    </a-form-item>
    <a-form-item
      field="days"
      label="有效期（天）"
      :rules="[{ required: true, message: '请输入有效期（天）' }]"
      :validate-trigger="['change', 'input']"
      showable
    >
      <a-input-number v-model="form.days" placeholder="请输入有效期（天）" />
    </a-form-item>
    <a-form-item field="calls" label="对话次数" :validate-trigger="['change', 'input']" showable>
      <a-input-number v-model="form.calls" placeholder="请输入对话次数" />
    </a-form-item>
    <a-form-item
      field="img_calls"
      label="绘图次数"
      :validate-trigger="['change', 'input']"
      showable
    >
      <a-input-number v-model="form.img_calls" placeholder="请输入绘图次数" />
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
