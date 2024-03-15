<template>
  <a-form ref="formRef" :model="form" :style="{ width: '600px' }" @submit="handleSubmit">
    <a-form-item
      field="name"
      label="角色名称"
      :rules="[{ required: true, message: '请输入角色名称' }]"
      :validate-trigger="['change', 'input']"
    >
      <a-input v-model="form.name" placeholder="请输入角色名称" />
    </a-form-item>
    <a-form-item
      field="key"
      label="角色标志"
      :rules="[{ required: true, message: '请输入角色标志' }]"
      :validate-trigger="['change', 'input']"
      showable
    >
      <a-input v-model="form.key" placeholder="请输入角色标志" />
    </a-form-item>
    <a-form-item
      field="icon"
      label="角色图标"
      :rules="[{ required: true, message: '请输入角色图标' }]"
      :validate-trigger="['change', 'input']"
    >
      <CustomUploader v-model="form.icon" placeholder="请输入角色图标" />
    </a-form-item>
    <a-form-item
      field="hello_msg"
      label="打招呼信息"
      :rules="[{ required: true, message: '请输入打招呼信息' }]"
      :validate-trigger="['change', 'input']"
    >
      <a-input v-model="form.hello_msg" placeholder="请输入打招呼信息" />
    </a-form-item>
    <a-form-item field="username" label="上下文信息" :validate-trigger="['change', 'input']">
      <a-table :data="form.context || []" :pagination="false">
        <template #columns>
          <a-table-column title="对话角色">
            <template #cell="{ record }">
              <a-input v-model="record.role" />
            </template>
          </a-table-column>
          <a-table-column width="350">
            <template #title>
              <div class="content-title">
                <span>对话内容</span>
                <a-button @click="addContext" type="primary">增加一行</a-button>
              </div>
            </template>
            <template #cell="{ record }">
              <div class="content-cell">
                <a-input v-model="record.content" /><icon-minus-circle
                  @click="removeContext(record)"
                />
              </div>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-form-item>
    <a-form-item field="enable" label="启用状态">
      <a-switch v-model="form.enable" />
    </a-form-item>
  </a-form>
</template>

<script setup>
import { ref, defineExpose, defineProps } from "vue";
import CustomUploader from "@/components/CustomUploader.vue";
const props = defineProps({
  data: {},
});

const formRef = ref();
const form = ref({});
if (props.data?.id) {
  form.value = Object.assign({}, props.data);
}

const addContext = () => {
  form.value.context = form.value.context || [];
  form.value.context.push({
    role: "",
    content: "",
  });
};

const removeContext = (record) => {
  const index = form.value.context.findIndex((item) => {
    return item === record;
  });
  if (index > -1) {
    form.value.context.splice(index, 1);
  }
};
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
