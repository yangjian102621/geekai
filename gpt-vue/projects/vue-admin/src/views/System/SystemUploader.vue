<script lang="ts" setup>
import { computed, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import type { UploadInstance, FileItem } from "@arco-design/web-vue";
import { uploadUrl } from "@/http/config";

defineProps({
  modelValue: String,
  placeholder: String,
});

const emits = defineEmits(["update:modelValue"]);

const uploadRef = ref();

const uploadProps = computed<UploadInstance["$props"]>(() => ({
  action: uploadUrl,
  name: "file",
  headers: { [__AUTH_KEY]: localStorage.getItem(__AUTH_KEY) },
  showFileList: false,
}));

const handleChange = (_, file: FileItem) => {
  if (file?.response) {
    emits("update:modelValue", file?.response?.data?.url);
    Message.success("上传成功");
  }
};
</script>
<template>
  <a-input-group style="width: 100%">
    <a-input
      :model-value="modelValue"
      :placeholder="placeholder"
      readonly
      @dblclick="uploadRef?.$el?.click()"
    />
    <a-upload ref="uploadRef" v-bind="uploadProps" @change="handleChange">
      <template #upload-button>
        <a-button type="primary" style="width: 100px">
          <icon-cloud />
        </a-button>
      </template>
    </a-upload>
  </a-input-group>
</template>
