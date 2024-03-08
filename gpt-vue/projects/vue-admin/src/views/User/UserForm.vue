<template>
  <a-form ref="formRef" :model="form" :style="{ width: '600px' }" @submit="handleSubmit">
    <a-form-item
      field="username"
      label="账号"
      :rules="[{ required: true, message: '请输入账号' }]"
      :validate-trigger="['change', 'input']"
    >
      <a-input v-model="form.username" placeholder="请输入账号" />
    </a-form-item>
    <a-form-item
      v-if="!props.data.id"
      field="password"
      label="密码"
      :rules="[{ required: true, message: '请输入密码' }]"
      :validate-trigger="['change', 'input']"
      showable
    >
      <a-input v-model="form.password" placeholder="请输入密码" />
    </a-form-item>
    <a-form-item
      field="calls"
      label="对话次数"
      :rules="[{ required: true, message: '请输入对话次数' }]"
    >
      <a-input-number v-model="form.calls" placeholder="请输入对话次数" />
    </a-form-item>
    <a-form-item
      field="img_calls"
      label="绘图次数"
      :rules="[{ required: true, message: '请输入绘图次数' }]"
    >
      <a-input-number v-model="form.img_calls" placeholder="请输入绘图次数" />
    </a-form-item>
    <a-form-item field="expired_time" label="有效期">
      <a-date-picker v-model="form.expired_time" placeholder="请选择有效期" />
    </a-form-item>
    <a-form-item field="chat_roles" label="聊天角色">
      <a-select
        :field-names="{ value: 'key', label: 'name' }"
        v-model="form.chat_roles"
        placeholder="请选择聊天角色"
        multiple
        :options="roleOption"
        :rules="[{ required: true, message: '请选择聊天角色' }]"
      >
      </a-select>
    </a-form-item>
    <a-form-item field="chat_models" label="模型角色">
      <a-select
        :field-names="{ value: 'value', label: 'name' }"
        v-model="form.chat_models"
        placeholder="请选择模型角色"
        multiple
        :options="modalOption"
        :rules="[{ required: true, message: '请选择模型角色' }]"
      >
      </a-select>
    </a-form-item>
    <a-form-item field="status" label="启用状态">
      <a-switch v-model="form.status" />
    </a-form-item>
    <a-form-item field="vip" label="开通VIP">
      <a-switch v-model="form.vip" />
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
  username: "",
  password: "",
  calls: "",
  img_calls: "",
  expired_time: "",
  chat_roles: [],
  chat_models: [],
  status: false,
  vip: false,
});
if (props.data?.id) {
  form.value = Object.assign({}, props.data);
  if (form.value.expired_time === 0) {
    form.value.expired_time = "";
  }
}

//拿选项
const modalOption = ref([]);
const roleOption = ref([]);
const getOption = (api, container) => {
  api().then(({ code, data }) => {
    if (code === 0) {
      container.value = data;
    }
  });
};
getOption(getModel, modalOption);
getOption(getRole, roleOption);

defineExpose({
  formRef,
  form,
});
</script>
