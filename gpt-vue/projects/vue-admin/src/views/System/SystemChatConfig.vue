<script lang="ts" setup>
import { onMounted } from "vue";
import { Message } from "@arco-design/web-vue";
import useSubmit from "@/composables/useSubmit";
import { getConfig, save } from "./api";

const {
  formRef,
  formData: chat,
  handleSubmit,
  submitting,
} = useSubmit({
  open_ai: { temperature: 1, max_tokens: 1024 },
  azure: { temperature: 1, max_tokens: 1024 },
  chat_gml: { temperature: 0.95, max_tokens: 1024 },
  baidu: { temperature: 0.95, max_tokens: 1024 },
  xun_fei: { temperature: 0.5, max_tokens: 1024 },
  context_deep: 0,
  enable_context: true,
  enable_history: true,
  dall_api_url: "",
});

const rules = {
  init_chat_calls: [{ required: true, message: "请输入赠送对话次数" }],
  user_img_calls: [{ required: true, message: "请输入赠送绘图次数" }],
};

const handleSave = async () => {
  await handleSubmit(
    () =>
      save({
        key: "chat",
        config: chat,
      }),
    {}
  );
  Message.success("保存成功");
};

const reload = async () => {
  const { data } = await getConfig({ key: "chat" });
  data && Object.assign(chat, data);
};

onMounted(reload);
</script>
<template>
  <a-form ref="formRef" :model="chat" :rules="rules" auto-label-width>
    <a-form-item label="开启聊天上下文">
      <a-switch v-model="chat['enable_context']" />
    </a-form-item>
    <a-form-item label="保存聊天记录">
      <a-switch v-model="chat['enable_history']" />
    </a-form-item>
    <a-form-item
      label="会话上下文深度"
      extra="会话上下文深度：在老会话中继续会话，默认加载多少条聊天记录作为上下文。如果设置为 0
        则不加载聊天记录，仅仅使用当前角色的上下文。该配置参数最好设置需要为偶数，否则将无法兼容百度的
        API。"
    >
      <a-input-number v-model="chat['context_deep']" :min="0" :max="10" />
    </a-form-item>

    <a-divider content-position="center">OpenAI</a-divider>
    <a-form-item label="模型创意度" extra="值越大 AI 回答越发散，值越小回答越保守，建议保持默认值">
      <a-slider v-model="chat['open_ai']['temperature']" :max="2" :step="0.1" />
    </a-form-item>
    <a-form-item label="最大响应长度">
      <a-input
        v-model.number="chat['open_ai']['max_tokens']"
        placeholder="回复的最大字数，最大4096"
      />
    </a-form-item>

    <a-divider content-position="center">Azure</a-divider>
    <a-form-item label="模型创意度" extra="值越大 AI 回答越发散，值越小回答越保守，建议保持默认值">
      <a-slider v-model="chat['azure']['temperature']" :max="2" :step="0.1" />
    </a-form-item>
    <a-form-item label="最大响应长度">
      <a-input
        v-model.number="chat['azure']['max_tokens']"
        placeholder="回复的最大字数，最大4096"
      />
    </a-form-item>

    <a-divider content-position="center">ChatGLM</a-divider>
    <a-form-item label="模型创意度" extra="值越大 AI 回答越发散，值越小回答越保守，建议保持默认值">
      <a-slider v-model="chat['chat_gml']['temperature']" :max="1" :step="0.01" />
    </a-form-item>
    <a-form-item label="最大响应长度">
      <a-input
        v-model.number="chat['chat_gml']['max_tokens']"
        placeholder="回复的最大字数，最大4096"
      />
    </a-form-item>

    <a-divider content-position="center">文心一言</a-divider>
    <a-form-item label="模型创意度" extra="值越大 AI 回答越发散，值越小回答越保守，建议保持默认值">
      <a-slider v-model="chat['baidu']['temperature']" :max="1" :step="0.01" />
    </a-form-item>
    <a-form-item label="最大响应长度">
      <a-input
        v-model.number="chat['baidu']['max_tokens']"
        placeholder="回复的最大字数，最大4096"
      />
    </a-form-item>

    <a-divider content-position="center">讯飞星火</a-divider>
    <a-form-item label="模型创意度" extra="值越大 AI 回答越发散，值越小回答越保守，建议保持默认值">
      <a-slider v-model="chat['xun_fei']['temperature']" :max="1" :step="0.1" />
    </a-form-item>
    <a-form-item label="最大响应长度">
      <a-input
        v-model.number="chat['xun_fei']['max_tokens']"
        placeholder="回复的最大字数，最大4096"
      />
    </a-form-item>

    <a-divider content-position="center">AI绘图</a-divider>
    <a-form-item label="DALL-E3出图数量">
      <a-input
        v-model.number="chat['dall_img_num']"
        placeholder="调用 DALL E3 API 传入的出图数量"
      />
    </a-form-item>
    <a-form-item>
      <a-button type="primary" :loading="submitting" @click="handleSave">提交</a-button>
    </a-form-item>
  </a-form>
</template>
