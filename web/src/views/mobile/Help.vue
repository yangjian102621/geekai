<template>
  <div class="help-page">
    <van-nav-bar title="帮助中心" left-arrow @click-left="router.back()" fixed>
      <template #right>
        <van-icon name="search" @click="showSearch = true" />
      </template>
    </van-nav-bar>

    <div class="help-content">
      <!-- 搜索框 -->
      <div class="search-section" v-if="showSearch">
        <van-search
          v-model="searchKeyword"
          placeholder="搜索帮助内容"
          @search="onSearch"
          @cancel="showSearch = false"
          show-action
        />
      </div>

      <!-- 常见问题 -->
      <div class="faq-section" v-if="!showSearch">
        <h3 class="section-title">常见问题</h3>
        <van-collapse v-model="activeNames" accordion>
          <van-collapse-item
            v-for="faq in frequentFAQs"
            :key="faq.id"
            :title="faq.question"
            :name="faq.id"
            class="faq-item"
          >
            <div class="faq-answer" v-html="faq.answer"></div>
          </van-collapse-item>
        </van-collapse>
      </div>

      <!-- 功能指南 -->
      <div class="guide-section" v-if="!showSearch">
        <h3 class="section-title">功能指南</h3>
        <van-grid :column-num="2" :gutter="12" :border="false">
          <van-grid-item
            v-for="guide in guides"
            :key="guide.id"
            @click="openGuide(guide)"
            class="guide-item"
          >
            <div class="guide-card">
              <div class="guide-icon" :style="{ backgroundColor: guide.color }">
                <i class="iconfont" :class="guide.icon"></i>
              </div>
              <div class="guide-title">{{ guide.title }}</div>
              <div class="guide-desc">{{ guide.desc }}</div>
            </div>
          </van-grid-item>
        </van-grid>
      </div>

      <!-- 问题分类 -->
      <div class="category-section" v-if="!showSearch">
        <h3 class="section-title">问题分类</h3>
        <van-cell-group inset>
          <van-cell
            v-for="category in categories"
            :key="category.id"
            :title="category.name"
            :value="`${category.count}个问题`"
            is-link
            @click="openCategory(category)"
          >
            <template #icon>
              <i class="iconfont" :class="category.icon" class="category-icon"></i>
            </template>
          </van-cell>
        </van-cell-group>
      </div>

      <!-- 搜索结果 -->
      <div class="search-results" v-if="showSearch && searchResults.length > 0">
        <h3 class="section-title">搜索结果</h3>
        <van-list>
          <van-cell
            v-for="result in searchResults"
            :key="result.id"
            :title="result.title"
            @click="openSearchResult(result)"
            is-link
          >
            <template #label>
              <div class="search-snippet" v-html="result.snippet"></div>
            </template>
          </van-cell>
        </van-list>
      </div>

      <!-- 空搜索结果 -->
      <van-empty
        v-if="showSearch && searchKeyword && searchResults.length === 0"
        description="没有找到相关内容"
      />

      <!-- 联系客服 -->
      <div class="contact-section" v-if="!showSearch">
        <h3 class="section-title">联系我们</h3>
        <van-cell-group inset>
          <van-cell title="在线客服" icon="service-o" is-link @click="openCustomerService">
            <template #value>
              <span class="online-status">在线</span>
            </template>
          </van-cell>
          <van-cell title="意见反馈" icon="chat-o" is-link @click="router.push('/mobile/feedback')" />
          <van-cell title="官方QQ群" icon="friends-o" is-link @click="joinQQGroup">
            <template #value>
              <span class="qq-number">123456789</span>
            </template>
          </van-cell>
          <van-cell title="官方微信" icon="wechat" is-link @click="showWeChatQR = true" />
        </van-cell-group>
      </div>

      <!-- 使用提示 -->
      <div class="tips-section" v-if="!showSearch">
        <h3 class="section-title">使用提示</h3>
        <van-swipe :autoplay="5000" class="tips-swipe">
          <van-swipe-item v-for="tip in tips" :key="tip.id">
            <div class="tip-card">
              <div class="tip-icon">
                <i class="iconfont" :class="tip.icon"></i>
              </div>
              <h4 class="tip-title">{{ tip.title }}</h4>
              <p class="tip-content">{{ tip.content }}</p>
            </div>
          </van-swipe-item>
        </van-swipe>
      </div>
    </div>

    <!-- 帮助详情弹窗 -->
    <van-action-sheet v-model:show="showHelpDetail" :title="selectedHelp?.title">
      <div class="help-detail" v-if="selectedHelp">
        <div class="detail-content" v-html="selectedHelp.content"></div>
        <div class="detail-actions">
          <van-button @click="likeHelp(selectedHelp)">
            <van-icon name="good-job-o" /> 有用
          </van-button>
          <van-button @click="shareHelp(selectedHelp)">
            <van-icon name="share-o" /> 分享
          </van-button>
        </div>
      </div>
    </van-action-sheet>

    <!-- 微信二维码弹窗 -->
    <van-dialog v-model:show="showWeChatQR" title="官方微信" :show-cancel-button="false">
      <div class="wechat-qr">
        <div class="qr-code">
          <img src="/images/wechat-qr.png" alt="微信二维码" @error="onQRError" />
        </div>
        <p class="qr-tip">扫描二维码添加官方微信</p>
      </div>
    </van-dialog>

    <!-- 客服聊天 -->
    <van-action-sheet v-model:show="showCustomerChat" title="在线客服" :close-on-click-overlay="false">
      <div class="customer-chat">
        <div class="chat-header">
          <div class="customer-info">
            <van-image src="/images/customer-service.png" round width="40" height="40" />
            <div class="customer-detail">
              <div class="customer-name">智能客服</div>
              <div class="customer-status online">在线</div>
            </div>
          </div>
          <van-button size="small" @click="showCustomerChat = false">结束</van-button>
        </div>
        
        <div class="chat-messages" ref="chatMessages">
          <div
            v-for="message in customerMessages"
            :key="message.id"
            class="message-item"
            :class="{ 'user-message': message.isUser }"
          >
            <div class="message-content">{{ message.content }}</div>
            <div class="message-time">{{ formatTime(message.time) }}</div>
          </div>
        </div>
        
        <div class="chat-input">
          <van-field
            v-model="customerMessage"
            placeholder="请输入您的问题..."
            @keyup.enter="sendCustomerMessage"
          >
            <template #button>
              <van-button size="small" type="primary" @click="sendCustomerMessage">
                发送
              </van-button>
            </template>
          </van-field>
        </div>
      </div>
    </van-action-sheet>
  </div>
</template>

<script setup>
import { showNotify, showSuccessToast } from 'vant'
import { nextTick, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const showSearch = ref(false)
const searchKeyword = ref('')
const searchResults = ref([])
const activeNames = ref([])
const selectedHelp = ref(null)
const showHelpDetail = ref(false)
const showWeChatQR = ref(false)
const showCustomerChat = ref(false)
const customerMessage = ref('')
const customerMessages = ref([])
const chatMessages = ref()

// 常见问题
const frequentFAQs = ref([
  {
    id: 1,
    question: '如何获得算力？',
    answer: '<p>您可以通过以下方式获得算力：</p><ul><li>注册即送算力</li><li>购买充值套餐</li><li>邀请好友注册</li><li>参与活动获得</li></ul>'
  },
  {
    id: 2,
    question: '如何使用AI绘画功能？',
    answer: '<p>使用AI绘画功能很简单：</p><ol><li>进入创作中心</li><li>选择绘画工具（MJ、SD、DALL-E）</li><li>输入描述文字</li><li>点击生成即可</li></ol>'
  },
  {
    id: 3,
    question: '为什么生成失败？',
    answer: '<p>生成失败可能的原因：</p><ul><li>算力不足</li><li>内容违规</li><li>网络不稳定</li><li>服务器繁忙</li></ul><p>请检查算力余额并重试。</p>'
  },
  {
    id: 4,
    question: '如何成为VIP会员？',
    answer: '<p>成为VIP会员的方式：</p><ol><li>进入会员中心</li><li>选择合适的套餐</li><li>完成支付</li><li>自动开通VIP权限</li></ol>'
  },
  {
    id: 5,
    question: '如何导出聊天记录？',
    answer: '<p>导出聊天记录步骤：</p><ol><li>进入对话页面</li><li>点击右上角菜单</li><li>选择"导出记录"</li><li>选择导出格式</li><li>确认导出</li></ol>'
  }
])

// 功能指南
const guides = ref([
  {
    id: 1,
    title: 'AI对话',
    desc: '与AI智能对话',
    icon: 'icon-chat',
    color: '#1989fa',
    content: 'AI对话使用指南详细内容...'
  },
  {
    id: 2,
    title: 'AI绘画',
    desc: '生成精美图片',
    icon: 'icon-mj',
    color: '#8B5CF6',
    content: 'AI绘画使用指南详细内容...'
  },
  {
    id: 3,
    title: 'AI音乐',
    desc: '创作美妙音乐',
    icon: 'icon-music',
    color: '#ee0a24',
    content: 'AI音乐创作指南详细内容...'
  },
  {
    id: 4,
    title: 'AI视频',
    desc: '制作精彩视频',
    icon: 'icon-video',
    color: '#07c160',
    content: 'AI视频制作指南详细内容...'
  }
])

// 问题分类
const categories = ref([
  { id: 1, name: '账户问题', icon: 'icon-user', count: 15 },
  { id: 2, name: '功能使用', icon: 'icon-apps', count: 23 },
  { id: 3, name: '充值支付', icon: 'icon-money', count: 12 },
  { id: 4, name: '技术问题', icon: 'icon-setting', count: 18 },
  { id: 5, name: '其他问题', icon: 'icon-help', count: 8 }
])

// 使用提示
const tips = ref([
  {
    id: 1,
    title: '提高绘画质量',
    content: '使用详细的描述词可以获得更好的绘画效果，建议加入风格、色彩、构图等关键词。',
    icon: 'icon-bulb'
  },
  {
    id: 2,
    title: '节省算力',
    content: '合理使用不同模型，简单问题使用GPT-3.5，复杂任务使用GPT-4。',
    icon: 'icon-flash'
  },
  {
    id: 3,
    title: '快速上手',
    content: '查看应用中心的预设角色，可以快速体验不同类型的AI对话。',
    icon: 'icon-star'
  }
])

onMounted(() => {
  // 初始化客服消息
  customerMessages.value = [
    {
      id: 1,
      content: '您好！欢迎使用我们的AI创作平台，有什么可以帮助您的吗？',
      isUser: false,
      time: new Date()
    }
  ]
})

// 搜索
const onSearch = (keyword) => {
  if (!keyword.trim()) {
    searchResults.value = []
    return
  }

  // 模拟搜索结果
  const allContent = [
    ...frequentFAQs.value.map(faq => ({
      id: faq.id,
      title: faq.question,
      content: faq.answer,
      type: 'faq'
    })),
    ...guides.value.map(guide => ({
      id: guide.id,
      title: guide.title,
      content: guide.content,
      type: 'guide'
    }))
  ]

  searchResults.value = allContent
    .filter(item => 
      item.title.includes(keyword) || item.content.includes(keyword)
    )
    .map(item => ({
      ...item,
      snippet: getSearchSnippet(item.content, keyword)
    }))
}

// 获取搜索摘要
const getSearchSnippet = (content, keyword) => {
  const cleanContent = content.replace(/<[^>]*>/g, '')
  const index = cleanContent.toLowerCase().indexOf(keyword.toLowerCase())
  if (index === -1) return cleanContent.substr(0, 100) + '...'
  
  const start = Math.max(0, index - 50)
  const end = Math.min(cleanContent.length, index + keyword.length + 50)
  let snippet = cleanContent.substr(start, end - start)
  
  // 高亮关键词
  const regex = new RegExp(`(${keyword})`, 'gi')
  snippet = snippet.replace(regex, '<mark>$1</mark>')
  
  return (start > 0 ? '...' : '') + snippet + (end < cleanContent.length ? '...' : '')
}

// 打开指南
const openGuide = (guide) => {
  selectedHelp.value = {
    title: guide.title,
    content: guide.content || '<p>该指南内容正在完善中，敬请期待。</p>'
  }
  showHelpDetail.value = true
}

// 打开分类
const openCategory = (category) => {
  showNotify({ type: 'primary', message: `正在加载${category.name}...` })
  // 这里可以跳转到分类详情页
}

// 打开搜索结果
const openSearchResult = (result) => {
  selectedHelp.value = {
    title: result.title,
    content: result.content
  }
  showHelpDetail.value = true
}

// 点赞帮助
const likeHelp = (help) => {
  showSuccessToast('感谢您的反馈！')
}

// 分享帮助
const shareHelp = (help) => {
  if (navigator.share) {
    navigator.share({
      title: help.title,
      text: help.content.replace(/<[^>]*>/g, ''),
      url: window.location.href
    })
  } else {
    showNotify({ type: 'primary', message: '该功能暂不支持' })
  }
}

// 打开客服
const openCustomerService = () => {
  showCustomerChat.value = true
}

// 发送客服消息
const sendCustomerMessage = () => {
  if (!customerMessage.value.trim()) return

  // 添加用户消息
  customerMessages.value.push({
    id: Date.now(),
    content: customerMessage.value,
    isUser: true,
    time: new Date()
  })

  const userMessage = customerMessage.value
  customerMessage.value = ''

  // 滚动到底部
  nextTick(() => {
    if (chatMessages.value) {
      chatMessages.value.scrollTop = chatMessages.value.scrollHeight
    }
  })

  // 模拟客服回复
  setTimeout(() => {
    let reply = '感谢您的问题，我们会尽快为您处理。'
    
    if (userMessage.includes('算力')) {
      reply = '关于算力问题，您可以在会员中心购买算力套餐，或者通过邀请好友获得免费算力。'
    } else if (userMessage.includes('绘画')) {
      reply = '关于AI绘画，建议您使用详细的描述词，这样可以获得更好的效果。'
    } else if (userMessage.includes('充值')) {
      reply = '充值问题请您检查支付方式是否正确，如有问题可以联系技术客服。'
    }

    customerMessages.value.push({
      id: Date.now(),
      content: reply,
      isUser: false,
      time: new Date()
    })

    nextTick(() => {
      if (chatMessages.value) {
        chatMessages.value.scrollTop = chatMessages.value.scrollHeight
      }
    })
  }, 1000)
}

// 加入QQ群
const joinQQGroup = () => {
  // 尝试打开QQ群链接
  const qqGroupUrl = 'mqqopensdkapi://bizAgent/qm/qr?url=http%3A%2F%2Fqm.qq.com%2Fcgi-bin%2Fqm%2Fqr%3Ffrom%3Dapp%26p%3Dandroid%26k%3D123456789'
  window.location.href = qqGroupUrl
  
  setTimeout(() => {
    showNotify({ type: 'primary', message: '请在QQ中搜索群号：123456789' })
  }, 1000)
}

// 格式化时间
const formatTime = (time) => {
  return time.toLocaleTimeString('zh-CN', { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}

// 二维码加载错误
const onQRError = (e) => {
  e.target.src = '/images/default-qr.png'
}
</script>

<style lang="scss" scoped>
.help-page {
  min-height: 100vh;
  background: var(--van-background);
  
  .help-content {
    padding: 54px 16px 20px;
    
    .search-section {
      margin-bottom: 20px;
    }
    
    .section-title {
      font-size: 18px;
      font-weight: 600;
      color: var(--van-text-color);
      margin: 0 0 16px 4px;
    }
    
    .faq-section {
      margin-bottom: 24px;
      
      :deep(.van-collapse-item) {
        background: var(--van-cell-background);
        border-radius: 12px;
        margin-bottom: 8px;
        overflow: hidden;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
        
        .van-collapse-item__title {
          padding: 16px;
          font-weight: 500;
        }
        
        .van-collapse-item__content {
          padding: 0 16px 16px;
          
          .faq-answer {
            color: var(--van-gray-7);
            line-height: 1.6;
            
            :deep(ul), :deep(ol) {
              padding-left: 20px;
              margin: 8px 0;
            }
            
            :deep(li) {
              margin: 4px 0;
            }
            
            :deep(p) {
              margin: 8px 0;
            }
          }
        }
      }
    }
    
    .guide-section {
      margin-bottom: 24px;
      
      .guide-item {
        .guide-card {
          background: var(--van-cell-background);
          border-radius: 12px;
          padding: 20px 16px;
          text-align: center;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
          cursor: pointer;
          transition: all 0.3s ease;
          
          &:active {
            transform: scale(0.98);
          }
          
          .guide-icon {
            width: 50px;
            height: 50px;
            border-radius: 12px;
            display: flex;
            align-items: center;
            justify-content: center;
            margin: 0 auto 12px;
            
            .iconfont {
              font-size: 24px;
              color: white;
            }
          }
          
          .guide-title {
            font-size: 16px;
            font-weight: 600;
            color: var(--van-text-color);
            margin-bottom: 8px;
          }
          
          .guide-desc {
            font-size: 13px;
            color: var(--van-gray-6);
          }
        }
      }
    }
    
    .category-section,
    .contact-section {
      margin-bottom: 24px;
      
      :deep(.van-cell-group) {
        border-radius: 12px;
        overflow: hidden;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
        
        .van-cell {
          padding: 16px;
          
          .category-icon {
            font-size: 18px;
            color: var(--van-primary-color);
            margin-right: 12px;
          }
          
          .van-cell__title {
            font-size: 15px;
            font-weight: 500;
          }
          
          .online-status {
            color: #07c160;
            font-size: 12px;
          }
          
          .qq-number {
            color: var(--van-gray-6);
            font-size: 13px;
          }
        }
      }
    }
    
    .search-results {
      .search-snippet {
        margin-top: 4px;
        color: var(--van-gray-6);
        font-size: 13px;
        line-height: 1.4;
        
        :deep(mark) {
          background: var(--van-primary-color);
          color: white;
          padding: 1px 2px;
          border-radius: 2px;
        }
      }
    }
    
    .tips-section {
      .tips-swipe {
        height: 140px;
        border-radius: 12px;
        overflow: hidden;
        
        .tip-card {
          background: linear-gradient(135deg, var(--van-primary-color), #8B5CF6);
          color: white;
          padding: 20px;
          text-align: center;
          height: 100%;
          display: flex;
          flex-direction: column;
          justify-content: center;
          
          .tip-icon {
            margin-bottom: 12px;
            
            .iconfont {
              font-size: 28px;
              opacity: 0.9;
            }
          }
          
          .tip-title {
            font-size: 16px;
            font-weight: 600;
            margin: 0 0 8px 0;
          }
          
          .tip-content {
            font-size: 13px;
            opacity: 0.9;
            line-height: 1.4;
            margin: 0;
          }
        }
      }
    }
  }
  
  .help-detail {
    padding: 20px;
    max-height: 60vh;
    overflow-y: auto;
    
    .detail-content {
      color: var(--van-text-color);
      line-height: 1.6;
      margin-bottom: 20px;
      
      :deep(p) {
        margin: 8px 0;
      }
      
      :deep(ul), :deep(ol) {
        padding-left: 20px;
        margin: 8px 0;
      }
    }
    
    .detail-actions {
      display: flex;
      gap: 12px;
      
      .van-button {
        flex: 1;
      }
    }
  }
  
  .wechat-qr {
    text-align: center;
    padding: 20px;
    
    .qr-code {
      width: 200px;
      height: 200px;
      margin: 0 auto 16px;
      border: 1px solid var(--van-border-color);
      border-radius: 8px;
      overflow: hidden;
      
      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }
    
    .qr-tip {
      font-size: 14px;
      color: var(--van-gray-6);
      margin: 0;
    }
  }
  
  .customer-chat {
    height: 500px;
    display: flex;
    flex-direction: column;
    
    .chat-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 16px;
      border-bottom: 1px solid var(--van-border-color);
      
      .customer-info {
        display: flex;
        align-items: center;
        
        .customer-detail {
          margin-left: 12px;
          
          .customer-name {
            font-size: 15px;
            font-weight: 500;
            color: var(--van-text-color);
          }
          
          .customer-status {
            font-size: 12px;
            
            &.online {
              color: #07c160;
            }
          }
        }
      }
    }
    
    .chat-messages {
      flex: 1;
      overflow-y: auto;
      padding: 16px;
      
      .message-item {
        margin-bottom: 16px;
        
        &.user-message {
          text-align: right;
          
          .message-content {
            background: var(--van-primary-color);
            color: white;
          }
        }
        
        .message-content {
          display: inline-block;
          max-width: 80%;
          padding: 10px 12px;
          background: var(--van-gray-1);
          border-radius: 8px;
          font-size: 14px;
          line-height: 1.4;
        }
        
        .message-time {
          font-size: 11px;
          color: var(--van-gray-5);
          margin-top: 4px;
        }
      }
    }
    
    .chat-input {
      padding: 16px;
      border-top: 1px solid var(--van-border-color);
    }
  }
}

// 深色主题优化
:deep(.van-theme-dark) {
  .help-page {
    .van-collapse-item,
    .guide-card,
    .van-cell-group {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    }
  }
}
</style>