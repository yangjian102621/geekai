<template>
  <div class="gem-list">
    <div
      v-for="gem in displayList"
      :key="gem.id"
      class="gem-row"
      :class="{ active: gem.id === activeId }"
      @click="$emit('select', gem)"
    >
      <div class="gem-row-icon">
        <el-image v-if="gem.icon" :src="gem.icon" class="gem-icon-img" />
        <span v-else class="gem-icon-placeholder" />
      </div>
      <span class="gem-row-name">{{ gem.name }}</span>
      <!-- 侧栏模式且使用 list 时显示更多（取消固定 / 发起新对话） -->
      <span v-if="mode === 'sidebar' && list != null" class="gem-row-more" @click.stop>
        <el-dropdown trigger="click" @command="(cmd) => onCommand(cmd, gem)">
          <span class="more-trigger">
            <i class="iconfont icon-more-horizontal"></i>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="chat">
                <el-icon><ChatDotRound /></el-icon>
                发起新对话
              </el-dropdown-item>
              <el-dropdown-item command="unpin">
                <i class="iconfont icon-unpin"></i>
                取消固定
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </span>
      <span v-else class="gem-row-pin" @click.stop>
        <i class="iconfont icon-more"></i>
      </span>
    </div>
  </div>
</template>

<script setup>
import { ChatDotRound, Rank } from '@element-plus/icons-vue'
import { httpGet } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { computed, onMounted, ref, watch } from 'vue'

const props = defineProps({
  activeId: {
    type: Number,
    default: 0,
  },
  /** 外部传入列表时使用，不请求接口；侧栏传入 pinned 列表 */
  list: {
    type: Array,
    default: null,
  },
  /** sidebar: 侧栏固定列表；manage: 右侧管理页列表 */
  mode: {
    type: String,
    default: 'sidebar',
  },
})

const emit = defineEmits(['select', 'unpin'])

const gems = ref([])

const displayList = computed(() => {
  if (props.list != null && Array.isArray(props.list)) {
    return props.list.length > 8 ? props.list.slice(0, 8) : props.list
  }
  return gems.value
})

function onCommand(cmd, gem) {
  if (cmd === 'chat') {
    emit('select', gem)
  } else if (cmd === 'unpin') {
    emit('unpin', gem.id)
  }
}

onMounted(async () => {
  if (props.list != null) return
  try {
    const res = await httpGet('/api/app/list/user')
    gems.value = res.data || []
  } catch (e) {
    ElMessage.error('加载智能体列表失败：' + (e.message || ''))
  }
})

watch(
  () => props.list,
  (v) => {
    if (v != null) gems.value = []
  },
  { immediate: true }
)
</script>

<style scoped lang="scss">
.gem-list {
  padding: 4px 0 12px;
}

.gem-row {
  display: flex;
  align-items: center;
  width: 100%;
  padding: 8px 10px;
  margin-bottom: 2px;
  border: none;
  border-radius: 8px;
  background: transparent;
  cursor: pointer;
  transition: background-color 0.2s ease;

  &:hover {
    background: #e5e7eb;
  }

  &.active {
    background: #e0e7ff;

    .gem-row-name {
      color: #3730a3;
      font-weight: 500;
    }
  }
}

.gem-row-icon {
  flex-shrink: 0;
  width: 28px;
  height: 28px;
  margin-right: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.gem-icon-img {
  width: 28px;
  height: 28px;
  border-radius: 6px;
}

.gem-icon-placeholder {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  background: linear-gradient(135deg, #a78bfa 0%, #7c3aed 100%);
  opacity: 0.9;
}

.gem-row-name {
  flex: 1;
  min-width: 0;
  font-size: 14px;
  color: #374151;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.gem-row-pin,
.gem-row-more {
  flex-shrink: 0;
  margin-left: 4px;
  color: #9ca3af;
  cursor: pointer;
  display: inline-flex;

  &:hover {
    color: #6b7280;
  }
}

.more-trigger {
  padding: 2px 4px;
  display: inline-flex;
  align-items: center;
}
</style>
