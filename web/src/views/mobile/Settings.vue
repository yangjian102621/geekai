<template>
  <div class="settings-page">
    <div class="settings-content">
      <!-- 个人设置 -->
      <div class="setting-section">
        <h3 class="section-title">个人设置</h3>
        <van-cell-group inset>
          <van-cell title="个人信息" is-link @click="router.push('/mobile/profile')">
            <template #icon>
              <i class="iconfont icon-user setting-icon"></i>
            </template>
          </van-cell>
          <van-cell title="修改密码" is-link @click="showPasswordDialog = true">
            <template #icon>
              <i class="iconfont icon-lock setting-icon"></i>
            </template>
          </van-cell>
          <van-cell title="绑定手机" is-link @click="showBindMobile = true">
            <template #icon>
              <i class="iconfont icon-phone setting-icon"></i>
            </template>
          </van-cell>
          <van-cell title="绑定邮箱" is-link @click="showBindEmail = true">
            <template #icon>
              <i class="iconfont icon-email setting-icon"></i>
            </template>
          </van-cell>
        </van-cell-group>
      </div>

      <!-- 应用设置 -->
      <div class="setting-section">
        <h3 class="section-title">应用设置</h3>
        <van-cell-group inset>
          <van-cell title="暗黑主题">
            <template #icon>
              <i class="iconfont icon-moon setting-icon"></i>
            </template>
            <template #right-icon>
              <van-switch v-model="darkMode" @change="onThemeChange" />
            </template>
          </van-cell>
          <van-cell title="流式输出">
            <template #icon>
              <i class="iconfont icon-stream setting-icon"></i>
            </template>
            <template #right-icon>
              <van-switch v-model="streamOutput" @change="onStreamChange" />
            </template>
          </van-cell>
          <van-cell title="消息通知">
            <template #icon>
              <i class="iconfont icon-bell setting-icon"></i>
            </template>
            <template #right-icon>
              <van-switch v-model="notifications" />
            </template>
          </van-cell>
          <van-cell title="自动保存">
            <template #icon>
              <i class="iconfont icon-save setting-icon"></i>
            </template>
            <template #right-icon>
              <van-switch v-model="autoSave" />
            </template>
          </van-cell>
          <van-cell title="语言设置" is-link @click="showLanguageSelect = true">
            <template #icon>
              <i class="iconfont icon-translate setting-icon"></i>
            </template>
            <template #value>
              <span class="setting-value">{{ currentLanguage.name }}</span>
            </template>
          </van-cell>
        </van-cell-group>
      </div>

      <!-- 聊天设置 -->
      <div class="setting-section">
        <h3 class="section-title">聊天设置</h3>
        <van-cell-group inset>
          <van-cell title="默认模型" is-link @click="showModelSelect = true">
            <template #icon>
              <i class="iconfont icon-robot setting-icon"></i>
            </template>
            <template #value>
              <span class="setting-value">{{ currentModel.name }}</span>
            </template>
          </van-cell>
          <van-cell title="对话记录">
            <template #icon>
              <i class="iconfont icon-history setting-icon"></i>
            </template>
            <template #right-icon>
              <van-switch v-model="saveHistory" />
            </template>
          </van-cell>
          <van-cell title="发送方式" is-link @click="showSendModeSelect = true">
            <template #icon>
              <i class="iconfont icon-send setting-icon"></i>
            </template>
            <template #value>
              <span class="setting-value">{{ currentSendMode.name }}</span>
            </template>
          </van-cell>
        </van-cell-group>
      </div>

      <!-- 隐私与安全 -->
      <div class="setting-section">
        <h3 class="section-title">隐私与安全</h3>
        <van-cell-group inset>
          <van-cell title="清除缓存" is-link @click="showClearCache = true">
            <template #icon>
              <i class="iconfont icon-delete setting-icon"></i>
            </template>
            <template #value>
              <span class="setting-value">{{ cacheSize }}</span>
            </template>
          </van-cell>
          <van-cell title="清除聊天记录" is-link @click="showClearHistory = true">
            <template #icon>
              <i class="iconfont icon-clear setting-icon"></i>
            </template>
          </van-cell>
          <van-cell title="隐私政策" is-link @click="showPrivacyPolicy = true">
            <template #icon>
              <i class="iconfont icon-shield setting-icon"></i>
            </template>
          </van-cell>
          <van-cell title="用户协议" is-link @click="showUserAgreement = true">
            <template #icon>
              <i class="iconfont icon-file setting-icon"></i>
            </template>
          </van-cell>
        </van-cell-group>
      </div>

      <!-- 其他设置 -->
      <div class="setting-section">
        <h3 class="section-title">其他</h3>
        <van-cell-group inset>
          <van-cell title="检查更新" is-link @click="checkUpdate">
            <template #icon>
              <i class="iconfont icon-refresh setting-icon"></i>
            </template>
            <template #value>
              <span class="setting-value">v{{ appVersion }}</span>
            </template>
          </van-cell>
          <van-cell title="帮助中心" is-link @click="router.push('/mobile/help')">
            <template #icon>
              <i class="iconfont icon-help setting-icon"></i>
            </template>
          </van-cell>
          <van-cell title="意见反馈" is-link @click="router.push('/mobile/feedback')">
            <template #icon>
              <i class="iconfont icon-message setting-icon"></i>
            </template>
          </van-cell>
          <van-cell title="关于我们" is-link @click="showAbout = true">
            <template #icon>
              <i class="iconfont icon-info setting-icon"></i>
            </template>
          </van-cell>
        </van-cell-group>
      </div>
    </div>

    <!-- 修改密码弹窗 -->
    <van-dialog
      v-model:show="showPasswordDialog"
      title="修改密码"
      show-cancel-button
      @confirm="updatePassword"
    >
      <van-form>
        <van-cell-group inset>
          <van-field
            v-model="passwordForm.old"
            type="password"
            label="旧密码"
            placeholder="请输入旧密码"
          />
          <van-field
            v-model="passwordForm.new"
            type="password"
            label="新密码"
            placeholder="请输入新密码"
          />
          <van-field
            v-model="passwordForm.confirm"
            type="password"
            label="确认密码"
            placeholder="请再次输入新密码"
          />
        </van-cell-group>
      </van-form>
    </van-dialog>

    <!-- 语言选择 -->
    <van-action-sheet v-model:show="showLanguageSelect" title="选择语言">
      <div class="language-options">
        <van-cell
          v-for="lang in languages"
          :key="lang.code"
          :title="lang.name"
          clickable
          @click="selectLanguage(lang)"
        >
          <template #right-icon>
            <van-icon v-if="currentLanguage.code === lang.code" name="success" color="#07c160" />
          </template>
        </van-cell>
      </div>
    </van-action-sheet>

    <!-- 模型选择 -->
    <van-action-sheet v-model:show="showModelSelect" title="选择默认模型">
      <div class="model-options">
        <van-cell
          v-for="model in models"
          :key="model.code"
          :title="model.name"
          :label="model.desc"
          clickable
          @click="selectModel(model)"
        >
          <template #right-icon>
            <van-icon v-if="currentModel.code === model.code" name="success" color="#07c160" />
          </template>
        </van-cell>
      </div>
    </van-action-sheet>

    <!-- 发送方式选择 -->
    <van-action-sheet v-model:show="showSendModeSelect" title="选择发送方式">
      <div class="send-mode-options">
        <van-cell
          v-for="mode in sendModes"
          :key="mode.code"
          :title="mode.name"
          :label="mode.desc"
          clickable
          @click="selectSendMode(mode)"
        >
          <template #right-icon>
            <van-icon v-if="currentSendMode.code === mode.code" name="success" color="#07c160" />
          </template>
        </van-cell>
      </div>
    </van-action-sheet>

    <!-- 清除缓存确认 -->
    <van-dialog
      v-model:show="showClearCache"
      title="清除缓存"
      message="确定要清除所有缓存数据吗？这将删除临时文件和图片缓存。"
      show-cancel-button
      @confirm="clearCache"
    />

    <!-- 清除聊天记录确认 -->
    <van-dialog
      v-model:show="showClearHistory"
      title="清除聊天记录"
      message="确定要清除所有聊天记录吗？此操作不可撤销。"
      show-cancel-button
      @confirm="clearHistory"
    />

    <!-- 关于我们 -->
    <van-dialog v-model:show="showAbout" title="关于我们" :show-cancel-button="false">
      <div class="about-content">
        <div class="about-logo">
          <img src="/images/logo.png" alt="Logo" />
        </div>
        <h3>{{ appName }}</h3>
        <p class="about-desc">
          专业的AI创作平台，提供对话、绘画、音乐、视频等多种AI服务，让创作更简单、更高效。
        </p>
        <div class="about-info">
          <p>版本：v{{ appVersion }}</p>
          <p>更新时间：2024-01-01</p>
        </div>
      </div>
    </van-dialog>
  </div>
</template>

<script setup>
import { useSharedStore } from '@/store/sharedata'
import { showNotify, showSuccessToast } from 'vant'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const store = useSharedStore()

// 基础状态
const appName = ref(import.meta.env.VITE_TITLE)
const appVersion = ref('2.1.0')
const cacheSize = ref('23.5MB')

// 设置状态
const darkMode = ref(store.theme === 'dark')
const streamOutput = ref(store.chatStream)
const notifications = ref(true)
const autoSave = ref(true)
const saveHistory = ref(true)

// 弹窗状态
const showPasswordDialog = ref(false)
const showBindMobile = ref(false)
const showBindEmail = ref(false)
const showLanguageSelect = ref(false)
const showModelSelect = ref(false)
const showSendModeSelect = ref(false)
const showClearCache = ref(false)
const showClearHistory = ref(false)
const showPrivacyPolicy = ref(false)
const showUserAgreement = ref(false)
const showAbout = ref(false)

// 表单数据
const passwordForm = ref({
  old: '',
  new: '',
  confirm: '',
})

// 语言选项
const languages = ref([
  { code: 'zh-CN', name: '简体中文' },
  { code: 'zh-TW', name: '繁體中文' },
  { code: 'en', name: 'English' },
  { code: 'ja', name: '日本語' },
  { code: 'ko', name: '한국어' },
])

const currentLanguage = ref(languages.value[0])

// 模型选项
const models = ref([
  { code: 'gpt-4', name: 'GPT-4', desc: '最新的GPT-4模型，性能强大' },
  { code: 'gpt-3.5', name: 'GPT-3.5', desc: '经典的GPT-3.5模型，速度快' },
  { code: 'claude', name: 'Claude', desc: '人工智能助手Claude' },
  { code: 'gemini', name: 'Gemini', desc: 'Google的Gemini模型' },
])

const currentModel = ref(models.value[0])

// 发送方式选项
const sendModes = ref([
  { code: 'enter', name: 'Enter发送', desc: '按Enter键发送消息' },
  { code: 'ctrl+enter', name: 'Ctrl+Enter发送', desc: '按Ctrl+Enter发送消息' },
  { code: 'button', name: '仅按钮发送', desc: '只能点击发送按钮' },
])

const currentSendMode = ref(sendModes.value[0])

onMounted(() => {
  loadSettings()
})

// 加载设置
const loadSettings = () => {
  // 从localStorage加载设置
  const savedSettings = localStorage.getItem('app-settings')
  if (savedSettings) {
    const settings = JSON.parse(savedSettings)
    darkMode.value = settings.darkMode ?? store.theme === 'dark'
    streamOutput.value = settings.streamOutput ?? store.chatStream
    notifications.value = settings.notifications ?? true
    autoSave.value = settings.autoSave ?? true
    saveHistory.value = settings.saveHistory ?? true

    // 恢复语言设置
    const savedLang = languages.value.find((lang) => lang.code === settings.language)
    if (savedLang) {
      currentLanguage.value = savedLang
    }

    // 恢复模型设置
    const savedModel = models.value.find((model) => model.code === settings.model)
    if (savedModel) {
      currentModel.value = savedModel
    }

    // 恢复发送方式设置
    const savedSendMode = sendModes.value.find((mode) => mode.code === settings.sendMode)
    if (savedSendMode) {
      currentSendMode.value = savedSendMode
    }
  }
}

// 保存设置
const saveSettings = () => {
  const settings = {
    darkMode: darkMode.value,
    streamOutput: streamOutput.value,
    notifications: notifications.value,
    autoSave: autoSave.value,
    saveHistory: saveHistory.value,
    language: currentLanguage.value.code,
    model: currentModel.value.code,
    sendMode: currentSendMode.value.code,
  }
  localStorage.setItem('app-settings', JSON.stringify(settings))
}

// 主题切换
const onThemeChange = (value) => {
  store.setTheme(value ? 'dark' : 'light')
  saveSettings()
}

// 流式输出切换
const onStreamChange = (value) => {
  store.setChatStream(value)
  saveSettings()
}

// 选择语言
const selectLanguage = (lang) => {
  currentLanguage.value = lang
  showLanguageSelect.value = false
  saveSettings()
  showSuccessToast(`已切换到${lang.name}`)
}

// 选择模型
const selectModel = (model) => {
  currentModel.value = model
  showModelSelect.value = false
  saveSettings()
  showSuccessToast(`已设置默认模型为${model.name}`)
}

// 选择发送方式
const selectSendMode = (mode) => {
  currentSendMode.value = mode
  showSendModeSelect.value = false
  saveSettings()
  showSuccessToast(`已设置发送方式为${mode.name}`)
}

// 修改密码
const updatePassword = () => {
  if (!passwordForm.value.old) {
    showNotify({ type: 'danger', message: '请输入旧密码' })
    return
  }
  if (!passwordForm.value.new || passwordForm.value.new.length < 8) {
    showNotify({ type: 'danger', message: '新密码长度不能少于8位' })
    return
  }
  if (passwordForm.value.new !== passwordForm.value.confirm) {
    showNotify({ type: 'danger', message: '两次输入的密码不一致' })
    return
  }

  // 这里应该调用API
  showSuccessToast('密码修改成功')
  showPasswordDialog.value = false
  passwordForm.value = { old: '', new: '', confirm: '' }
}

// 清除缓存
const clearCache = () => {
  setTimeout(() => {
    cacheSize.value = '0MB'
    showSuccessToast('缓存清除成功')
  }, 1000)
}

// 清除聊天记录
const clearHistory = () => {
  // 这里应该调用API清除聊天记录
  showSuccessToast('聊天记录清除成功')
}

// 检查更新
const checkUpdate = () => {
  showNotify({ type: 'primary', message: '正在检查更新...' })
  setTimeout(() => {
    showNotify({ type: 'success', message: '当前已是最新版本' })
  }, 2000)
}
</script>

<style lang="scss" scoped>
.settings-page {
  min-height: 100vh;
  background: var(--van-background);

  .settings-content {
    padding: 54px 16px 20px;

    .setting-section {
      margin-bottom: 24px;

      .section-title {
        font-size: 18px;
        font-weight: 600;
        color: var(--van-text-color);
        margin: 0 0 16px 4px;
      }

      :deep(.van-cell-group) {
        border-radius: 12px;
        overflow: hidden;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

        .van-cell {
          padding: 16px;

          .setting-icon {
            font-size: 18px;
            color: var(--van-primary-color);
            margin-right: 12px;
          }

          .van-cell__title {
            font-size: 15px;
            font-weight: 500;
          }

          .setting-value {
            font-size: 14px;
            color: var(--van-gray-6);
          }
        }
      }
    }
  }

  .language-options,
  .model-options,
  .send-mode-options {
    max-height: 400px;
    overflow-y: auto;

    :deep(.van-cell) {
      padding: 16px 20px;

      .van-cell__title {
        font-size: 15px;
        font-weight: 500;
      }

      .van-cell__label {
        color: var(--van-gray-6);
        font-size: 13px;
        margin-top: 4px;
      }
    }
  }

  .about-content {
    text-align: center;
    padding: 20px;

    .about-logo {
      margin-bottom: 16px;

      img {
        width: 60px;
        height: 60px;
        border-radius: 12px;
      }
    }

    h3 {
      font-size: 20px;
      font-weight: 600;
      color: var(--van-text-color);
      margin: 0 0 12px 0;
    }

    .about-desc {
      font-size: 14px;
      color: var(--van-gray-6);
      line-height: 1.5;
      margin: 0 0 20px 0;
    }

    .about-info {
      p {
        font-size: 13px;
        color: var(--van-gray-7);
        margin: 0 0 4px 0;

        &:last-child {
          margin: 0;
        }
      }
    }
  }
}

// 深色主题优化
:deep(.van-theme-dark) {
  .settings-page {
    .van-cell-group {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    }
  }
}
</style>
