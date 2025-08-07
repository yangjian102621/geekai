<template>
  <div class="empty-state">
    <van-empty :image="getImage()" :description="description" :image-size="imageSize">
      <template #bottom>
        <slot name="action">
          <van-button
            v-if="showAction"
            round
            type="primary"
            class="action-btn"
            @click="$emit('action')"
          >
            {{ actionText }}
          </van-button>
        </slot>
      </template>
    </van-empty>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  type: {
    type: String,
    default: 'search', // search, error, network, default
    validator: (value) => ['search', 'error', 'network', 'default'].includes(value),
  },
  description: {
    type: String,
    default: '暂无数据',
  },
  imageSize: {
    type: [String, Number],
    default: 120,
  },
  showAction: {
    type: Boolean,
    default: false,
  },
  actionText: {
    type: String,
    default: '刷新',
  },
})

defineEmits(['action'])

// 根据类型获取对应的图标
const getImage = () => {
  const imageMap = {
    search: 'search',
    error: 'error',
    network: 'network',
    default: 'default',
  }
  return imageMap[props.type] || 'search'
}
</script>

<style scoped lang="scss">
.empty-state {
  padding: 40px 20px;
  text-align: center;

  .action-btn {
    margin-top: 16px;
    min-width: 120px;
    height: 36px;
    font-size: 14px;
  }
}
</style>
