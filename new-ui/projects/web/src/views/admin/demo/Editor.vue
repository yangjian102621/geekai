<template>
  <div class="container">
    <div class="plugins-tips">
      md-editor-v3：vue3版本的 markdown 编辑器，配置丰富，请详看文档。 访问地址：
      <a href="https://imzbf.github.io/md-editor-v3/index" target="_blank">md-editor-v3</a>
    </div>
    <md-editor class="mgb20" v-model="text" @on-upload-img="onUploadImg"/>
    <el-button type="primary">提交</el-button>
  </div>
</template>

<script setup>
import MdEditor from 'md-editor-v3';
import 'md-editor-v3/lib/style.css';
import {ref} from "vue";

const text = ref(`
ChatGPT-PLUS 是基于 OpenAI API 实现的 ChatGPT Web 系统。主要有如下特性：
* 完整的开源系统，前端应用和后台管理系统皆可开箱即用。
* 聊天体验跟 ChatGPT 官方版本完全一致。
* 内置了各种预训练好的角色，比如小红书写手，英语翻译大师，苏格拉底，孔子，乔布斯，周报助手等。轻松满足你的各种聊天和应用需求。
`);
const onUploadImg = (files) => {
  console.log(files);
};
</script>
