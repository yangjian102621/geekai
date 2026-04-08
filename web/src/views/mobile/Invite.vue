<template>
  <div class="invite-page">
    <div class="invite-content">
      <!-- 邀请头图 -->
      <div class="invite-header">
        <div class="header-bg">
          <img src="/images/logo.png" alt="邀请背景" @error="onImageError" />
        </div>
        <div class="header-content">
          <h2 class="invite-title">邀请好友获得奖励</h2>
          <p class="invite-desc">邀请好友注册即可获得算力奖励</p>
        </div>
      </div>

      <!-- 用户统计 -->
      <div class="stats-section">
        <van-row :gutter="12">
          <van-col :span="8">
            <div class="stat-card">
              <div class="stat-number">{{ userStats.inviteCount }}</div>
              <div class="stat-label">累计邀请</div>
            </div>
          </van-col>
          <van-col :span="8">
            <div class="stat-card">
              <div class="stat-number">{{ userStats.rewardTotal }}</div>
              <div class="stat-label">获得奖励</div>
            </div>
          </van-col>
          <van-col :span="8">
            <div class="stat-card">
              <div class="stat-number">{{ userStats.todayInvite }}</div>
              <div class="stat-label">今日邀请</div>
            </div>
          </van-col>
        </van-row>
      </div>

      <!-- 奖励规则 -->
      <div class="rules-section">
        <h3 class="section-title">奖励规则</h3>
        <div class="rules-list">
          <div class="rule-item" v-for="rule in rewardRules" :key="rule.id">
            <div class="rule-icon">
              <i class="iconfont" :class="rule.icon" :style="{ color: rule.color }"></i>
            </div>
            <div class="rule-content">
              <div class="rule-title">{{ rule.title }}</div>
              <div class="rule-desc">{{ rule.desc }}</div>
            </div>
            <div class="rule-reward">
              <span class="reward-value">+{{ rule.reward }}</span>
              <span class="reward-unit">算力</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 邀请码 -->
      <div class="invite-code-section">
        <div class="code-card">
          <div class="code-header">
            <span class="code-label">我的邀请码</span>
            <van-button
              size="small"
              type="primary"
              plain
              class="copy-invite-code"
              :data-clipboard-text="inviteCode"
            >
              复制
            </van-button>
          </div>
          <div class="code-value">{{ inviteCode }}</div>
          <div class="code-link">
            <van-field v-model="inviteLink" readonly placeholder="邀请链接">
              <template #button>
                <van-button
                  size="small"
                  type="primary"
                  class="copy-invite-link"
                  :data-clipboard-text="inviteLink"
                >
                  复制链接
                </van-button>
              </template>
            </van-field>
          </div>
        </div>
      </div>

      <!-- 邀请记录 -->
      <div class="invite-records">
        <div class="records-header">
          <h3 class="section-title">邀请记录</h3>
        </div>

        <div class="records-list">
          <div v-if="inviteRecords.length > 0">
            <div v-for="record in inviteRecords" :key="record.id" class="record-item">
              <div class="record-avatar">
                <van-image
                  :src="record.avatar || '/images/avatar/default.jpg'"
                  round
                  width="40"
                  height="40"
                  @error="onAvatarError"
                />
              </div>
              <div class="record-info">
                <div class="record-name">{{ record.username }}</div>
                <div class="record-time">{{ formatTime(record.created_at) }}</div>
              </div>
              <div class="record-status">
                <van-tag type="success">已获得奖励</van-tag>
              </div>
            </div>
          </div>

          <van-empty v-else description="暂无邀请记录" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { checkSession } from '@/store/cache'
import { showLoginDialog } from '@/utils/libs'
import { httpGet } from '@/utils/http'
import { showNotify, showSuccessToast } from 'vant'
import { onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import Clipboard from 'clipboard'

const router = useRouter()
const userStats = ref({
  inviteCount: 0,
  rewardTotal: 0,
  todayInvite: 0,
})
const inviteCode = ref('')
const inviteLink = ref('')
const inviteRecords = ref([])

// 奖励规则配置
const rewardRules = ref([])

// clipboard实例
const clipboard = ref(null)

onMounted(() => {
  initPage()
})

onUnmounted(() => {
  // 清理clipboard实例
  if (clipboard.value) {
    clipboard.value.destroy()
  }
})

const initPage = async () => {
  try {
    const user = await checkSession()

    // 获取邀请统计（包含邀请码和链接）
    fetchInviteStats()

    // 获取奖励规则
    fetchRewardRules()

    // 一次性加载邀请记录
    await loadInviteRecords()

    // 初始化clipboard
    initClipboard()
  } catch (error) {
    showLoginDialog(router)
  }
}

const fetchInviteStats = async () => {
  try {
    const res = await httpGet('/api/invite/stats')
    userStats.value = {
      inviteCount: res.data.invite_count,
      rewardTotal: res.data.reward_total,
      todayInvite: res.data.today_invite,
    }
    inviteCode.value = res.data.invite_code
    inviteLink.value = res.data.invite_link
  } catch (error) {
    console.error('获取邀请统计失败:', error)
    // 使用默认值
    userStats.value = {
      inviteCount: 0,
      rewardTotal: 0,
      todayInvite: 0,
    }
  }
}

const fetchRewardRules = async () => {
  try {
    const res = await httpGet('/api/invite/rules')
    rewardRules.value = res.data
  } catch (error) {
    console.error('获取奖励规则失败:', error)
  }
}

const loadInviteRecords = async () => {
  try {
    const res = await httpGet('/api/invite/list', {
      page: 1, // 加载第一页
      page_size: 20,
    })
    console.log('邀请记录API返回:', res.data) // 调试日志
    inviteRecords.value = res.data.items || []
    console.log('设置邀请记录:', inviteRecords.value) // 调试日志

    // 调试头像信息
    if (inviteRecords.value.length > 0) {
      console.log('第一条记录头像:', inviteRecords.value[0].avatar)
    }
  } catch (error) {
    console.error('获取邀请记录失败:', error)
  }
}

const formatTime = (timestamp) => {
  const date = new Date(timestamp * 1000) // 转换为毫秒
  return date.toLocaleDateString()
}

// 初始化clipboard
const initClipboard = () => {
  // 销毁之前的实例
  if (clipboard.value) {
    clipboard.value.destroy()
  }

  // 创建新的clipboard实例
  clipboard.value = new Clipboard('.copy-invite-code, .copy-invite-link')

  clipboard.value.on('success', () => {
    showSuccessToast('复制成功')
  })

  clipboard.value.on('error', () => {
    showNotify({ type: 'danger', message: '复制失败' })
  })
}

// 图片加载错误处理
const onImageError = (e) => {
  e.target.src = '/images/img-holder.png'
}

// 头像加载错误处理
const onAvatarError = (e) => {
  e.target.src = '/images/avatar/default.jpg'
}
</script>

<style lang="scss" scoped>
.invite-page {
  min-height: 100vh;
  background: var(--van-background);

  .invite-content {
    .invite-header {
      position: relative;
      height: 200px;
      overflow: hidden;
      background: linear-gradient(135deg, var(--van-primary-color), #8b5cf6);

      .header-bg {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        opacity: 0.3;

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }
      }

      .header-content {
        position: relative;
        z-index: 2;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        height: 100%;
        color: white;
        text-align: center;

        .invite-title {
          font-size: 24px;
          font-weight: 700;
          margin: 0 0 8px 0;
        }

        .invite-desc {
          font-size: 14px;
          opacity: 0.9;
          margin: 0;
        }
      }
    }

    .stats-section {
      padding: 16px;
      margin-top: -20px;
      position: relative;
      z-index: 3;

      .stat-card {
        background: var(--van-cell-background);
        border-radius: 12px;
        padding: 16px;
        text-align: center;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);

        .stat-number {
          font-size: 20px;
          font-weight: 700;
          color: var(--van-primary-color);
          margin-bottom: 4px;
        }

        .stat-label {
          font-size: 12px;
          color: var(--van-gray-6);
        }
      }
    }

    .rules-section,
    .invite-code-section,
    .invite-records {
      padding: 0 16px 16px;

      .section-title {
        font-size: 18px;
        font-weight: 600;
        color: var(--van-text-color);
        margin: 0 0 16px 0;
      }
    }

    .rules-list {
      .rule-item {
        display: flex;
        align-items: center;
        background: var(--van-cell-background);
        border-radius: 12px;
        padding: 16px;
        margin-bottom: 12px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

        .rule-icon {
          width: 40px;
          height: 40px;
          border-radius: 10px;
          background: rgba(25, 137, 250, 0.1);
          display: flex;
          align-items: center;
          justify-content: center;
          margin-right: 12px;

          .iconfont {
            font-size: 20px;
          }
        }

        .rule-content {
          flex: 1;

          .rule-title {
            font-size: 15px;
            font-weight: 600;
            color: var(--van-text-color);
            margin-bottom: 4px;
          }

          .rule-desc {
            font-size: 13px;
            color: var(--van-gray-6);
          }
        }

        .rule-reward {
          text-align: right;

          .reward-value {
            font-size: 16px;
            font-weight: 600;
            color: #07c160;
          }

          .reward-unit {
            font-size: 12px;
            color: var(--van-gray-6);
            margin-left: 2px;
          }
        }
      }
    }

    .invite-code-section {
      .code-card {
        background: var(--van-cell-background);
        border-radius: 12px;
        padding: 20px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

        .code-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 12px;

          .code-label {
            font-size: 16px;
            font-weight: 600;
            color: var(--van-text-color);
          }
        }

        .code-value {
          font-size: 24px;
          font-weight: 700;
          color: var(--van-primary-color);
          text-align: center;
          padding: 16px;
          background: rgba(25, 137, 250, 0.1);
          border-radius: 8px;
          margin-bottom: 16px;
          letter-spacing: 2px;
        }
      }
    }

    .invite-records {
      .records-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16px;
      }

      .records-list {
        .record-item {
          display: flex;
          align-items: center;
          background: var(--van-cell-background);
          border-radius: 12px;
          padding: 16px;
          margin-bottom: 12px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

          .record-avatar {
            margin-right: 12px;
          }

          .record-info {
            flex: 1;

            .record-name {
              font-size: 15px;
              font-weight: 500;
              color: var(--van-text-color);
              margin-bottom: 4px;
            }

            .record-time {
              font-size: 12px;
              color: var(--van-gray-6);
            }
          }
        }
      }
    }
  }
}

// 深色主题优化
:deep(.van-theme-dark) {
  .invite-page {
    .stat-card,
    .rule-item,
    .method-item,
    .code-card,
    .record-item {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    }
  }
}
</style>
