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
    <a-input-group>
      <a-input :model-value="modelValue" :placeholder="placeholder" readonly />
      <a-upload v-bind="uploadProps" @change="handleChange">
        <template #upload-button>
          <a-button type="primary">
            <icon-cloud />
          </a-button>
        </template>
      </a-upload>
    </a-input-group>
  </a-space>
</template>
