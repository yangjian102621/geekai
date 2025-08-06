<template>
  <van-cell class="app-cell">
    <div class="app-card">
      <div class="app-info">
        <div class="app-image">
          <van-image :src="app.icon" round />
        </div>
        <div class="app-detail">
          <div class="app-title">{{ app.name }}</div>
          <div class="app-desc">{{ app.hello_msg }}</div>
        </div>
      </div>

      <div class="app-actions">
        <van-button
          size="small"
          type="primary"
          class="action-btn"
          @click="$emit('use-role', app.id)"
          >开始对话</van-button
        >
        <van-button
          size="small"
          :type="hasRole ? 'danger' : 'success'"
          class="action-btn"
          @click="$emit('update-role', app, hasRole ? 'remove' : 'add')"
        >
          {{ hasRole ? '移出工作台' : '添加到工作台' }}
        </van-button>
      </div>
    </div>
  </van-cell>
</template>

<script setup>
defineProps({
  app: {
    type: Object,
    required: true,
  },
  hasRole: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['use-role', 'update-role'])
</script>

<style scoped lang="scss">
.app-cell {
  padding: 0;
  margin-bottom: 15px;

  .app-card {
    background: var(--van-cell-background);
    border-radius: 12px;
    padding: 15px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);

    .app-info {
      display: flex;
      align-items: center;
      margin-bottom: 15px;

      .app-image {
        width: 60px;
        height: 60px;
        margin-right: 15px;

        :deep(.van-image) {
          width: 100%;
          height: 100%;
        }
      }

      .app-detail {
        flex: 1;

        .app-title {
          font-size: 16px;
          font-weight: 600;
          margin-bottom: 5px;
          color: var(--van-text-color);
        }

        .app-desc {
          font-size: 13px;
          color: var(--van-gray-6);
          display: -webkit-box;
          -webkit-box-orient: vertical;
          -webkit-line-clamp: 2;
          overflow: hidden;
        }
      }
    }

    .app-actions {
      display: flex;
      gap: 10px;

      .action-btn {
        flex: 1;
        border-radius: 20px;
        padding: 0 10px;
      }
    }
  }
}
</style>
