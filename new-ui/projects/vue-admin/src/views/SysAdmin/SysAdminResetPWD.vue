<template>
  <a-form ref="formRef" :model="formData" auto-label-width>
    <a-form-item
      field="password"
      label="密码"
      :rules="[{ required: true, message: '请输入密码' }]"
      showable
    >
      <a-input-password
        v-model="formData.password"
        placeholder="请输入密码"
        autocomplete="new-password"
      />
    </a-form-item>
  </a-form>
</template>

<script lang="ts" setup>
import useSubmit from "@/composables/useSubmit";

const { formRef, formData, handleSubmit } = useSubmit({
  password: "",
});

defineExpose({
  handleSubmit,
});
</script>
