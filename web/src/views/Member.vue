<template>
  <div>
    <div
      class="member custom-scroll"
      v-loading="loading"
      element-loading-background="rgba(255,255,255,.3)"
      :element-loading-text="loadingText"
    >
      <div class="inner">
        <el-card class="profile-card">
          <el-row class="user-opt" :gutter="16">
            <el-col :span="24">
              <el-button class="profile-btn email" @click="showBindEmailDialog = true">
                <i class="iconfont icon-email"></i> 绑定邮箱
              </el-button>
            </el-col>
            <el-col :span="24">
              <el-button class="profile-btn mobile" @click="showBindMobileDialog = true">
                <i class="iconfont icon-mobile"></i> 绑定手机
              </el-button>
            </el-col>
            <el-col :span="24">
              <el-button class="profile-btn password" @click="showPasswordDialog = true">
                <i class="iconfont icon-password"></i> 修改密码
              </el-button>
            </el-col>
            <el-divider />
            <el-col :span="24">
              <el-button class="profile-btn redeem" @click="showRedeemVerifyDialog = true">
                <i class="iconfont icon-redeem"></i> 卡密兑换
              </el-button>
            </el-col>
          </el-row>
        </el-card>
        <div class="profile-bg"></div>

        <div class="product-box">
          <el-row v-if="list.length > 0" :gutter="24" class="list-box">
            <el-col
              v-for="item in list"
              :key="item"
              :xs="24"
              :sm="12"
              :md="8"
              :lg="6"
              class="product-col"
            >
              <div class="product-item">
                <div class="product-header">
                  <div class="image-container">
                    <el-image :src="vipImg" fit="cover" />
                    <div class="image-overlay">
                      <div class="vip-badge">热销</div>
                    </div>
                  </div>
                  <div class="product-title">
                    <h3 class="name">{{ item.name }}</h3>
                    <p class="description">算力值：{{ item.power }}</p>
                  </div>
                </div>

                <div class="product-content">
                  <div class="price-section">
                    <div class="price-info">
                      <span class="currency">￥</span>
                      <span class="price-value">{{ item.price }}</span>
                    </div>
                  </div>

                  <div class="features-list" v-if="item.features">
                    <div class="feature-item" v-for="feature in item.features" :key="feature">
                      <i class="iconfont icon-check"></i>
                      <span>{{ feature }}</span>
                    </div>
                  </div>
                </div>

                <div class="product-actions">
                  <div class="payment-buttons">
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
            </el-col>
          </el-row>
          <el-empty description="暂无数据" v-else :image="nodata" />
          <div class="box-card">
            <h2 class="headline">消费账单</h2>

            <div class="user-order">
              <user-order v-if="isLogin" :key="userOrderKey" />
            </div>
          </div>
        </div>
      </div>

      <password-dialog
        v-if="isLogin"
        :show="showPasswordDialog"
        @hide="showPasswordDialog = false"
      />

      <!-- 绑定手机弹窗 -->
      <el-dialog
        v-model="showBindMobileDialog"
        title="绑定手机"
        width="400px"
        :close-on-click-modal="true"
        @close="showBindMobileDialog = false"
      >
        <bind-mobile @hide="showBindMobileDialog = false" />
      </el-dialog>

      <!-- 绑定邮箱弹窗 -->
      <el-dialog
        v-model="showBindEmailDialog"
        title="绑定邮箱"
        width="400px"
        :close-on-click-modal="true"
        @close="showBindEmailDialog = false"
      >
        <bind-email @hide="showBindEmailDialog = false" />
      </el-dialog>

      <!-- 卡密兑换弹窗 -->
      <el-dialog
        v-model="showRedeemVerifyDialog"
        title="卡密兑换"
        width="450px"
        :close-on-click-modal="true"
        @close="showRedeemVerifyDialog = false"
      >
        <redeem-verify @hide="redeemCallback" />
      </el-dialog>
    </div>

    <!--支付二维码-->
    <el-dialog
      v-model="showQrCode"
      :show-close="true"
      style="width: 334px; height: 368px"
      class="pay-dialog"
    >
      <template #header>
        <div class="flex items-center justify-center text-base">
          <span style="color: var(--el-text-color-regular)">{{ title }}</span>
        </div>
      </template>
      <div class="qr-container">
        <el-image :src="qrImg" style="height: 300px; width: 300px" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import nodata from '@/assets/img/no-data.png'
import BindEmail from '@/components/BindEmail.vue'
import BindMobile from '@/components/BindMobile.vue'
import PasswordDialog from '@/components/PasswordDialog.vue'
import RedeemVerify from '@/components/RedeemVerify.vue'
import UserOrder from '@/components/UserOrder.vue'
import { checkSession, getSystemInfo } from '@/store/cache'
import { useSharedStore } from '@/store/sharedata'
import { closeLoading, showLoading } from '@/utils/dialog'
import { httpGet, httpPost } from '@/utils/http'
import { isMobile } from '@/utils/libs'
import { ElMessage } from 'element-plus'
import QRCode from 'qrcode'
import { onMounted, onUnmounted, ref } from 'vue'

const list = ref([])
const vipImg = ref('/images/menu/member.png')
const enableReward = ref(false) // 是否启用众筹功能
const rewardImg = ref('/images/reward.png')
const showPasswordDialog = ref(false)
const showBindMobileDialog = ref(false)
const showBindEmailDialog = ref(false)
const showRedeemVerifyDialog = ref(false)
const user = ref(null)
const isLogin = ref(false)
const orderTimeout = ref(1800)
const loading = ref(true)
const loadingText = ref('加载中...')

const store = useSharedStore()
const userOrderKey = ref(0)
const showQrCode = ref(false)
const qrImg = ref('')
const title = ref('')
const handler = ref(null)

onMounted(() => {
  checkSession()
    .then((_user) => {
      user.value = _user
      isLogin.value = true
    })
    .catch(() => {
      store.setShowLoginDialog(true)
    })

  httpGet('/api/product/list')
    .then((res) => {
      list.value = res.data
      loading.value = false
    })
    .catch((e) => {
      ElMessage.error('获取产品套餐失败：' + e.message)
    })

  getSystemInfo()
    .then((res) => {
      rewardImg.value = res.data['reward_img']
      enableReward.value = res.data['enabled_reward']
      if (res.data['order_pay_timeout'] > 0) {
        orderTimeout.value = res.data['order_pay_timeout']
      }
    })
    .catch((e) => {
      ElMessage.error('获取系统配置失败：' + e.message)
    })
})

const selectedPid = ref(0)
const wxPay = (product) => {
  selectedPid.value = product.id
  title.value = '请打开微信扫码支付'
  generateOrder('wxpay')
}

const alipay = (product) => {
  selectedPid.value = product.id
  title.value = '请打开支付宝扫码支付'
  generateOrder('alipay')
}

const generateOrder = (payWay) => {
  showLoading('正在生成支付订单...')
  // 生成支付订单
  httpPost('/api/payment/create', {
    pid: selectedPid.value,
    pay_way: payWay,
    domain: `${window.location.protocol}//${window.location.host}`,
    device: isMobile() ? 'mobile' : 'pc',
  })
    .then((res) => {
      closeLoading()

      if (isMobile()) {
        window.location.href = res.data.pay_url
      } else {
        QRCode.toDataURL(res.data.pay_url, { width: 300, height: 300, margin: 2 }, (error, url) => {
          if (error) {
            console.error(error)
          } else {
            qrImg.value = url
          }
        })
        // 查询订单状态
        if (handler.value) {
          clearTimeout(handler.value)
        }
        handler.value = setTimeout(() => queryOrder(res.data.order_no), 3000)
        showQrCode.value = true
      }
    })
    .catch((e) => {
      closeLoading()
      ElMessage.error('生成支付订单失败：' + e.message)
    })
}

const queryOrder = async (orderNo) => {
  const res = await httpGet('/api/order/query?order_no=' + orderNo)
  if (res?.data.status === 2) {
    // 订单支付成功
    clearTimeout(handler.value)
    ElMessage.success('支付成功')
    showQrCode.value = false
    // 更新用户积分
    user.value.scores += res.data.credit
  } else {
    handler.value = setTimeout(() => queryOrder(orderNo), 3000)
  }
}

const redeemCallback = (success) => {
  showRedeemVerifyDialog.value = false

  if (success) {
    userOrderKey.value += 1
  }
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
@use '@/assets/css/custom-scroll.scss' as *;
@use '@/assets/css/member.scss' as *;

// 支付弹窗样式优化
.pay-dialog {
  .qr-container {
    text-align: center;
    position: relative;

    .qr-overlay {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(255, 255, 255, 0.9);
      display: flex;
      align-items: center;
      justify-content: center;

      .success-text {
        background: #67c23a;
        color: white;
        border-radius: 4px;
        font-size: 14px;
      }
    }
  }
}

// 支付按钮样式
.pay-way {
  .row {
    margin: 0;

    .col {
      padding: 0 5px;

      button {
        border: none;
        cursor: pointer;
        transition: all 0.3s ease;
        font-size: 14px;
        font-weight: 500;

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
        }
      }
    }
  }
}
</style>
