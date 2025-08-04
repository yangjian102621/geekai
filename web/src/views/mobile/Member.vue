<template>
  <div class="member-page">
    <div class="member-content" v-loading="loading" :element-loading-text="loadingText">
      <!-- 用户信息卡片 -->
      <div class="user-card" v-if="isLogin">
        <div class="user-header">
          <div class="user-avatar">
            <van-image :src="userAvatar" round width="60" height="60" />
          </div>
          <div class="user-info">
            <h3 class="username">{{ userInfo.nickname || userInfo.username }}</h3>
            <div class="user-meta">
              <van-tag type="primary" v-if="isVip">VIP会员</van-tag>
              <van-tag type="default" v-else>普通用户</van-tag>
              <span class="user-id">ID: {{ userInfo.id }}</span>
            </div>
          </div>
        </div>
        <div class="user-stats">
          <div class="stat-item">
            <div class="stat-value">{{ userInfo.power || 0 }}</div>
            <div class="stat-label">剩余算力</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ vipDays }}</div>
            <div class="stat-label">VIP天数</div>
          </div>
        </div>
      </div>

      <!-- 产品套餐 -->
      <div class="products-section">
        <h3 class="section-title">会员套餐</h3>
        <div class="info-alert" v-if="vipInfoText">
          <van-notice-bar
            :text="vipInfoText"
            color="#1989fa"
            background="#ecf9ff"
            :scrollable="false"
          />
        </div>

        <div class="products-grid" v-if="list.length > 0">
          <div v-for="item in list" :key="item.id" class="product-card">
            <div class="product-header">
              <div class="product-image">
                <van-image :src="vipImg" fit="cover" />
              </div>
              <div class="product-title">{{ item.name }}</div>
            </div>

            <div class="product-info">
              <div class="price-info">
                <div class="original-price">
                  <span class="label">原价：</span>
                  <span class="price">￥{{ item.price }}</span>
                </div>
                <div class="discount-price">
                  <span class="label">优惠价：</span>
                  <span class="price">￥{{ item.discount }}</span>
                </div>
              </div>

              <div class="product-details">
                <div class="detail-item">
                  <span class="label">有效期：</span>
                  <span class="value">{{ item.days > 0 ? item.days + '天' : '长期有效' }}</span>
                </div>
                <div class="detail-item">
                  <span class="label">算力值：</span>
                  <span class="value">{{ item.power }}</span>
                </div>
              </div>

              <div class="payment-methods">
                <div class="methods-title">支付方式：</div>
                <div class="methods-grid">
                  <van-button
                    v-for="payWay in payWays"
                    :key="payWay.pay_type"
                    size="small"
                    :color="getPayButtonColor(payWay.pay_type)"
                    @click="pay(item, payWay)"
                    class="pay-button"
                  >
                    <i class="iconfont" :class="getPayIcon(payWay.pay_type)"></i>
                    {{ getPayButtonText(payWay.pay_type) }}
                  </van-button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <van-empty v-else description="暂无套餐" />
      </div>

      <!-- 卡密兑换 -->
      <div class="redeem-section" v-if="isLogin">
        <h3 class="section-title">卡密兑换</h3>
        <van-cell-group inset>
          <van-cell title="卡密兑换" is-link @click="showRedeemVerifyDialog = true">
            <template #icon>
              <i class="iconfont icon-redeem menu-icon"></i>
            </template>
          </van-cell>
        </van-cell-group>
      </div>

      <!-- 消费账单 -->
      <div class="bills-section" v-if="isLogin">
        <h3 class="section-title">消费账单</h3>
        <div class="bills-content">
          <user-order :key="userOrderKey" />
        </div>
      </div>
    </div>

    <!-- 支付弹窗 -->
    <van-dialog
      v-model="showPayDialog"
      title="支付确认"
      :show-cancel-button="false"
      :close-on-click-overlay="false"
    >
      <div class="pay-dialog-content">
        <div v-if="qrImg" class="qr-section">
          <p class="pay-tip">请使用微信扫码支付：</p>
          <div class="qr-container">
            <van-image :src="qrImg" width="200" height="200" />
          </div>
          <p class="pay-amount">￥{{ currentPrice }}</p>
        </div>
        <div class="pay-actions">
          <van-button type="success" @click="payCallback(true)">支付成功</van-button>
          <van-button type="danger" @click="payCallback(false)">支付失败</van-button>
        </div>
      </div>
    </van-dialog>

    <!-- 组件弹窗 -->
    <redeem-verify v-if="isLogin" :show="showRedeemVerifyDialog" @hide="redeemCallback" />
  </div>
</template>

<script setup>
import RedeemVerify from '@/components/RedeemVerify.vue'
import UserOrder from '@/components/UserOrder.vue'
import { checkSession, getSystemInfo } from '@/store/cache'
import { useSharedStore } from '@/store/sharedata'
import { httpGet, httpPost } from '@/utils/http'
import QRCode from 'qrcode'
import { showFailToast, showLoadingToast, showSuccessToast } from 'vant'
import { computed, onMounted, ref } from 'vue'

// 响应式数据
const list = ref([])
const vipImg = ref('/images/menu/member.png')
const userInfo = ref({})
const isLogin = ref(false)
const loading = ref(true)
const loadingText = ref('加载中...')
const vipInfoText = ref('')
const payWays = ref([])
const userOrderKey = ref(0)

// 弹窗控制
const showRedeemVerifyDialog = ref(false)
const showPayDialog = ref(false)

// 支付相关
const qrImg = ref('')
const currentPrice = ref(0)
const currentProduct = ref(null)
const currentPayWay = ref(null)

const store = useSharedStore()

// 计算属性
const isVip = computed(() => {
  const now = Date.now()
  const expiredTime = userInfo.value.expired_time ? userInfo.value.expired_time * 1000 : 0
  return expiredTime > now
})

const vipDays = computed(() => {
  if (!isVip.value) return 0
  const now = Date.now()
  const expiredTime = userInfo.value.expired_time * 1000
  return Math.ceil((expiredTime - now) / (24 * 60 * 60 * 1000))
})

const userAvatar = computed(() => {
  return userInfo.value.avatar || '/images/avatar/default.jpg'
})

// 支付按钮颜色
const getPayButtonColor = (payType) => {
  const colors = {
    alipay: '#15A6E8',
    wxpay: '#07C160',
    qqpay: '#12B7F5',
    paypal: '#0070BA',
    jdpay: '#E1251B',
    douyin: '#000000',
  }
  return colors[payType] || '#1989fa'
}

// 支付按钮图标
const getPayIcon = (payType) => {
  const icons = {
    alipay: 'icon-alipay',
    wxpay: 'icon-wechat-pay',
    qqpay: 'icon-qq',
    paypal: 'icon-paypal',
    jdpay: 'icon-jd-pay',
    douyin: 'icon-douyin',
  }
  return icons[payType] || 'icon-money'
}

// 支付按钮文本
const getPayButtonText = (payType) => {
  const texts = {
    alipay: '支付宝',
    wechat: '微信支付',
    qqpay: 'QQ钱包',
    paypal: 'PayPal',
    jdpay: '京东支付',
    douyin: '抖音支付',
  }
  return texts[payType] || '支付'
}

// 初始化
onMounted(() => {
  checkSession()
    .then((user) => {
      userInfo.value = user
      isLogin.value = true
    })
    .catch(() => {
      store.setShowLoginDialog(true)
    })

  // 获取产品列表
  httpGet('/api/product/list')
    .then((res) => {
      list.value = res.data
      loading.value = false
    })
    .catch((e) => {
      showFailToast('获取产品套餐失败：' + e.message)
      loading.value = false
    })

  // 获取系统配置
  getSystemInfo()
    .then((res) => {
      vipInfoText.value = res.data['vip_info_text']
    })
    .catch((e) => {
      console.error('获取系统配置失败：', e.message)
    })

  // 获取支付方式
  httpGet('/api/payment/payWays')
    .then((res) => {
      payWays.value = res.data
    })
    .catch((e) => {
      showFailToast('获取支付方式失败：' + e.message)
    })
})

// 支付处理
const pay = (product, payWay) => {
  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return
  }

  currentProduct.value = product
  currentPayWay.value = payWay
  currentPrice.value = Number(product.discount)

  showLoadingToast({
    message: '正在生成支付订单...',
    forbidClick: true,
  })

  let host = import.meta.env.VITE_API_HOST
  if (host === '') {
    host = `${location.protocol}//${location.host}`
  }

  httpPost(`${import.meta.env.VITE_API_HOST}/api/payment/doPay`, {
    product_id: product.id,
    pay_way: payWay.pay_way,
    pay_type: payWay.pay_type,
    user_id: userInfo.value.id,
    host: host,
    device: 'mobile',
  })
    .then((res) => {
      if (payWay.pay_way === 'wechat') {
        // 生成二维码
        QRCode.toDataURL(res.data, { width: 200, height: 200, margin: 2 }, (error, url) => {
          if (error) {
            showFailToast('生成二维码失败')
          } else {
            qrImg.value = url
            showPayDialog.value = true
          }
        })
      } else {
        // 跳转支付
        window.open(res.data, '_blank')
      }
    })
    .catch((e) => {
      showFailToast('生成支付订单失败：' + e.message)
    })
}

// 支付回调
const payCallback = (success) => {
  showPayDialog.value = false
  qrImg.value = ''

  if (success) {
    showSuccessToast('支付成功！')
    userOrderKey.value += 1
    // 刷新用户信息
    checkSession().then((user) => {
      userInfo.value = user
    })
  }
}

// 卡密兑换回调
const redeemCallback = () => {
  showRedeemVerifyDialog.value = false
  showSuccessToast('卡密兑换成功！')
  // 刷新用户信息
  checkSession().then((user) => {
    userInfo.value = user
  })
}
</script>

<style lang="scss" scoped>
.member-page {
  min-height: 100vh;
  background: var(--van-background);

  .member-content {
    padding: 20px 16px;

    .user-card {
      background: linear-gradient(135deg, var(--van-primary-color), #8b5cf6);
      border-radius: 16px;
      padding: 24px;
      margin-bottom: 24px;
      color: white;
      box-shadow: 0 8px 32px rgba(139, 92, 246, 0.3);

      .user-header {
        display: flex;
        align-items: center;
        margin-bottom: 20px;

        .user-avatar {
          margin-right: 16px;
        }

        .user-info {
          flex: 1;

          .username {
            font-size: 20px;
            font-weight: 600;
            margin: 0 0 8px 0;
          }

          .user-meta {
            display: flex;
            align-items: center;
            gap: 12px;

            .user-id {
              font-size: 12px;
              opacity: 0.8;
            }
          }
        }
      }

      .user-stats {
        display: flex;
        justify-content: space-around;

        .stat-item {
          text-align: center;

          .stat-value {
            font-size: 24px;
            font-weight: 700;
            margin-bottom: 4px;
          }

          .stat-label {
            font-size: 12px;
            opacity: 0.8;
          }
        }
      }
    }

    .section-title {
      font-size: 18px;
      font-weight: 600;
      color: var(--van-text-color);
      margin: 0 0 16px 4px;
    }

    .redeem-section {
      margin-bottom: 24px;

      :deep(.van-cell-group) {
        border-radius: 12px;
        overflow: hidden;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

        .van-cell {
          padding: 16px;

          .menu-icon {
            font-size: 18px;
            margin-right: 12px;
            color: var(--van-primary-color);
          }

          .van-cell__title {
            font-size: 15px;
            font-weight: 500;
          }
        }
      }
    }

    .products-section {
      margin-bottom: 24px;

      .info-alert {
        margin-bottom: 16px;
      }

      .products-grid {
        display: flex;
        flex-direction: column;
        gap: 16px;

        .product-card {
          background: var(--van-cell-background);
          border-radius: 16px;
          padding: 20px;
          box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
          transition: all 0.3s ease;

          &:hover {
            transform: translateY(-2px);
            box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
          }

          .product-header {
            display: flex;
            align-items: center;
            margin-bottom: 16px;

            .product-image {
              width: 48px;
              height: 48px;
              margin-right: 12px;
              border-radius: 8px;
              overflow: hidden;
            }

            .product-title {
              font-size: 18px;
              font-weight: 600;
              color: var(--van-text-color);
            }
          }

          .product-info {
            .price-info {
              margin-bottom: 16px;

              .original-price,
              .discount-price {
                display: flex;
                justify-content: space-between;
                align-items: center;
                margin-bottom: 8px;

                .label {
                  font-size: 14px;
                  color: var(--van-gray-6);
                }

                .price {
                  font-weight: 600;
                }
              }

              .original-price .price {
                color: var(--van-gray-5);
                text-decoration: line-through;
              }

              .discount-price .price {
                color: #ff6b35;
                font-size: 18px;
              }
            }

            .product-details {
              margin-bottom: 20px;

              .detail-item {
                display: flex;
                justify-content: space-between;
                align-items: center;
                margin-bottom: 8px;

                .label {
                  font-size: 14px;
                  color: var(--van-gray-6);
                }

                .value {
                  font-size: 14px;
                  color: var(--van-text-color);
                  font-weight: 500;
                }
              }
            }

            .payment-methods {
              .methods-title {
                font-size: 14px;
                color: var(--van-gray-6);
                margin-bottom: 12px;
              }

              .methods-grid {
                display: grid;
                grid-template-columns: repeat(2, 1fr);
                gap: 8px;

                .pay-button {
                  height: 36px;
                  font-size: 12px;

                  .iconfont {
                    margin-right: 4px;
                  }
                }
              }
            }
          }
        }
      }
    }

    .bills-section {
      .bills-content {
        background: var(--van-cell-background);
        border-radius: 12px;
        overflow: hidden;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
      }
    }
  }

  // 支付弹窗样式
  .pay-dialog-content {
    padding: 20px;

    .qr-section {
      text-align: center;
      margin-bottom: 20px;

      .pay-tip {
        font-size: 16px;
        color: var(--van-text-color);
        margin-bottom: 16px;
      }

      .qr-container {
        margin-bottom: 12px;
      }

      .pay-amount {
        font-size: 20px;
        font-weight: 600;
        color: #ff6b35;
      }
    }

    .pay-actions {
      display: flex;
      gap: 12px;
      justify-content: center;
    }
  }
}

// 深色主题优化
:deep(.van-theme-dark) {
  .member-page {
    .product-card {
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);

      &:hover {
        box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3);
      }
    }

    .van-cell-group {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    }
  }
}

// 响应式优化
@media (max-width: 375px) {
  .member-page {
    .member-content {
      padding: 16px 12px;

      .user-card {
        padding: 20px;

        .user-header .user-info .username {
          font-size: 18px;
        }

        .user-stats .stat-item .stat-value {
          font-size: 20px;
        }
      }

      .products-section .products-grid .product-card {
        padding: 16px;

        .product-header .product-title {
          font-size: 16px;
        }

        .payment-methods .methods-grid {
          grid-template-columns: 1fr;
        }
      }
    }
  }
}
</style>
