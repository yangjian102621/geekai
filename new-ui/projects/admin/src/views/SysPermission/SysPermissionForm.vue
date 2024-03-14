<template>
  <a-form ref="formRef" :model="formData" auto-label-width>
    <a-form-item
      field="name"
      label="权限名称"
      :rules="[{ required: true, message: '请输入权限名称' }]"
    >
      <a-input v-model="formData.name" placeholder="请输入权限名称" />
    </a-form-item>
    <a-form-item
      field="pid"
      label="上级权限"
      :rules="[{ required: true, message: '请选择上级权限' }]"
      extra="默认为顶级权限"
    >
      <a-tree-select
        v-model="formData.pid"
        :data="_options"
        :field-names="{ key: 'id', title: 'name' }"
        placeholder="请选择上级权限"
      />
    </a-form-item>
    <a-form-item field="slug" label="权限标识" extra="规则：controller_method (注意大小写)">
      <a-input v-model="formData.slug" placeholder="请输入权限标识" />
    </a-form-item>
    <a-form-item
      field="sort"
      label="排序"
      :rules="[
        { required: true, message: '请输入排序' },
        { min: 0, message: '请输入大于0的正整数', type: 'number' },
      ]"
      extra="数字越小越靠前"
    >
      <a-input-number v-model="formData.sort" placeholder="请输入排序" :precision="0" />
    </a-form-item>
  </a-form>
</template>

<script lang="ts" setup>
import { computed, unref } from "vue";
import useSubmit from "@/composables/useSubmit";

const props = defineProps({
  record: Object,
  options: Array,
});

const { formRef, formData, handleSubmit } = useSubmit({
  name: "",
  pid: "0",
  slug: "",
  sort: null,
});

const _options = computed(() => {
  return [{ id: "0", name: "顶部权限", children: unref(props.options ?? []) }];
});

Object.assign(formData, props.record);

defineExpose({ handleSubmit });
</script>
