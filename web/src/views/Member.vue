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
          <!-- <div class="info" v-if="orderPayInfoText !== ''">
            <el-alert type="success" show-icon :closable="false" effect="dark">
              <strong>说明:</strong> {{ vipInfoText }}
            </el-alert>
          </div> -->

          <el-row v-if="list.length > 0" :gutter="20" class="list-box">
            <el-col v-for="item in list" :key="item" :span="6">
              <div class="product-item">
                <div class="image-container">
                  <el-image :src="vipImg" fit="cover" />
                </div>
                <div class="product-title">
                  <span class="name">{{ item.name }}</span>
                </div>
                <div class="product-info">
                  <div class="info-line">
                    <span class="label">商品原价：</span>
                    <span class="price"
                      ><del>￥{{ item.price }}</del></span
                    >
                  </div>
                  <div class="info-line">
                    <span class="label">优惠价：</span>
                    <span class="discount">￥{{ item.discount }}</span>
                  </div>
                  <div class="info-line">
                    <span class="label">有效期：</span>
                    <span class="expire" v-if="item.days > 0">{{ item.days }}天</span>
                    <span class="expire" v-else>长期有效</span>
                  </div>

                  <div class="info-line">
                    <span class="label">算力值：</span>
                    <span class="power">{{ item.power }}</span>
                  </div>

                  <div class="pay-way">
                    <span
                      type="primary"
                      v-for="payWay in payWays"
                      @click="pay(item, payWay)"
                      :key="payWay"
                    >
                      <el-button v-if="payWay.pay_type === 'alipay'" color="#15A6E8" circle>
                        <i class="iconfont icon-alipay"></i>
                      </el-button>
                      <el-button v-else-if="payWay.pay_type === 'qqpay'" circle>
                        <i class="iconfont icon-qq"></i>
                      </el-button>
                      <el-button v-else-if="payWay.pay_type === 'paypal'" class="paypal" round>
                        <i class="iconfont icon-paypal"></i>
                      </el-button>
                      <el-button v-else-if="payWay.pay_type === 'jdpay'" color="#E1251B" circle>
                        <i class="iconfont icon-jd-pay"></i>
                      </el-button>
                      <el-button v-else-if="payWay.pay_type === 'douyin'" class="douyin" circle>
                        <i class="iconfont icon-douyin"></i>
                      </el-button>
                      <el-button v-else circle class="wechat" color="#67C23A">
                        <i class="iconfont icon-wechat-pay"></i>
                      </el-button>
                    </span>
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

    <el-dialog
      v-model="showDialog"
      :show-close="false"
      :close-on-click-modal="false"
      hide-footer
      width="auto"
      class="pay-dialog"
    >
      <div v-if="qrImg !== ''">
        <div class="product-info">
          请使用微信扫码支付：<span class="price">￥{{ price }}</span>
        </div>
        <el-image :src="qrImg" fit="cover" />
      </div>
      <div style="padding-bottom: 10px; text-align: center">
        <el-button type="success" @click="payCallback(true)">支付成功</el-button>
        <el-button type="danger" @click="payCallback(false)">支付失败</el-button>
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
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import QRCode from 'qrcode'
import { onMounted, ref } from 'vue'

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
const orderPayInfoText = ref('')

const payWays = ref([])
const vipInfoText = ref('')
const store = useSharedStore()
const userOrderKey = ref(0)
const showDialog = ref(false)
const qrImg = ref('')
const price = ref(0)

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
      orderPayInfoText.value = res.data['order_pay_info_text']
      if (res.data['order_pay_timeout'] > 0) {
        orderTimeout.value = res.data['order_pay_timeout']
      }
      vipInfoText.value = res.data['vip_info_text']
    })
    .catch((e) => {
      ElMessage.error('获取系统配置失败：' + e.message)
    })

  httpGet('/api/payment/payWays')
    .then((res) => {
      payWays.value = res.data
    })
    .catch((e) => {
      ElMessage.error('获取支付方式失败：' + e.message)
    })
})

const pay = (product, payWay) => {
  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return
  }
  loading.value = true
  loadingText.value = '正在生成支付订单...'
  let host = import.meta.env.VITE_API_HOST
  if (host === '') {
    host = `${location.protocol}//${location.host}`
  }
  httpPost(`${import.meta.env.VITE_API_HOST}/api/payment/doPay`, {
    product_id: product.id,
    pay_way: payWay.pay_way,
    pay_type: payWay.pay_type,
    user_id: user.value.id,
    host: host,
    device: 'jump',
  })
    .then((res) => {
      showDialog.value = true
      loading.value = false
      if (payWay.pay_way === 'wechat') {
        price.value = Number(product.discount)
        QRCode.toDataURL(res.data, { width: 300, height: 300, margin: 2 }, (error, url) => {
          if (error) {
            console.error(error)
          } else {
            qrImg.value = url
          }
        })
      } else {
        window.open(res.data, '_blank')
      }
    })
    .catch((e) => {
      setTimeout(() => {
        ElMessage.error('生成支付订单失败：' + e.message)
        loading.value = false
      }, 500)
    })
}

const redeemCallback = (success) => {
  showRedeemVerifyDialog.value = false

  if (success) {
    userOrderKey.value += 1
  }
}

const payCallback = (success) => {
  showDialog.value = false
  if (success) {
    userOrderKey.value += 1
  }
}
</script>

<style lang="scss" scoped>
@use '../assets/css/custom-scroll.scss' as *;
@use '../assets/css/member.scss' as *;
</style>
