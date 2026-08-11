<template>
  <div class="agent-manage-page">
    <!-- 系统预设的智能体：默认一行，可展开 -->
    <section class="agent-section">
      <div class="agent-section-head">
        <h3 class="agent-section-title">系统预设的智能体</h3>
        <el-button
          text
          type="primary"
          size="small"
          class="expand-btn"
          @click="systemExpanded = !systemExpanded"
        >
          {{ systemExpanded ? '收起' : '展开' }}
          <el-icon class="expand-icon" :class="{ expanded: systemExpanded }">
            <ArrowDown />
          </el-icon>
        </el-button>
      </div>
      <div class="system-cards" :class="{ 'one-row': !systemExpanded }">
        <button
          v-for="item in systemAgentsDisplay"
          :key="item.id"
          type="button"
          class="system-card"
          @click="$emit('select', item)"
        >
          <div class="system-card-icon-wrap">
            <el-image v-if="item.icon" :src="item.icon" class="system-card-icon" fit="cover" />
            <span v-else class="system-card-icon-ph" />
          </div>
          <div class="system-card-body pl-2">
            <span class="system-card-name">{{ item.name }}</span>
            <p class="system-card-desc">{{ item.hello_msg || '' }}</p>
          </div>
          <span class="system-card-more" @click.stop>
            <el-dropdown trigger="click" @command="(cmd) => onSystemCommand(cmd, item)">
              <span class="more-trigger">
                <i class="iconfont icon-more-horizontal"></i>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="chat">
                    <el-icon><ChatDotRound /></el-icon>
                    发起新对话
                  </el-dropdown-item>
                  <el-dropdown-item command="pin">
                    <template v-if="isPinned(item.id)">
                      <i class="iconfont icon-unpin"></i>
                      {{ '取消固定' }}
                    </template>
                    <template v-else>
                      <i class="iconfont icon-pin"></i>
                      {{ '固定' }}
                    </template>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </span>
        </button>
      </div>
    </section>

    <!-- 我的智能体 -->
    <section class="agent-section">
      <div class="agent-section-head">
        <h3 class="agent-section-title">
          我的智能体
          <el-icon class="title-info"><InfoFilled /></el-icon>
        </h3>
        <el-button type="primary" size="small" @click="$emit('new')">+ 新建智能体</el-button>
      </div>
      <div class="my-list">
        <div
          v-for="item in myAgents"
          :key="item.id"
          class="my-list-row"
          @click="$emit('select', item)"
        >
          <div class="my-list-icon">
            <el-image v-if="item.icon" :src="item.icon" class="my-list-icon-img" fit="cover" />
            <span v-else class="my-list-icon-ph" />
          </div>
          <div class="my-list-body">
            <span class="my-list-name">{{ item.name }}</span>
            <p class="my-list-desc">{{ item.hello_msg || '' }}</p>
          </div>
          <div class="my-list-actions" @click.stop>
            <span class="my-list-action" title="编辑" @click="emit('edit', item)">
              <el-icon><Edit /></el-icon>
            </span>
            <el-dropdown trigger="click" @command="(cmd) => onMyCommand(cmd, item)">
              <span class="my-list-action more-trigger">
                <i class="iconfont icon-more-horizontal"></i>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="chat">
                    <el-icon><ChatDotRound /></el-icon>
                    发起新对话
                  </el-dropdown-item>
                  <el-dropdown-item command="pin">
                    <template v-if="isPinned(item.id)">
                      <i class="iconfont icon-unpin"></i>
                      {{ '取消固定' }}
                    </template>
                    <template v-else>
                      <i class="iconfont icon-pin"></i>
                      {{ '固定' }}
                    </template>
                  </el-dropdown-item>
                  <el-dropdown-item command="edit">
                    <el-icon><Edit /></el-icon>
                    编辑
                  </el-dropdown-item>
                  <el-dropdown-item command="copy">
                    <el-icon><DocumentCopy /></el-icon>
                    复制
                  </el-dropdown-item>
                  <el-dropdown-item command="delete" divided>
                    <el-icon><Delete /></el-icon>
                    <span style="color: var(--el-color-danger)">删除</span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
        <div v-if="myAgents.length === 0" class="my-list-empty">
          暂无自己创建的智能体，点击「新建智能体」创建
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import {
  ArrowDown,
  ChatDotRound,
  Delete,
  DocumentCopy,
  Edit,
  InfoFilled,
  Rank,
  Share,
} from '@element-plus/icons-vue'
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { computed, onMounted, ref, watch } from 'vue'

const props = defineProps({
  profile: {
    type: Object,
    default: () => ({ gem_ids: [] }),
  },
  refreshTrigger: {
    type: Number,
    default: 0,
  },
})

const emit = defineEmits(['select', 'new', 'edit', 'profile-update'])

const allAgents = ref([])
const systemExpanded = ref(false)

const gemIds = computed(() => props.profile?.gem_ids || [])

// 统一按数字比较，避免 API 返回字符串导致 pin/unpin 失效
const toNum = (v) => (typeof v === 'number' && !Number.isNaN(v) ? v : Number(v))
const normalizedGemIds = computed(() =>
  (gemIds.value || []).map(toNum).filter((n) => !Number.isNaN(n) && n > 0)
)

const systemAgents = computed(() => allAgents.value.filter((a) => a.user_id === 0))
const systemAgentsDisplay = computed(() =>
  systemExpanded.value ? systemAgents.value : systemAgents.value.slice(0, 3)
)
const myAgents = computed(() => allAgents.value.filter((a) => a.user_id > 0))

function isPinned(id) {
  const n = toNum(id)
  return normalizedGemIds.value.some((x) => x === n)
}

async function updateGemIds(newIds) {
  const raw = Array.isArray(newIds) ? newIds.slice(0, 8) : []
  const ids = raw.map(toNum).filter((n) => !Number.isNaN(n) && n > 0)
  try {
    await httpPost('/api/user/profile/update', { ...props.profile, gem_ids: ids })
    emit('profile-update', ids)
  } catch (e) {
    ElMessage.error('操作失败：' + (e.message || ''))
  }
}

function doPin(id) {
  const n = toNum(id)
  if (normalizedGemIds.value.some((x) => x === n)) return
  if (normalizedGemIds.value.length >= 8) {
    ElMessage.warning('最多固定 8 个智能体')
    return
  }
  const ids = [...normalizedGemIds.value, n]
  updateGemIds(ids)
}

function doUnpin(id) {
  const target = toNum(id)
  const ids = normalizedGemIds.value.filter((x) => x !== target)
  updateGemIds(ids)
}

async function fetchAgentList() {
  try {
    const res = await httpGet('/api/app/list/user')
    allAgents.value = res.data || []
  } catch (e) {
    ElMessage.error('加载智能体列表失败：' + (e.message || ''))
  }
}

onMounted(fetchAgentList)
watch(() => props.refreshTrigger, fetchAgentList)

function onSystemCommand(cmd, item) {
  if (cmd === 'chat') emit('select', item)
  else if (cmd === 'pin') {
    if (isPinned(item.id)) doUnpin(item.id)
    else doPin(item.id)
  }
}

async function onMyCommand(cmd, item) {
  if (cmd === 'chat') {
    emit('select', item)
    return
  }
  if (cmd === 'edit') {
    emit('edit', item)
    return
  }
  if (cmd === 'pin') {
    if (isPinned(item.id)) doUnpin(item.id)
    else doPin(item.id)
    return
  }
  if (cmd === 'copy') {
    try {
      const res = await httpPost('/api/app/copy', { source_id: item.id })
      ElMessage.success('已复制到「我的智能体」')
      const listRes = await httpGet('/api/app/list/user')
      allAgents.value = listRes.data || []
      if (res.data?.id) doPin(res.data.id)
    } catch (e) {
      ElMessage.error('复制失败：' + (e.message || ''))
    }
    return
  }
  if (cmd === 'delete') {
    try {
      await httpPost('/api/app/remove', { id: item.id })
      ElMessage.success('已删除')
      const listRes = await httpGet('/api/app/list/user')
      allAgents.value = listRes.data || []
      doUnpin(item.id)
    } catch (e) {
      ElMessage.error('删除失败：' + (e.message || ''))
    }
  }
}
</script>

<style scoped lang="scss">
.agent-manage-page {
  padding: 0 0 32px;
}

.agent-section {
  margin-bottom: 28px;
}

.agent-section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.agent-section-title {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-fb, #111827);
  display: inline-flex;
  align-items: center;
  gap: 6px;

  .title-info {
    font-size: 14px;
    color: #9ca3af;
    cursor: help;
  }
}

.expand-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
}

.expand-icon {
  transition: transform 0.2s ease;
  &.expanded {
    transform: rotate(180deg);
  }
}

.system-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;

  &.one-row {
    grid-template-columns: repeat(3, 1fr);
    overflow: hidden;
  }
}

.system-card {
  display: flex;
  align-items: flex-start;
  position: relative;
  width: 100%;
  padding: 14px;
  border: none;
  border-radius: 12px;
  background: #f9fafb;
  cursor: pointer;
  transition:
    background-color 0.2s ease,
    box-shadow 0.2s ease;
  text-align: left;

  &:hover {
    background: #f3f4f6;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  }
}

.system-card-icon-wrap {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  overflow: hidden;
  background: #e5e7eb;
}

.system-card-icon {
  width: 40px;
  height: 40px;
  display: block;
}

.system-card-icon-ph {
  display: block;
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #a78bfa 0%, #7c3aed 100%);
  opacity: 0.9;
}

.system-card-body {
  flex: 1;
  min-width: 0;
  padding-right: 28px;
}

.system-card-name {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--text-fb, #111827);
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.system-card-desc {
  margin: 4px 0 0;
  font-size: 0.75rem;
  color: #6b7280;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.system-card-pin,
.system-card-more {
  position: absolute;
  top: 10px;
  right: 8px;
  color: #9ca3af;
  cursor: pointer;
  padding: 4px;

  &:hover {
    color: #374151;
  }
}

.system-card-more {
  right: 8px;
}

.system-card-pin {
  right: 36px;
}

.more-trigger {
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

/* 我的智能体列表 */
.my-list {
  background: #fff;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  overflow: hidden;
}

.my-list-row {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  transition: background-color 0.2s ease;
  border-bottom: 1px solid #f3f4f6;

  &:last-child {
    border-bottom: none;
  }

  &:hover {
    background: #f9fafb;
  }
}

.my-list-icon {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  margin-right: 12px;
  border-radius: 10px;
  overflow: hidden;
  background: #e5e7eb;
}

.my-list-icon-img {
  width: 40px;
  height: 40px;
  display: block;
}

.my-list-icon-ph {
  display: block;
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #a78bfa 0%, #7c3aed 100%);
  opacity: 0.9;
}

.my-list-body {
  flex: 1;
  min-width: 0;
}

.my-list-name {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--text-fb, #111827);
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.my-list-desc {
  margin: 4px 0 0;
  font-size: 0.8125rem;
  color: #6b7280;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.my-list-actions {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 4px;
  margin-left: 8px;
}

.my-list-action {
  color: #9ca3af;
  padding: 6px;
  cursor: pointer;
  display: inline-flex;
  transition: color 0.2s ease;

  &:hover {
    color: var(--el-color-primary);
  }
}

.my-list-empty {
  padding: 24px 16px;
  text-align: center;
  color: #9ca3af;
  font-size: 0.875rem;
}

@media (max-width: 768px) {
  .system-cards.one-row {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
