<template>
  <a-alert type="warning">
    <div class="warning">
      <div>注意：如果是百度文心一言平台，API-KEY 为 APIKey|SecretKey，中间用竖线（|）连接</div>
      <div>注意：如果是讯飞星火大模型，API-KEY 为 AppId|APIKey|APISecret，中间用竖线（|）连接</div>
    </div>
  </a-alert>
  <a-form
    ref="formRef"
    :model="form"
    :style="{ width: '600px', 'margin-top': '10px' }"
    @submit="handleSubmit"
  >
    <a-form-item
      field="platform"
      label="所属平台"
      :rules="[{ required: true, message: '请输入所属平台' }]"
      :validate-trigger="['change', 'input']"
    >
      <a-select
        v-model="form.platform"
        placeholder="请输入所属平台"
        :options="platformOptions"
        @change="handlePlatformChange"
      />
    </a-form-item>
    <a-form-item
      field="name"
      label="名称"
      :rules="[{ required: true, message: '请输入名称' }]"
      :validate-trigger="['change', 'input']"
      showable
    >
      <a-input v-model="form.name" placeholder="请输入名称" />
    </a-form-item>
    <a-form-item
      field="type"
      label="用途"
      :rules="[{ required: true, message: '请输入用途' }]"
      :validate-trigger="['change', 'input']"
    >
      <a-select
        v-model="form.type"
        placeholder="请输入用途"
        :options="typeOptions"
        @change="handlePlatformChange"
      />
    </a-form-item>
    <a-form-item
      field="value"
      label="API KEY"
      :rules="[{ required: true, message: '请输入API KEY' }]"
      :validate-trigger="['change', 'input']"
    >
      <a-input v-model="form.value" placeholder="请输入API KEY" />
    </a-form-item>
    <a-form-item
      field="api_url"
      label="API URL"
      :rules="[{ required: true, message: '请输入API URL' }]"
      :validate-trigger="['change', 'input']"
    >
      <a-input v-model="form.api_url" placeholder="请输入API URL" />
    </a-form-item>

    <a-form-item field="use_proxy" label="使用代理">
      <a-space>
        <a-switch v-model="form.use_proxy" />
        <a-tooltip
          content="是否使用代理访问 API URL，OpenAI 官方API需要开启代理访问"
          position="right"
        >
          <icon-info-circle-fill />
        </a-tooltip>
      </a-space>
    </a-form-item>
    <a-form-item field="enable" label="启用状态">
      <a-switch v-model="form.enable" />
    </a-form-item>
  </a-form>
</template>

<script setup>
import { ref, defineExpose, defineProps } from "vue";
const props = defineProps({
  data: {},
});

const formRef = ref();
const form = ref({});
if (props.data?.id) {
  form.value = Object.assign({}, props.data);
}

defineExpose({
  formRef,
  form,
});

const typeOptions = [
  {
    label: "聊天",
    value: "chart",
  },
  {
    label: "绘图",
    value: "img",
  },
];

const platformOptions = [
  {
    label: "【OpenAI】ChatGPT",
    value: "OpenAI",
    api_url: "https://gpt.bemore.lol/v1/chat/completions",
    img_url: "https://gpt.bemore.lol/v1/images/generations",
  },
  {
    label: "【讯飞】星火大模型",
    value: "XunFei",
    api_url: "wss://spark-api.xf-yun.com/{version}/chat",
  },
  {
    label: "【清华智普】ChatGLM",
    value: "ChatGLM",
    api_url: "https://open.bigmodel.cn/api/paas/v3/model-api/{model}/sse-invoke",
  },
  {
    label: "【百度】文心一言",
    value: "Baidu",
    api_url: "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/{model}",
  },
  {
    label: "【微软】Azure",
    value: "Azure",
    api_url:
      "https://chat-bot-api.openai.azure.com/openai/deployments/{model}/chat/completions?api-version=2023-05-15",
  },
  {
    label: "【阿里】千义通问",
    value: "QWen",
    api_url: "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation",
  },
];

const handlePlatformChange = () => {
  const obj = platformOptions.find((item) => item.value === form.value.platform);
  if (obj) {
    form.value.api_url = form.value.type === "img" ? obj.img_url : obj.api_url;
  }
};
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
.warning {
  color: #e6a23c;
  white-space: pre;
}
</style>
