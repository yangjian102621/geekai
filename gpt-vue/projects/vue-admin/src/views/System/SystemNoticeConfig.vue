<script lang="ts" setup>
import { onMounted } from "vue";
import { Message } from "@arco-design/web-vue";
import MdEditor from "md-editor-v3";
import "md-editor-v3/lib/style.css";
import http from "@/http/config";
import useSubmit from "@/composables/useSubmit";
import { getConfig, save } from "./api";

const { formRef, formData, handleSubmit, submitting } = useSubmit({
  content: "",
  updated: true,
});

const handleSave = async () => {
  await handleSubmit(
    () =>
      save({
        key: "notice",
        config: formData,
      }),
    {}
  );
  Message.success("保存成功");
};

const reload = async () => {
  const { data } = await getConfig({ key: "notice" });
  data && Object.assign(formData, data);
};

const onUploadImg = (files, callback) => {
  Promise.all(
    files.map((file) => {
      return new Promise((rev, rej) => {
        const formData = new FormData();
        formData.append("file", file, file.name);
        http({
          url: `/api/upload`,
          data: formData,
        })
          .then((res) => rev(res))
          .catch((e) => rej(e));
      });
    })
  )
    .then((res) => {
      Message.success({ content: "上传成功", duration: 500 });
      callback(res.map((item) => item.data.url));
    })
    .catch((e) => {
      Message.error("图片上传失败:" + e.message);
    });
};

onMounted(reload);
</script>
<template>
  <a-form ref="formRef" :model="formData" auto-label-width :disabled="submitting">
    <a-form-item>
      <md-editor v-model="formData.content" @on-upload-img="onUploadImg" />
    </a-form-item>
    <a-form-item>
      <a-button type="primary" :loading="submitting" @click="handleSave">提交</a-button>
    </a-form-item>
  </a-form>
</template>
