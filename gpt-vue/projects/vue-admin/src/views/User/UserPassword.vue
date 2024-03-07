<template>
  <a-form ref="formRef" :model="form" :style="{ width: '600px' }" @submit="handleSubmit">
    <a-form-item
      field="username"
      label="账号"
      :rules="[{ required: true, message: 'name is required' }]"
      :validate-trigger="['change', 'input']"
      :disabled="true"
    >
      <a-input v-model="form.username" placeholder="请输入账号" />
    </a-form-item>
    <a-form-item
      field="password"
      label="新密码"
      :rules="[{ required: true, message: 'password is required' }]"
      :validate-trigger="['change', 'input']"
      showable
    >
      <a-input v-model="form.password" placeholder="请输入密码" />
    </a-form-item>
  </a-form>
</template>

<script setup>
import { ref, defineExpose, defineProps } from "vue";
import { getModel, getRole } from "./api";
const props = defineProps({
  data: {},
});

const formRef = ref();
const form = ref({
  id: "",
  username: "",
  password: "",
});
form.value.id = props.data.id;
form.value.username = props.data.username;

defineExpose({
  formRef,
  form,
});
</script>
