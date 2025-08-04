<template>
  <div class="invite-page">
    <div class="invite-content">
      <!-- 邀请头图 -->
      <div class="invite-header">
        <div class="header-bg">
          <img src="/images/invite-bg.png" alt="邀请背景" @error="onImageError" />
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

      <!-- 邀请方式 -->
      <div class="invite-methods">
        <h3 class="section-title">邀请方式</h3>
        <div class="methods-grid">
          <div class="method-item" @click="shareToWeChat">
            <div class="method-icon wechat">
              <i class="iconfont icon-wechat"></i>
            </div>
            <div class="method-name">微信分享</div>
          </div>
          <div class="method-item" @click="copyInviteLink">
            <div class="method-icon link">
              <i class="iconfont icon-link"></i>
            </div>
            <div class="method-name">复制链接</div>
          </div>
          <div class="method-item" @click="shareQRCode">
            <div class="method-icon qr">
              <i class="iconfont icon-qrcode"></i>
            </div>
            <div class="method-name">二维码</div>
          </div>
          <div class="method-item" @click="shareToFriends">
            <div class="method-icon more">
              <i class="iconfont icon-share"></i>
            </div>
            <div class="method-name">更多</div>
          </div>
        </div>
      </div>

      <!-- 邀请码 -->
      <div class="invite-code-section">
        <div class="code-card">
          <div class="code-header">
            <span class="code-label">我的邀请码</span>
            <van-button size="small" type="primary" plain @click="copyInviteCode">
              复制
            </van-button>
          </div>
          <div class="code-value">{{ inviteCode }}</div>
          <div class="code-link">
            <van-field v-model="inviteLink" readonly placeholder="邀请链接">
              <template #button>
                <van-button size="small" type="primary" @click="copyInviteLink">
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
          <van-button size="small" type="primary" plain @click="showAllRecords = !showAllRecords">
            {{ showAllRecords ? '收起' : '查看全部' }}
          </van-button>
        </div>

        <div class="records-list">
          <van-list
            v-model:loading="recordsLoading"
            :finished="recordsFinished"
            finished-text="没有更多记录"
            @load="loadInviteRecords"
          >
            <div v-for="record in displayRecords" :key="record.id" class="record-item">
              <div class="record-avatar">
                <van-image :src="record.avatar" round width="40" height="40" />
              </div>
              <div class="record-info">
                <div class="record-name">{{ record.username }}</div>
                <div class="record-time">{{ formatTime(record.created_at) }}</div>
              </div>
              <div class="record-status">
                <van-tag :type="record.status === 'completed' ? 'success' : 'warning'">
                  {{ record.status === 'completed' ? '已获得奖励' : '待获得奖励' }}
                </van-tag>
              </div>
            </div>

            <van-empty
              v-if="!recordsLoading && inviteRecords.length === 0"
              description="暂无邀请记录"
            />
          </van-list>
        </div>
      </div>
    </div>

    <!-- 二维码弹窗 -->
    <van-dialog
      v-model:show="showQRDialog"
      title="邀请二维码"
      :show-cancel-button="false"
      confirm-button-text="保存图片"
      @confirm="saveQRCode"
    >
      <div class="qr-content">
        <div ref="qrCodeRef" class="qr-code">
          <!-- 这里应该生成实际的二维码 -->
          <div class="qr-placeholder">
            <i class="iconfont icon-qrcode"></i>
            <p>邀请二维码</p>
          </div>
        </div>
        <p class="qr-tip">扫描二维码或长按保存分享给好友</p>
      </div>
    </van-dialog>
  </div>
</template>

<script setup>
import { checkSession } from '@/store/cache'
import { showLoginDialog } from '@/utils/libs'
import { showNotify, showSuccessToast } from 'vant'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const userStats = ref({
  inviteCount: 0,
  rewardTotal: 0,
  todayInvite: 0,
})
const inviteCode = ref('')
const inviteLink = ref('')
const inviteRecords = ref([])
const recordsLoading = ref(false)
const recordsFinished = ref(false)
const showAllRecords = ref(false)
const showQRDialog = ref(false)
const qrCodeRef = ref()

// 奖励规则配置
const rewardRules = ref([
  {
    id: 1,
    title: '好友注册',
    desc: '好友通过邀请链接成功注册',
    icon: 'icon-user-plus',
    color: '#1989fa',
    reward: 50,
  },
  {
    id: 2,
    title: '好友首次充值',
    desc: '好友首次充值任意金额',
    icon: 'icon-money',
    color: '#07c160',
    reward: 100,
  },
  {
    id: 3,
    title: '好友活跃使用',
    desc: '好友连续使用7天',
    icon: 'icon-star',
    color: '#ff9500',
    reward: 200,
  },
])

// 显示的记录（根据showAllRecords决定）
const displayRecords = computed(() => {
  return showAllRecords.value ? inviteRecords.value : inviteRecords.value.slice(0, 5)
})

onMounted(() => {
  initPage()
})

const initPage = async () => {
  try {
    const user = await checkSession()

    // 生成邀请码和链接
    inviteCode.value = user.invite_code || generateInviteCode()
    inviteLink.value = `${location.origin}/register?invite=${inviteCode.value}`

    // 获取用户邀请统计
    fetchInviteStats()

    // 加载邀请记录
    loadInviteRecords()
  } catch (error) {
    showLoginDialog(router)
  }
}

const generateInviteCode = () => {
  return Math.random().toString(36).substr(2, 8).toUpperCase()
}

const fetchInviteStats = () => {
  // 这里应该调用实际的API
  // httpGet('/api/user/invite/stats').then(res => {
  //   userStats.value = res.data
  // })

  // 临时使用模拟数据
  userStats.value = {
    inviteCount: Math.floor(Math.random() * 50),
    rewardTotal: Math.floor(Math.random() * 5000),
    todayInvite: Math.floor(Math.random() * 5),
  }
}

const loadInviteRecords = () => {
  if (recordsFinished.value) return

  recordsLoading.value = true

  // 模拟API调用
  setTimeout(() => {
    const mockRecords = generateMockRecords()
    inviteRecords.value.push(...mockRecords)

    recordsLoading.value = false

    // 模拟数据加载完成
    if (inviteRecords.value.length >= 20) {
      recordsFinished.value = true
    }
  }, 1000)
}

const generateMockRecords = () => {
  const records = []
  const names = ['张三', '李四', '王五', '赵六', '钱七', '孙八', '周九', '吴十']

  for (let i = 0; i < 10; i++) {
    records.push({
      id: Date.now() + i,
      username: names[i % names.length] + (i + 1),
      avatar: '/images/avatar/default.jpg',
      status: Math.random() > 0.3 ? 'completed' : 'pending',
      created_at: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
    })
  }

  return records
}

const formatTime = (timeStr) => {
  const date = new Date(timeStr)
  return date.toLocaleDateString()
}

// 分享到微信
const shareToWeChat = () => {
  if (typeof WeixinJSBridge !== 'undefined') {
    // 在微信中分享
    WeixinJSBridge.invoke('sendAppMessage', {
      title: '邀请你使用AI创作平台',
      desc: '强大的AI工具，让创作更简单',
      link: inviteLink.value,
      imgUrl: `${location.origin}/images/share-logo.png`,
    })
  } else {
    // 复制链接提示
    copyInviteLink()
    showNotify({ type: 'primary', message: '请在微信中打开链接进行分享' })
  }
}

// 复制邀请码
const copyInviteCode = async () => {
  try {
    await navigator.clipboard.writeText(inviteCode.value)
    showSuccessToast('邀请码已复制')
  } catch (err) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = inviteCode.value
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    showSuccessToast('邀请码已复制')
  }
}

// 复制邀请链接
const copyInviteLink = async () => {
  try {
    await navigator.clipboard.writeText(inviteLink.value)
    showSuccessToast('邀请链接已复制')
  } catch (err) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = inviteLink.value
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    showSuccessToast('邀请链接已复制')
  }
}

// 显示二维码
const shareQRCode = () => {
  showQRDialog.value = true
}

// 保存二维码
const saveQRCode = () => {
  showNotify({ type: 'primary', message: '请长按二维码保存到相册' })
}

// 更多分享方式
const shareToFriends = () => {
  if (navigator.share) {
    navigator.share({
      title: '邀请你使用AI创作平台',
      text: '强大的AI工具，让创作更简单',
      url: inviteLink.value,
    })
  } else {
    copyInviteLink()
  }
}

// 图片加载错误处理
const onImageError = (e) => {
  e.target.src = '/images/default-bg.png'
}
</script>

<style lang="scss" scoped>
.invite-page {
  min-height: 100vh;
  background: var(--van-background);

  .invite-content {
    padding-top: 46px;

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
    .invite-methods,
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

    .methods-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 16px;

      .method-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 20px 16px;
        background: var(--van-cell-background);
        border-radius: 12px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
        cursor: pointer;
        transition: all 0.3s ease;

        &:active {
          transform: scale(0.95);
        }

        .method-icon {
          width: 44px;
          height: 44px;
          border-radius: 12px;
          display: flex;
          align-items: center;
          justify-content: center;
          margin-bottom: 8px;

          &.wechat {
            background: #07c160;
          }

          &.link {
            background: #1989fa;
          }

          &.qr {
            background: #8b5cf6;
          }

          &.more {
            background: #ff9500;
          }

          .iconfont {
            font-size: 20px;
            color: white;
          }
        }

        .method-name {
          font-size: 12px;
          color: var(--van-text-color);
          text-align: center;
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

  .qr-content {
    text-align: center;
    padding: 20px;

    .qr-code {
      width: 200px;
      height: 200px;
      margin: 0 auto 16px;
      border: 1px solid var(--van-border-color);
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: var(--van-background-2);

      .qr-placeholder {
        text-align: center;
        color: var(--van-gray-6);

        .iconfont {
          font-size: 48px;
          margin-bottom: 8px;
        }

        p {
          margin: 0;
          font-size: 14px;
        }
      }
    }

    .qr-tip {
      font-size: 13px;
      color: var(--van-gray-6);
      margin: 0;
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
