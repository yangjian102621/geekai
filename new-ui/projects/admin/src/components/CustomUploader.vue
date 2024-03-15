<script lang="ts" setup>
import { computed } from "vue";
import { Message } from "@arco-design/web-vue";
import type { UploadInstance, FileItem } from "@arco-design/web-vue";
import { uploadUrl } from "@/http/config";

defineProps({
  modelValue: String,
  placeholder: String,
});

const emits = defineEmits(["update:modelValue"]);

const uploadProps = computed<UploadInstance["$props"]>(() => {
  const TOKEN = JSON.parse(localStorage.getItem(__AUTH_KEY))?.token;
  return {
    accept: "image/*",
    action: uploadUrl,
    name: "file",
    headers: { [__AUTH_KEY]: TOKEN },
    showFileList: false,
  };
});

const handleChange = (_, file: FileItem) => {
  if (file?.response) {
    emits("update:modelValue", file?.response?.data?.url);
    Message.success("上传成功");
  }
};
</script>
<template>
  <a-upload v-bind="uploadProps" style="width: 100%" @change="handleChange">
    <template #upload-button>
      <a-input-group style="width: 100%">
        <a-input
          :model-value="modelValue"
          :placeholder="placeholder"
          readonly
          allow-clear
          @clear.stop="emits('update:modelValue')"
        />
        <a-button type="primary" style="width: 100px">
          <icon-cloud />
        </a-button>
      </a-input-group>
    </template>
  </a-upload>
</template>
