<template>
  <a-form ref="formRef" :model="formData" auto-label-width>
    <a-form-item field="username" label="账号" :rules="[{ required: true, message: '请输入账号' }]">
      <a-input v-model="formData.username" placeholder="请输入账号" />
    </a-form-item>
    <a-form-item
      v-if="!props.record?.id"
      field="password"
      label="密码"
      :rules="[{ required: true, message: '请输入密码' }]"
    >
      <a-input-password
        v-model="formData.password"
        placeholder="请输入密码"
        autocomplete="new-password"
      />
    </a-form-item>
    <a-form-item field="status" label="启用状态">
      <a-switch v-model="formData.status" />
    </a-form-item>
  </a-form>
</template>

<script lang="ts" setup>
import useSubmit from "@/composables/useSubmit";
const props = defineProps({
  record: Object,
});

const { formRef, formData, handleSubmit } = useSubmit({
  username: "",
  password: "",
  status: true,
});

defineExpose({
  handleSubmit,
});
</script>
