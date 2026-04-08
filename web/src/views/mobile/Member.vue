<template>
  <div class="member-page">
    <div class="member-content" v-loading="loading" :element-loading-text="loadingText">
      <!-- 产品套餐 -->
      <div class="products-section">
        <div class="text-center bg-[#7c3aed] text-white rounded-lg p-3 mb-4">充值套餐</div>

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
                <div class="discount-price">
                  <span class="label">商品价格：</span>
                  <span class="price">￥{{ item.price }}</span>
                </div>

                <div class="product-details">
                  <div class="detail-item">
                    <span class="label">算力值：</span>
                    <span class="value">{{ item.power }}</span>
                  </div>
                </div>
              </div>

              <div class="payment-methods">
                <div class="methods-grid">
                  <button class="payment-btn wechat-btn" @click="wxPay(item)">
                    <i class="iconfont icon-wechat-pay"></i>
                    <span>微信支付</span>
                  </button>
                  <button class="payment-btn alipay-btn" @click="alipay(item)">
                    <i class="iconfont icon-alipay"></i>
                    <span>支付宝</span>
                  </button>
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
        <van-cell-group>
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
      :show="showPayDialog"
      :title="title"
      :show-cancel-button="false"
      confirm-button-text="确认支付成功"
      :close-on-click-overlay="false"
      @confirm="paySuccess"
    >
      <div class="pay-dialog-content">
        <div v-if="qrImg" class="qr-section">
          <div class="qr-container">
            <van-image :src="qrImg" width="200" height="200" />
          </div>
          <p class="pay-amount">￥{{ currentPrice }}</p>
          <p class="pay-tip">支付成功之后点击确定按钮</p>
        </div>
      </div>
    </van-dialog>

    <!-- 卡密兑换弹窗 -->
    <van-dialog
      :show="showRedeemVerifyDialog"
      title="卡密兑换"
      :show-cancel-button="false"
      :show-confirm-button="false"
      width="90%"
      :close-on-click-overlay="true"
    >
      <div class="p-4">
        <redeem-verify @hide="redeemCallback" />
      </div>
    </van-dialog>
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
import { onMounted, onUnmounted, ref } from 'vue'

// 响应式数据
const list = ref([])
const vipImg = ref('/images/menu/member.png')
const isLogin = ref(false)
const loading = ref(true)
const loadingText = ref('加载中...')
const vipInfoText = ref('')
const userOrderKey = ref(0)

// 弹窗控制
const showRedeemVerifyDialog = ref(false)
const showPayDialog = ref(false)

// 支付相关
const qrImg = ref('')
const currentPrice = ref(0)
const currentProduct = ref(null)
const selectedPid = ref(0)
const orderTimeout = ref(1800)
const handler = ref(null)
const title = ref('')

const store = useSharedStore()

// 初始化
onMounted(() => {
  checkSession()
    .then(() => {
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
      if (res.data['order_pay_timeout'] > 0) {
        orderTimeout.value = res.data['order_pay_timeout']
      }
    })
    .catch((e) => {
      console.error('获取系统配置失败：', e.message)
    })
})

// 支付处理
const wxPay = (product) => {
  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return
  }

  selectedPid.value = product.id
  currentProduct.value = product
  currentPrice.value = Number(product.price)
  title.value = '请打开微信扫码支付'

  showLoadingToast({
    message: '正在生成支付订单...',
    forbidClick: true,
  })

  // 生成支付订单
  GenerateOrder('wxpay')
}

const alipay = (product) => {
  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return
  }

  selectedPid.value = product.id
  currentProduct.value = product
  currentPrice.value = Number(product.price)
  title.value = '请打开支付宝扫码支付'

  showLoadingToast({
    message: '正在生成支付订单...',
    forbidClick: true,
  })

  // 生成支付订单
  GenerateOrder('alipay')
}

function GenerateOrder(payWay) {
  // 生成支付订单
  httpPost('/api/payment/create', {
    pid: selectedPid.value,
    pay_way: payWay,
    domain: `${window.location.protocol}//${window.location.host}`,
    device: 'pc',
  })
    .then((res) => {
      if (res.data.pay_url) {
        // 生成二维码
        QRCode.toDataURL(res.data.pay_url, { width: 200, height: 200, margin: 2 }, (error, url) => {
          if (error) {
            showFailToast('生成二维码失败')
          } else {
            qrImg.value = url
            showPayDialog.value = true
            // 开始查询订单状态
            if (handler.value) {
              clearTimeout(handler.value)
            }
            handler.value = setTimeout(() => queryOrder(res.data.order_no), 3000)
          }
        })
      } else {
        showFailToast('支付链接生成失败')
      }
    })
    .catch((e) => {
      showFailToast('生成支付订单失败：' + e.message)
    })
}

// 查询订单状态
const queryOrder = async (orderNo) => {
  try {
    const res = await httpGet('/api/order/query?order_no=' + orderNo)
    if (res?.data.status === 2) {
      paySuccess(true)
    } else {
      // 继续查询，但设置最大查询次数
      const maxQueries = Math.floor(orderTimeout.value / 3) // 每3秒查询一次
      if (handler.value && maxQueries > 0) {
        handler.value = setTimeout(() => queryOrder(orderNo), 3000)
      } else {
        // 查询超时
        showFailToast('支付超时，请重新发起支付')
        showPayDialog.value = false
        qrImg.value = ''
      }
    }
  } catch (error) {
    console.error('查询订单状态失败:', error)
    // 继续查询，但设置最大重试次数
    if (handler.value) {
      handler.value = setTimeout(() => queryOrder(orderNo), 3000)
    }
  }
}

const paySuccess = () => {
  showPayDialog.value = false
  showSuccessToast('支付成功！')
  clearTimeout(handler.value)
  userOrderKey.value += 1
}

// 组件卸载时清理定时器
onUnmounted(() => {
  if (handler.value) {
    clearTimeout(handler.value)
    handler.value = null
  }
})
</script>

<style lang="scss" scoped>
.member-page {
  min-height: 100vh;
  background: var(--van-background);

  .member-content {
    padding: 20px 16px;

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
                gap: 12px;

                .payment-btn {
                  height: 44px;
                  border: none;
                  border-radius: 8px;
                  font-size: 14px;
                  font-weight: 500;
                  cursor: pointer;
                  transition: all 0.3s ease;
                  display: flex;
                  align-items: center;
                  justify-content: center;
                  gap: 8px;

                  &:hover {
                    transform: translateY(-2px);
                    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
                  }

                  .iconfont {
                    font-size: 18px;
                  }

                  &.wechat-btn {
                    background: #07c160;
                    color: white;

                    &:hover {
                      background: #06ad56;
                    }
                  }

                  &.alipay-btn {
                    background: #15a6e8;
                    color: white;

                    &:hover {
                      background: #1395d1;
                    }
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
        border-radius: 12px;
        overflow: hidden;
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
        font-weight: 500;
      }

      .qr-container {
        margin-bottom: 12px;
        padding: 16px;
        background: #f8f9fa;
        border-radius: 12px;
        display: inline-block;
      }

      .pay-amount {
        font-size: 20px;
        font-weight: 600;
        color: #ff6b35;
        margin-bottom: 16px;
      }

      .pay-status {
        margin-top: 15px;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        padding: 12px;
        background: #f0f9ff;
        border-radius: 8px;
        border: 1px solid #e0f2fe;

        .success-status {
          display: flex;
          align-items: center;
          gap: 8px;
          color: #07c160;
          font-weight: 500;
        }
      }
    }

    .pay-actions {
      display: flex;
      gap: 12px;
      justify-content: center;
      padding-top: 16px;
      border-top: 1px solid #f0f0f0;
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
