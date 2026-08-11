<template>
  <div class="chat-input-container">
    <div class="chat-input-wrapper">
      <div ref="textHeightRef" class="text-height-calculator">{{ localPrompt }}</div>
      <div class="input-content">
        <div class="input-area">
          <div v-if="files.length > 0" class="file-preview-list">
            <file-list :files="files" @remove-file="handleRemoveFile" />
          </div>
          <textarea
            ref="inputRef"
            class="prompt-textarea"
            :rows="row"
            v-model="localPrompt"
            @keydown="handleInput"
            @input="handleInput"
            @paste="handlePaste"
            :placeholder="placeholder"
            autofocus
          ></textarea>
        </div>
        <div class="toolbar-container">
          <div class="toolbar-left">
            <span class="toolbar-btn" title="上传附件">
              <file-select :user-id="userId" @selected="handleFileSelected" />
            </span>
            <el-dropdown :hide-on-click="false" trigger="click">
              <span class="toolbar-btn" title="插件">
                <i class="iconfont icon-plugin"></i>
              </span>
              <template #dropdown>
                <el-dropdown-menu class="tools-dropdown">
                  <el-checkbox-group v-model="localToolSelected">
                    <el-dropdown-item v-for="item in tools" :key="item.id">
                      <el-checkbox :value="item.id" :label="item.label" />
                      <el-tooltip :content="item.description" placement="right">
                        <el-icon><InfoFilled /></el-icon>
                      </el-tooltip>
                    </el-dropdown-item>
                  </el-checkbox-group>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <span class="toolbar-btn" title="配置" @click="$emit('show-setting')">
              <i class="iconfont icon-config"></i>
            </span>
            <span v-if="agentName" class="agent-label">
              {{ agentName }}
            </span>
          </div>
          <div class="toolbar-right">
            <span class="model-selector-wrapper">
              <el-select
                v-model="localModelId"
                size="small"
                class="model-select"
                :disabled="disableModel"
                placeholder="选择模型"
                filterable
              >
                <template #label="{ label }">
                  <span class="text-sm font-bold">{{ selectedModel?.name || label }}</span>
                </template>
                <el-option v-for="m in models" :key="m.id" :value="m.id">
                  <div class="flex items-center justify-center gap-1.5 h-full">
                    <span class="text-sm font-bold">{{ m.name }}</span>
                    <el-tag size="small" :type="m.power > 0 ? 'info' : 'success'" effect="plain">{{
                      m.power > 0 ? m.power + '算力/次' : '免费'
                    }}</el-tag>
                  </div>
                </el-option>
              </el-select>
            </span>
            <span class="send-btn-wrapper">
              <el-button
                type="info"
                v-if="isGenerating"
                @click="$emit('stop')"
                plain
                class="action-btn stop-btn"
              >
                <el-icon><VideoPause /></el-icon>
              </el-button>
              <el-button
                @click="$emit('send')"
                :disabled="isGenerating"
                v-else
                title="发送"
                class="action-btn send-btn"
              >
                <el-icon><Promotion /></el-icon>
              </el-button>
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import FileList from '@/components/FileList.vue'
import FileSelect from '@/components/FileSelect.vue'
import { InfoFilled, Promotion, VideoPause } from '@element-plus/icons-vue'
import { computed, ref } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  files: {
    type: Array,
    default: () => [],
  },
  modelId: {
    type: Number,
    required: true,
  },
  models: {
    type: Array,
    default: () => [],
  },
  tools: {
    type: Array,
    default: () => [],
  },
  toolSelected: {
    type: Array,
    default: () => [],
  },
  isGenerating: {
    type: Boolean,
    default: false,
  },
  disableModel: {
    type: Boolean,
    default: false,
  },
  placeholder: {
    type: String,
    default: '按 Enter 键发送消息，使用 Shift + Enter 换行',
  },
  userId: {
    type: Number,
    default: 0,
  },
  agentName: {
    type: String,
    default: '',
  },
})

const emit = defineEmits([
  'update:modelValue',
  'update:modelId',
  'update:toolSelected',
  'send',
  'stop',
  'paste',
  'input',
  'file-selected',
  'file-removed',
  'show-setting',
])

const inputRef = ref(null)
const textHeightRef = ref(null)
const row = ref(3)

const localPrompt = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
})

const localModelId = computed({
  get: () => props.modelId,
  set: (value) => emit('update:modelId', value),
})

const localToolSelected = computed({
  get: () => props.toolSelected,
  set: (value) => emit('update:toolSelected', value),
})

const selectedModel = computed(() => {
  return props.models.find((model) => model.id === props.modelId)
})

const handleInput = (e) => {
  if (!textHeightRef.value || !inputRef.value) return

  const lineHeight = parseFloat(window.getComputedStyle(inputRef.value).lineHeight)
  textHeightRef.value.style.width = inputRef.value.clientWidth + 'px'
  const lines = Math.floor(textHeightRef.value.clientHeight / lineHeight)

  if (localPrompt.value.length < 10) {
    row.value = 1
  } else if (lines <= 7) {
    row.value = lines
  } else {
    row.value = 7
  }

  emit('input', e)

  // 输入回车自动提交
  if (e.keyCode === 13) {
    // Shift + Enter 换行
    if (e.shiftKey) {
      return
    }
    e.preventDefault()
    emit('send')
  }
}

const handlePaste = (e) => {
  emit('paste', e)
}

const handleFileSelected = (file) => {
  emit('file-selected', file)
}

const handleRemoveFile = (file) => {
  emit('file-removed', file)
}
</script>

<style scoped lang="scss">
.chat-input-container {
  width: 100%;
  padding: 0 16px;
}

.chat-input-wrapper {
  width: 100%;
  background: var(--el-bg-color, #ffffff);
  border: 2px solid var(--el-border-color-light, #e4e7ed);
  border-radius: 12px;
  padding: 12px 16px;
  transition: all 0.2s ease;

  &:hover {
    border-color: var(--el-color-primary-light-5, #a0cfff);
  }

  &:focus-within {
    border-color: var(--el-color-primary, #409eff);
    box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.1);
  }
}

.text-height-calculator {
  position: absolute;
  visibility: hidden;
  white-space: pre-wrap;
  word-wrap: break-word;
  overflow-wrap: break-word;
  line-height: 24px;
  font-size: 14px;
  pointer-events: none;
}

.input-content {
  width: 100%;
}

.input-area {
  display: flex;
  flex-direction: column;
  width: 100%;
  margin-bottom: 8px;
}

.file-preview-list {
  padding-bottom: 12px;
  border-bottom: 1px solid var(--el-border-color-lighter, #ebeef5);
  margin-bottom: 12px;
}

.prompt-textarea {
  width: 100%;
  min-height: 56px;
  line-height: 24px;
  font-size: 14px;
  color: var(--el-text-color-primary, #303133);
  background: transparent;
  border: none;
  outline: none;
  resize: none;
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: inherit;

  &::placeholder {
    color: var(--el-text-color-placeholder, #a8abb2);
  }

  &::-webkit-scrollbar {
    width: 0;
    height: 0;
  }
}

.toolbar-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  padding-top: 8px;
  border-top: 1px solid var(--el-border-color-lighter, #ebeef5);
}

.toolbar-left,
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.agent-label {
  margin-left: 4px;
  padding: 2px 8px;
  border-radius: 999px;
  border: 1px solid var(--el-color-primary, #409eff);
  font-size: 12px;
  color: var(--el-color-primary, #409eff);
  max-width: 140px;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}

.toolbar-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: var(--el-text-color-regular, #606266);

  &:hover {
    background: var(--el-fill-color-light, #f5f7fa);
    color: var(--el-color-primary, #409eff);
  }

  .iconfont {
    font-size: 18px;
  }
}

.model-selector-wrapper {
  display: inline-flex;
}

.model-select {
  min-width: 180px;

  :deep(.el-input__wrapper) {
    border-radius: 6px;
    transition: all 0.2s ease;
  }
}

.model-label-tags {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.send-btn-wrapper {
  display: inline-flex;
}

.action-btn {
  width: 36px;
  height: 36px;
  padding: 0;
  border-radius: 50%;
  border: none;
  transition: all 0.2s ease;

  &.send-btn {
    background: var(--el-color-primary, #409eff);
    color: #ffffff;

    &:hover:not(:disabled) {
      background: var(--el-color-primary-dark-2, #337ecc);
      transform: scale(1.05);
    }

    &:active:not(:disabled) {
      transform: scale(0.95);
    }

    &:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  }

  &.stop-btn {
    background: var(--el-color-info-light-3, #b1b3b8);
    color: #ffffff;

    &:hover {
      background: var(--el-color-info, #909399);
    }
  }

  :deep(.el-icon) {
    font-size: 18px;
  }
}

// Dark mode support
@media (prefers-color-scheme: dark) {
  .chat-input-wrapper {
    background: var(--el-bg-color, #1a1a1a);
    border-color: var(--el-border-color-light, #414243);

    &:hover {
      border-color: var(--el-color-primary-light-3, #79bbff);
    }

    &:focus-within {
      border-color: var(--el-color-primary, #409eff);
      box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.2);
    }
  }

  .prompt-textarea {
    color: var(--el-text-color-primary, #e5eaf3);

    &::placeholder {
      color: var(--el-text-color-placeholder, #8d9095);
    }
  }

  .toolbar-btn {
    color: var(--el-text-color-regular, #cfd3dc);

    &:hover {
      background: var(--el-fill-color-light, #262727);
      color: var(--el-color-primary, #409eff);
    }
  }
}

// Responsive design
@media (max-width: 768px) {
  .chat-input-container {
    padding: 0 12px;
  }

  .chat-input-wrapper {
    padding: 10px 12px;
  }

  .toolbar-container {
    gap: 8px;
  }

  .toolbar-left,
  .toolbar-right {
    gap: 6px;
  }

  .toolbar-btn {
    width: 28px;
    height: 28px;

    .iconfont {
      font-size: 16px;
    }
  }

  .model-select {
    min-width: 140px;
  }

  .action-btn {
    width: 32px;
    height: 32px;

    :deep(.el-icon) {
      font-size: 16px;
    }
  }
}
</style>
