<template>
  <a-form ref="formRef" :model="formData" autoLabelWidth>
    <a-formItem
      label="角色名称"
      field="name"
      :rules="[{ required: true, message: '请输入角色名称' }]"
    >
      <a-input v-model="formData.name" />
    </a-formItem>
    <a-formItem
      label="角色描述"
      field="description"
      :rules="[{ required: true, message: '请输入角色描述' }]"
    >
      <a-input v-model="formData.description" />
    </a-formItem>
    <a-formItem
      label="平台权限"
      field="permissions"
      :rules="[{ required: true, message: '请选择平台权限' }]"
    >
      <a-tree
        v-model:checked-keys="formData.permissions"
        :data="options"
        :field-names="{ key: 'id', title: 'name' }"
      />
    </a-formItem>
  </a-form>
</template>

<script lang="ts" setup>
import { getList } from "@/views/SysPermission/api";
import useRequest from "@/composables/useRequest";
import useSubmit from "@/composables/useSubmit";

const props = defineProps({
  record: Object,
});

const [getOptions, options, loading] = useRequest(getList);

const { formRef, formData, handleSubmit } = useSubmit({
  name: "",
  description: "",
  permissions: [],
});

Object.assign(formData, props.record);
getOptions();
defineExpose({ handleSubmit });
</script>
