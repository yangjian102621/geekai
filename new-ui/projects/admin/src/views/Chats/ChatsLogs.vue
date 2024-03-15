<script lang="ts" setup>
import { onMounted } from "vue";
import { Message } from "@arco-design/web-vue";
import { dateFormat } from "@chatgpt-plus/packages/utils";
import useRequest from "@/composables/useRequest";
import { history } from "./api";

const props = defineProps({
  id: String,
});

const [getData, data, loading] = useRequest(history);
onMounted(async () => {
  await getData({ chat_id: props.id });
});
</script>
<template>
  <template v-if="loading">
    <div class="custom-skeleton">
      <a-skeleton-shape />
      <div style="flex: 1">
        <a-skeleton-line :rows="2" />
      </div>
    </div>
  </template>
  <template v-else>
    <div v-for="item in data" :key="item.id">
      <div class="item-container" :class="item.type">
        <div class="left">
          <a-avatar shape="square">
            <img :src="item.icon" />
          </a-avatar>
        </div>
        <a-space class="right" direction="vertical">
          <div>{{ item.content }}</div>
          <a-space>
            <div class="code">
              <icon-clock-circle />
              {{ dateFormat(item.created_at) }}
            </div>
            <div class="code">算力消耗: {{ item.tokens }}</div>
            <a-typography-text
              v-if="item.type === 'reply'"
              copyable
              :copy-delay="1000"
              :copy-text="item.content"
              @copy="Message.success('复制成功')"
            >
              <template #copy-icon>
                <a-button class="code" size="mini">
                  <icon-copy />
                </a-button>
              </template>
              <template #copy-tooltip>复制回答</template>
            </a-typography-text>
          </a-space>
        </a-space>
      </div>
    </div>
  </template>
</template>
<style lang="less" scoped>
.item-container {
  display: flex;
  padding: 20px 10px;
  width: 100%;
  gap: 20px;
  border-bottom: 1px solid #d9d9e3;
  box-sizing: border-box;
  align-items: flex-start;
  &.reply {
    background: #f7f7f8;
  }
  .left {
    width: 40px;
  }
  .right {
    flex: 1;
    overflow: hidden;
  }
  .code {
    background-color: #e7e7e8;
    color: #888;
    padding: 3px 5px;
    margin-right: 10px;
    border-radius: 5px;
  }
}
.custom-skeleton {
  display: flex;
  width: 100%;
  gap: 12px;
}
</style>
