<script lang="ts" setup>
import { computed } from "vue";
import type { UploadInstance, FileItem } from "@arco-design/web-vue";
import { uploadUrl } from "@/http/config";

defineProps({
  modelValue: String,
  placeholder: String,
});

const emits = defineEmits(["update:modelValue"]);

const uploadProps = computed<UploadInstance["$props"]>(() => ({
  action: uploadUrl,
  name: "file",
  headers: { [__AUTH_KEY]: localStorage.getItem(__AUTH_KEY) },
  showFileList: false,
}));

const handleChange = (_, file: FileItem) => {
  console.log(file.response);
};
</script>
<template>
  <a-space>
    <a-input :model-value="modelValue" :placeholder="placeholder" readonly>
      <template #append>
        <a-upload v-bind="uploadProps" @change="handleChange">
          <template #upload-button>
            <icon-upload />
          </template>
        </a-upload>
      </template>
    </a-input>
  </a-space>
</template>
