<template>
  <div class="member custom-scroll">
    <div class="title">
      会员充值中心
    </div>
    <div class="inner" :style="{height: listBoxHeight + 'px'}">

      <div class="user-profile">
        <user-profile/>
      </div>

      <div class="product-box">
        <div class="info">
          <el-alert type="info" show-icon :closable="false" effect="dark">
            <strong>说明:</strong> 成为本站会员后每月有500次对话额度，50次 AI 绘画额度，限制下月1号解除，若在期间超过次数后可单独购买点卡。
            当月充值的点卡有效期可以延期到下个月底。
          </el-alert>
        </div>

        <ItemList :items="list" v-if="list.length > 0" :gap="30" :width="200">
          <template #default="scope">
            <div class="product-item" :style="{width: scope.width+'px'}" @click="orderPay(scope.item)">
              <div class="image-container">
                <el-image :src="vipImg" fit="cover"/>
              </div>
              <div class="product-title">
                <span class="name">{{ scope.item.name }}</span>
              </div>
              <div class="product-info">
                <div class="info-line">
                  <span class="label">商品原价：</span>
                  <span class="price">￥{{ scope.item.price }}</span>
                </div>
                <div class="info-line">
                  <span class="label">促销立减：</span>
                  <span class="price">￥{{ scope.item.discount }}</span>
                </div>
                <div class="info-line">
                  <span class="label">有效期：</span>
                  <span class="expire" v-if="scope.item.days > 0">{{ scope.item.days }}天</span>
                  <span class="expire" v-else>当月有效</span>
                </div>
              </div>
            </div>
          </template>
        </ItemList>
      </div>
    </div>

    <login-dialog :show="showLoginDialog" @hide="showLoginDialog = false"/>

    <el-dialog
        v-model="showPayDialog"
        :close-on-click-modal="false"
        :show-close="true"
        :width="400"
        title="用户登录">
      <div class="pay-container">
        <div class="pay-qrcode">
          <el-image :src="qrcode"/>
        </div>

        <div class="tip success" v-if="text !== ''">
          <el-icon>
            <SuccessFilled/>
          </el-icon>
          <span class="text">{{ text }}</span>
        </div>
        <div class="tip" v-else>
          <el-icon>
            <InfoFilled/>
          </el-icon>
          <span class="text">请打开手机支付宝扫码支付</span>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import {nextTick, onMounted, ref} from "vue"
import {ElMessage} from "element-plus";
import {httpGet, httpPost} from "@/utils/http";
import ItemList from "@/components/ItemList.vue";
import {Delete, InfoFilled, Plus, SuccessFilled} from "@element-plus/icons-vue";
import LoginDialog from "@/components/LoginDialog.vue";
import {checkSession} from "@/action/session";
import {arrayContains, removeArrayItem, substr} from "@/utils/libs";
import router from "@/router";
import UserProfile from "@/components/UserProfile.vue";

const listBoxHeight = window.innerHeight - 97
const list = ref([])
const showLoginDialog = ref(false)
const showPayDialog = ref(false)
const elements = ref(null)
const vipImg = ref("/images/vip.png")
const qrcode = ref("")
const amount = ref(0)
const discount = ref(0)
const text = ref("")
onMounted(() => {
  httpGet("/api/product/list").then((res) => {
    list.value = res.data
  }).catch(e => {
    ElMessage.error("获取产品套餐失败：" + e.message)
  })
})

const orderPay = (row) => {
  checkSession().then(user => {
    console.log(row)
    httpPost("/api/payment/alipay/qrcode", {product_id: row.id, user_id: user.id}).then(res => {
      console.log(res)
      showPayDialog.value = true
      qrcode.value = res.data['image']
      queryOrder(res.data['order_no'])
    }).catch(e => {
      ElMessage.error("生成支付订单失败：" + e.message)
    })
  }).catch(e => {
    console.log(e)
    showLoginDialog.value = true
  })
}

const queryOrder = (orderNo) => {
  httpPost("/api/payment/query", {order_no: orderNo}).then(res => {
    if (res.data.status === 1) {
      text.value = "扫码成功，请在手机上进行支付！"
      queryOrder(orderNo)
    } else if (res.data.status === 2) {
      text.value = "支付成功，正在刷新页面"
      setTimeout(() => location.reload(), 500)
    } else {
      queryOrder(orderNo)
    }
  }).catch(e => {
    ElMessage.error("查询支付状态失败：" + e.message)
  })
}

</script>

<style lang="stylus">
@import "@/assets/css/custom-scroll.styl"
.member {
  background-color: #282c34;
  height 100vh

  .el-dialog {
    .el-dialog__body {
      padding-top 0

      .pay-container {
        .pay-qrcode {
          display flex
          justify-content center

          .el-image {
            width 360px;
            height 360px;
          }
        }

        .tip {
          display flex
          justify-content center

          .el-icon {
            font-size 24px
          }

          .text {
            font-size: 16px
            margin-left 10px
          }
        }

        .tip.success {
          color #07c160
        }
      }
    }
  }

  .title {
    text-align center
    background-color #25272d
    font-size 24px
    color #ffffff
    padding 10px
    border-bottom 1px solid #3c3c3c
  }

  .inner {
    display flex
    color #ffffff
    padding 15px;
    overflow-y visible
    overflow-x hidden

    .user-profile {
      padding 10px 20px
      background-color #393F4A
      color #ffffff
      border-radius 10px

      .el-form-item__label {
        color #ffffff
        justify-content start
      }
    }

    .product-box {
      padding 0 10px

      .info {
        .el-alert__description {
          font-size 14px !important
          margin 0
        }
        padding 10px 20px
      }

      .list-box {
        .product-item {
          border 1px solid #666666
          border-radius 6px
          overflow hidden
          cursor pointer
          transition: all 0.3s ease; /* 添加过渡效果 */

          .image-container {
            display flex
            justify-content center

            .el-image {
              padding 6px

              .el-image__inner {
                border-radius 10px
              }
            }
          }

          .product-title {
            display flex
            padding 10px

            .name {
              width 100%
              text-align center
              font-size 16px
              font-weight bold
              color #47fff1
            }
          }

          .product-info {
            padding 10px 20px
            font-size 14px
            color #999999

            .info-line {
              display flex
              width 100%
              padding 5px 0

              .label {
                display flex
                width 100%
              }

              .price, .expire {
                display flex
                width 90px
                justify-content right
              }

              .price {
                color #f56c6c
              }

              .expire {
                color #409eff
              }
            }

          }

          &:hover {
            box-shadow: 0 0 10px rgba(71, 255, 241, 0.6); /* 添加阴影效果 */
            transform: translateY(-10px); /* 向上移动10像素 */
          }
        }
      }
    }
  }

}
</style>
