<template>
  <div>
    <div class="member custom-scroll" v-loading="loading" element-loading-background="rgba(255,255,255,.3)" :element-loading-text="loadingText">
      <div class="inner">
        <div class="user-profile">
          <user-profile :key="profileKey"/>

          <el-row class="user-opt" :gutter="20">
            <el-col :span="12">
              <el-button type="primary" @click="showBindEmailDialog = true">绑定邮箱</el-button>
            </el-col>
            <el-col :span="12">
              <el-button type="primary" @click="showBindMobileDialog = true">绑定手机</el-button>
            </el-col>
            <el-col :span="12">
              <el-button type="primary" @click="showThirdLoginDialog = true">第三方登录</el-button>
            </el-col>
            <el-col :span="12">
              <el-button type="primary" @click="showPasswordDialog = true">修改密码</el-button>
            </el-col>
            <el-col :span="24">
              <el-button type="primary" @click="showRedeemVerifyDialog = true">卡密兑换
              </el-button>
            </el-col>
          </el-row>
        </div>

        <div class="product-box">
          <div class="info" v-if="orderPayInfoText !== ''">
            <el-alert type="success" show-icon :closable="false" effect="dark">
              <strong>说明:</strong> {{ vipInfoText }}
            </el-alert>
          </div>

          <el-row v-if="list.length > 0" :gutter="20" class="list-box">
            <el-col v-for="item in list" :key="item" :span="6">
              <div class="product-item">
                <div class="image-container">
                  <el-image :src="vipImg" fit="cover"/>
                </div>
                <div class="product-title">
                  <span class="name">{{ item.name }}</span>
                </div>
                <div class="product-info">
                  <div class="info-line">
                    <span class="label">商品原价：</span>
                    <span class="price"><del>￥{{ item.price }}</del></span>
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

                    <span type="primary" v-for="payWay in payWays" @click="pay(item,payWay)" :key="payWay">
                      <el-button v-if="payWay.pay_type==='alipay'" color="#15A6E8" circle>
                        <i class="iconfont icon-alipay" ></i>
                      </el-button>
                      <el-button v-else-if="payWay.pay_type==='qqpay'" circle>
                        <i class="iconfont icon-qq"></i>
                      </el-button>
                      <el-button v-else-if="payWay.pay_type==='paypal'" class="paypal" round>
                        <i class="iconfont icon-paypal"></i>
                      </el-button>
                      <el-button v-else-if="payWay.pay_type==='jdpay'" color="#E1251B" circle>
                        <i class="iconfont icon-jd-pay"></i>
                      </el-button>
                      <el-button v-else-if="payWay.pay_type==='douyin'" class="douyin" circle>
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
          <el-empty description="暂无数据" v-else />

          <h2 class="headline">消费账单</h2>

          <div class="user-order">
            <user-order v-if="isLogin" :key="userOrderKey"/>
          </div>
        </div>
      </div>

      <password-dialog v-if="isLogin" :show="showPasswordDialog" @hide="showPasswordDialog = false"/>
      <bind-mobile v-if="isLogin" :show="showBindMobileDialog" @hide="showBindMobileDialog = false"/>
      <bind-email v-if="isLogin" :show="showBindEmailDialog" @hide="showBindEmailDialog = false"/>
      <third-login v-if="isLogin" :show="showThirdLoginDialog" @hide="showThirdLoginDialog = false"/>
      <redeem-verify v-if="isLogin" :show="showRedeemVerifyDialog" @hide="redeemCallback"/>

    </div>

    <el-dialog v-model="showDialog" :show-close=false :close-on-click-modal="false" hide-footer width="auto" class="pay-dialog">
      <div v-if="qrImg !== ''">
        <div class="product-info">请使用微信扫码支付：<span class="price">￥{{price}}</span></div>
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
import {onMounted, ref} from "vue"
import {ElMessage} from "element-plus";
import {httpGet, httpPost} from "@/utils/http";
import {checkSession, getSystemInfo} from "@/store/cache";
import UserProfile from "@/components/UserProfile.vue";
import PasswordDialog from "@/components/PasswordDialog.vue";
import BindMobile from "@/components/BindMobile.vue";
import RedeemVerify from "@/components/RedeemVerify.vue";
import UserOrder from "@/components/UserOrder.vue";
import {useSharedStore} from "@/store/sharedata";
import BindEmail from "@/components/BindEmail.vue";
import ThirdLogin from "@/components/ThirdLogin.vue";
import QRCode from "qrcode";

const list = ref([])
const vipImg = ref("/images/vip.png")
const enableReward = ref(false) // 是否启用众筹功能
const rewardImg = ref('/images/reward.png')
const showPasswordDialog = ref(false)
const showBindMobileDialog = ref(false)
const showBindEmailDialog = ref(false)
const showRedeemVerifyDialog = ref(false)
const showThirdLoginDialog = ref(false)
const user = ref(null)
const isLogin = ref(false)
const orderTimeout = ref(1800)
const loading = ref(true)
const loadingText = ref("加载中...")
const orderPayInfoText = ref("")

const payWays = ref([])
const vipInfoText = ref("")
const store = useSharedStore()
const profileKey = ref(0)
const userOrderKey = ref(0)
const showDialog = ref(false)
const qrImg = ref("")
const price = ref(0)


onMounted(() => {
  checkSession().then(_user => {
    user.value = _user
    isLogin.value = true
  }).catch(() => {
    store.setShowLoginDialog(true)
  })

  httpGet("/api/product/list").then((res) => {
    list.value = res.data
    loading.value = false
  }).catch(e => {
    ElMessage.error("获取产品套餐失败：" + e.message)
  })

  getSystemInfo().then(res => {
    rewardImg.value = res.data['reward_img']
    enableReward.value = res.data['enabled_reward']
    orderPayInfoText.value = res.data['order_pay_info_text']
    if (res.data['order_pay_timeout'] > 0) {
      orderTimeout.value = res.data['order_pay_timeout']
    }
    vipInfoText.value = res.data['vip_info_text']
  }).catch(e => {
    ElMessage.error("获取系统配置失败：" + e.message)
  })

  httpGet("/api/payment/payWays").then(res => {
    payWays.value = res.data
  }).catch(e => {
    ElMessage.error("获取支付方式失败：" + e.message)
  })
})

const pay = (product, payWay) => {
  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return
  }
  loading.value = true
  loadingText.value = "正在生成支付订单..."
  let host = process.env.VUE_APP_API_HOST
  if (host === '') {
    host = `${location.protocol}//${location.host}`;
  }
  httpPost(`${process.env.VUE_APP_API_HOST}/api/payment/doPay`, {
    product_id: product.id,
    pay_way: payWay.pay_way,
    pay_type: payWay.pay_type,
    user_id: user.value.id,
    host: host,
    device: "jump"
  }).then(res => {
    showDialog.value = true
    loading.value = false
    if (payWay.pay_way === 'wechat') {
      price.value = Number(product.discount)
      QRCode.toDataURL(res.data, {width: 300, height: 300, margin: 2}, (error, url) => {
        if (error) {
          console.error(error)
        } else {
          qrImg.value = url;
        }
      })
    } else {
      window.open(res.data, '_blank');
    }
  }).catch(e => {
    setTimeout(() => {
      ElMessage.error("生成支付订单失败：" + e.message)
      loading.value = false
    }, 500)
  })
}

const redeemCallback = (success) => {
  showRedeemVerifyDialog.value = false
  if (success) {
    profileKey.value += 1
  }
}

const payCallback = (success) => {
  showDialog.value = false
  if (success) {
    profileKey.value += 1
    userOrderKey.value += 1
  }
}


</script>

<style lang="stylus">
@import "@/assets/css/custom-scroll.styl"
@import "@/assets/css/member.styl"
</style>
