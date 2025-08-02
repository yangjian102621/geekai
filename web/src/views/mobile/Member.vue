<template>
  <div class="member-center">
    <van-nav-bar title="会员中心" left-arrow @click-left="router.back()" fixed>
      <template #right>
        <van-icon name="question-o" @click="showHelp = true" />
      </template>
    </van-nav-bar>

    <div class="member-content">
      <!-- 用户信息卡片 -->
      <div class="user-card" v-if="isLogin">
        <div class="user-info">
          <van-image :src="userInfo.avatar" round width="60" height="60" />
          <div class="user-detail">
            <div class="user-name">{{ userInfo.nickname || userInfo.username }}</div>
            <div class="user-level">
              <van-tag :type="vipInfo.isVip ? 'primary' : 'default'" size="medium">
                {{ vipInfo.isVip ? 'VIP会员' : '普通用户' }}
              </van-tag>
            </div>
          </div>
        </div>
        
        <div class="user-stats">
          <div class="stat-item">
            <div class="stat-value">{{ userInfo.power || 0 }}</div>
            <div class="stat-label">剩余算力</div>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item" v-if="vipInfo.isVip">
            <div class="stat-value">{{ vipInfo.daysLeft }}天</div>
            <div class="stat-label">VIP剩余</div>
          </div>
          <div class="stat-item" v-else>
            <div class="stat-value">--</div>
            <div class="stat-label">VIP到期</div>
          </div>
        </div>
      </div>

      <!-- 会员特权 -->
      <div class="privileges-section">
        <h3 class="section-title">会员特权</h3>
        <div class="privileges-grid">
          <div 
            v-for="privilege in privileges" 
            :key="privilege.key"
            class="privilege-item"
            :class="{ active: vipInfo.isVip }"
          >
            <div class="privilege-icon" :style="{ backgroundColor: privilege.color }">
              <i class="iconfont" :class="privilege.icon"></i>
            </div>
            <div class="privilege-info">
              <div class="privilege-title">{{ privilege.title }}</div>
              <div class="privilege-desc">{{ privilege.desc }}</div>
            </div>
            <div class="privilege-status">
              <van-icon 
                :name="vipInfo.isVip ? 'success' : 'cross'" 
                :color="vipInfo.isVip ? '#07c160' : '#ee0a24'"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- 充值套餐 -->
      <div class="packages-section">
        <h3 class="section-title">充值套餐</h3>
        <div class="packages-list">
          <div 
            v-for="pkg in packages" 
            :key="pkg.id"
            class="package-item"
            :class="{ recommended: pkg.recommended }"
            @click="selectPackage(pkg)"
          >
            <div class="package-tag" v-if="pkg.recommended">推荐</div>
            <div class="package-header">
              <div class="package-name">{{ pkg.name }}</div>
              <div class="package-price">
                <span class="current-price">￥{{ pkg.discount || pkg.price }}</span>
                <span class="original-price" v-if="pkg.discount">￥{{ pkg.price }}</span>
              </div>
            </div>
            <div class="package-features">
              <div class="feature-item">
                <van-icon name="checked" color="#07c160" />
                <span>{{ pkg.power }}算力值</span>
              </div>
              <div class="feature-item" v-if="pkg.days > 0">
                <van-icon name="checked" color="#07c160" />
                <span>{{ pkg.days }}天有效期</span>
              </div>
              <div class="feature-item" v-else>
                <van-icon name="checked" color="#07c160" />
                <span>长期有效</span>
              </div>
              <div class="feature-item" v-if="pkg.features">
                <van-icon name="checked" color="#07c160" />
                <span>{{ pkg.features }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 支付方式 -->
      <div class="payment-section" v-if="selectedPackage">
        <h3 class="section-title">支付方式</h3>
        <van-radio-group v-model="selectedPayment">
          <van-cell-group inset>
            <van-cell 
              v-for="payment in paymentMethods" 
              :key="payment.type"
              clickable 
              @click="selectedPayment = payment.type"
            >
              <template #title>
                <div class="payment-info">
                  <i class="iconfont" :class="payment.icon" :style="{ color: payment.color }"></i>
                  <span>{{ payment.name }}</span>
                </div>
              </template>
              <template #right-icon>
                <van-radio :name="payment.type" />
              </template>
            </van-cell>
          </van-radio-group>
        </van-radio-group>
      </div>

      <!-- 支付按钮 -->
      <div class="pay-button" v-if="selectedPackage">
        <van-button 
          type="primary" 
          size="large" 
          round 
          block
          :loading="payLoading"
          @click="processPay"
        >
          立即支付 ￥{{ selectedPackage.discount || selectedPackage.price }}
        </van-button>
      </div>
    </div>

    <!-- 帮助弹窗 -->
    <van-action-sheet v-model:show="showHelp" title="会员帮助">
      <div class="help-content">
        <div class="help-item">
          <h4>什么是算力？</h4>
          <p>算力是使用AI功能时消耗的虚拟货币，不同功能消耗的算力不同。</p>
        </div>
        <div class="help-item">
          <h4>如何获得算力？</h4>
          <p>通过充值套餐可以获得算力，会员用户还可享受每月赠送的算力。</p>
        </div>
        <div class="help-item">
          <h4>VIP特权说明</h4>
          <p>VIP会员享有更多功能权限、优先处理、专属客服等特权服务。</p>
        </div>
      </div>
    </van-action-sheet>
  </div>
</template>

<script setup>
import { checkSession } from '@/store/cache'
import { httpGet, httpPost } from '@/utils/http'
import { showLoginDialog } from '@/utils/libs'
import { showFailToast, showLoadingToast, showNotify } from 'vant'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const isLogin = ref(false)
const userInfo = ref({})
const packages = ref([])
const paymentMethods = ref([])
const selectedPackage = ref(null)
const selectedPayment = ref('alipay')
const payLoading = ref(false)
const showHelp = ref(false)

// VIP信息计算
const vipInfo = computed(() => {
  const now = Date.now()
  const expiredTime = userInfo.value.expired_time ? userInfo.value.expired_time * 1000 : 0
  const isVip = expiredTime > now
  const daysLeft = isVip ? Math.ceil((expiredTime - now) / (24 * 60 * 60 * 1000)) : 0
  
  return { isVip, daysLeft }
})

// 会员特权配置
const privileges = ref([
  {
    key: 'unlimited',
    title: '无限对话',
    desc: '不限制对话次数',
    icon: 'icon-chat',
    color: '#07c160'
  },
  {
    key: 'priority',
    title: '优先处理',
    desc: '请求优先处理',
    icon: 'icon-flash',
    color: '#ff9500'
  },
  {
    key: 'models',
    title: '高级模型',
    desc: '使用最新AI模型',
    icon: 'icon-star',
    color: '#ffd700'
  },
  {
    key: 'support',
    title: '专属客服',
    desc: '7×24小时客服支持',
    icon: 'icon-service',
    color: '#1989fa'
  }
])

onMounted(() => {
  initPage()
})

const initPage = async () => {
  try {
    const user = await checkSession()
    isLogin.value = true
    userInfo.value = user
    
    // 获取用户详细信息
    const profileRes = await httpGet('/api/user/profile')
    userInfo.value = { ...userInfo.value, ...profileRes.data }
    
  } catch (error) {
    showLoginDialog(router)
    return
  }

  // 获取充值套餐
  fetchPackages()
  // 获取支付方式
  fetchPaymentMethods()
}

const fetchPackages = () => {
  httpGet('/api/product/list')
    .then((res) => {
      // 添加推荐标签和特权描述
      packages.value = res.data.map((pkg, index) => ({
        ...pkg,
        recommended: index === 1, // 第二个套餐设为推荐
        features: pkg.days > 30 ? 'VIP专属权益' : null
      }))
    })
    .catch((e) => {
      showFailToast('获取套餐失败：' + e.message)
    })
}

const fetchPaymentMethods = () => {
  httpGet('/api/payment/payWays')
    .then((res) => {
      paymentMethods.value = res.data.map(item => ({
        type: item.pay_type,
        name: item.pay_type === 'alipay' ? '支付宝' : '微信支付',
        icon: item.pay_type === 'alipay' ? 'icon-alipay' : 'icon-wechat-pay',
        color: item.pay_type === 'alipay' ? '#1677ff' : '#07c160',
        payWay: item.pay_way
      }))
      
      if (paymentMethods.value.length > 0) {
        selectedPayment.value = paymentMethods.value[0].type
      }
    })
    .catch((e) => {
      showFailToast('获取支付方式失败：' + e.message)
    })
}

const selectPackage = (pkg) => {
  selectedPackage.value = pkg
}

const processPay = () => {
  if (!selectedPackage.value || !selectedPayment.value) {
    showNotify({ type: 'warning', message: '请选择套餐和支付方式' })
    return
  }

  const paymentMethod = paymentMethods.value.find(p => p.type === selectedPayment.value)
  if (!paymentMethod) {
    showNotify({ type: 'danger', message: '支付方式无效' })
    return
  }

  payLoading.value = true
  showLoadingToast({
    message: '正在创建订单...',
    forbidClick: true,
  })

  const host = `${location.protocol}//${location.host}`
  
  httpPost('/api/payment/doPay', {
    product_id: selectedPackage.value.id,
    pay_way: paymentMethod.payWay,
    pay_type: paymentMethod.type,
    user_id: userInfo.value.id,
    host: host,
    device: 'mobile',
  })
    .then((res) => {
      location.href = res.data
    })
    .catch((e) => {
      payLoading.value = false
      showFailToast('创建订单失败：' + e.message)
    })
}
</script>

<style lang="scss" scoped>
.member-center {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  
  .member-content {
    padding: 54px 16px 20px;
    
    .user-card {
      background: rgba(255, 255, 255, 0.95);
      border-radius: 16px;
      padding: 20px;
      margin-bottom: 20px;
      backdrop-filter: blur(10px);
      box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
      
      .user-info {
        display: flex;
        align-items: center;
        margin-bottom: 16px;
        
        .user-detail {
          margin-left: 16px;
          
          .user-name {
            font-size: 18px;
            font-weight: 600;
            color: #333;
            margin-bottom: 8px;
          }
        }
      }
      
      .user-stats {
        display: flex;
        align-items: center;
        justify-content: space-around;
        
        .stat-item {
          text-align: center;
          
          .stat-value {
            font-size: 20px;
            font-weight: 700;
            color: var(--van-primary-color);
            margin-bottom: 4px;
          }
          
          .stat-label {
            font-size: 12px;
            color: #666;
          }
        }
        
        .stat-divider {
          width: 1px;
          height: 30px;
          background: #e5e5e5;
        }
      }
    }
    
    .section-title {
      font-size: 18px;
      font-weight: 600;
      color: white;
      margin: 0 0 16px 4px;
      text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
    }
    
    .privileges-section {
      margin-bottom: 24px;
      
      .privileges-grid {
        background: rgba(255, 255, 255, 0.95);
        border-radius: 12px;
        padding: 16px;
        backdrop-filter: blur(10px);
        
        .privilege-item {
          display: flex;
          align-items: center;
          padding: 12px 0;
          border-bottom: 1px solid #f5f5f5;
          
          &:last-child {
            border-bottom: none;
          }
          
          &.active {
            .privilege-icon {
              box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
            }
          }
          
          .privilege-icon {
            width: 40px;
            height: 40px;
            border-radius: 10px;
            display: flex;
            align-items: center;
            justify-content: center;
            margin-right: 12px;
            transition: all 0.3s ease;
            
            .iconfont {
              font-size: 20px;
              color: white;
            }
          }
          
          .privilege-info {
            flex: 1;
            
            .privilege-title {
              font-size: 15px;
              font-weight: 600;
              color: #333;
              margin-bottom: 4px;
            }
            
            .privilege-desc {
              font-size: 13px;
              color: #666;
            }
          }
          
          .privilege-status {
            margin-left: 8px;
          }
        }
      }
    }
    
    .packages-section {
      margin-bottom: 24px;
      
      .packages-list {
        .package-item {
          background: rgba(255, 255, 255, 0.95);
          border-radius: 12px;
          padding: 16px;
          margin-bottom: 12px;
          backdrop-filter: blur(10px);
          position: relative;
          cursor: pointer;
          transition: all 0.3s ease;
          border: 2px solid transparent;
          
          &:active {
            transform: scale(0.98);
          }
          
          &.recommended {
            border-color: var(--van-primary-color);
            box-shadow: 0 8px 24px rgba(25, 137, 250, 0.3);
          }
          
          .package-tag {
            position: absolute;
            top: -1px;
            right: 16px;
            background: var(--van-primary-color);
            color: white;
            font-size: 12px;
            padding: 4px 12px;
            border-radius: 0 0 8px 8px;
          }
          
          .package-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 12px;
            
            .package-name {
              font-size: 16px;
              font-weight: 600;
              color: #333;
            }
            
            .package-price {
              text-align: right;
              
              .current-price {
                font-size: 18px;
                font-weight: 700;
                color: var(--van-primary-color);
              }
              
              .original-price {
                font-size: 14px;
                color: #999;
                text-decoration: line-through;
                margin-left: 8px;
              }
            }
          }
          
          .package-features {
            .feature-item {
              display: flex;
              align-items: center;
              margin-bottom: 8px;
              font-size: 14px;
              color: #666;
              
              &:last-child {
                margin-bottom: 0;
              }
              
              .van-icon {
                margin-right: 8px;
              }
            }
          }
        }
      }
    }
    
    .payment-section {
      margin-bottom: 24px;
      
      :deep(.van-cell-group) {
        border-radius: 12px;
        overflow: hidden;
        backdrop-filter: blur(10px);
        
        .van-cell {
          background: rgba(255, 255, 255, 0.95);
          
          .payment-info {
            display: flex;
            align-items: center;
            
            .iconfont {
              font-size: 20px;
              margin-right: 12px;
            }
          }
        }
      }
    }
    
    .pay-button {
      position: sticky;
      bottom: 20px;
      z-index: 100;
    }
  }
  
  .help-content {
    padding: 20px;
    
    .help-item {
      margin-bottom: 20px;
      
      &:last-child {
        margin-bottom: 0;
      }
      
      h4 {
        font-size: 16px;
        font-weight: 600;
        color: var(--van-text-color);
        margin: 0 0 8px 0;
      }
      
      p {
        font-size: 14px;
        color: var(--van-gray-6);
        line-height: 1.5;
        margin: 0;
      }
    }
  }
}
</style>