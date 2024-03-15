<template>
  <div class="mobile-user-profile container">
    <van-nav-bar :title="title"/>

    <div class="content">
      <van-form>
        <van-cell-group inset v-model="form">
          <van-field
              v-model="form.username"
              name="账号"
              label="账号"
              readonly
              disabled
          />
          <van-field label="头像">
            <template #input>
              <van-uploader v-model="fileList"
                            reupload max-count="1"
                            :deletable="false"
                            :after-read="afterRead"/>
            </template>
          </van-field>

          <van-field label="剩余对话次数">
            <template #input>
              <van-tag type="primary">{{ form.calls }}</van-tag>
            </template>
          </van-field>

          <van-field label="剩余绘图次数">
            <template #input>
              <van-tag type="primary">{{ form.img_calls }}</van-tag>
            </template>
          </van-field>

          <van-field label="累计算力消耗">
            <template #input>
              <van-tag type="primary">{{ form.total_tokens }}</van-tag>
            </template>
          </van-field>
        </van-cell-group>
      </van-form>

      <div class="modify-pass">
        <van-button round block type="primary" @click="showPasswordDialog = true">修改密码</van-button>
      </div>

      <div class="product-list">
        <h3>充值套餐</h3>
        <div class="item" v-for="item in products" :key="item.id">
          <h4 class="title">
            <span>{{ item.name }}</span>
            <van-button type="primary" size="small" icon="cart">立即购买</van-button>
          </h4>

          <van-cell-group>
            <van-cell title="商品价格">
              <span class="price">
                ￥{{ (item.price - item.discount).toFixed(2) }}
              </span>
              （
              <del>￥{{ item.price }}</del>
              ）
            </van-cell>
            <van-cell title="有效期">
              <span v-if="item.days > 0">{{ item.days }}天</span>
              <van-tag type="primary" v-else>长期有效</van-tag>
            </van-cell>
            <van-cell title="对话次数">
              <span v-if="item.calls > 0">{{ item.calls }}</span>
              <span v-else>{{ vipMonthCalls }}</span>
            </van-cell>
            <van-cell title="绘图次数">
              <span v-if="item.img_calls > 0">{{ item.img_calls }}</span>
              <span v-else>{{ vipMonthImgCalls }}</span>
            </van-cell>
          </van-cell-group>
        </div>
      </div>

    </div>

    <van-dialog v-model:show="showPasswordDialog" title="修改密码" show-cancel-button
                @confirm="updatePass"
                @cancel="showPasswordDialog = false"
                :before-close="beforeClose">
      <van-form>
        <van-cell-group inset>
          <van-field
              v-model="pass.old"
              placeholder="旧密码"
          />
          <van-field
              v-model="pass.new"
              type="password"
              placeholder="新密码"
          />
          <van-field
              v-model="pass.renew"
              type="password"
              placeholder="确认密码"
          />
        </van-cell-group>
      </van-form>
    </van-dialog>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue";
import {showFailToast, showNotify, showSuccessToast} from "vant";
import {httpGet, httpPost} from "@/utils/http";
import Compressor from 'compressorjs';
import {ElMessage} from "element-plus";

const title = ref('用户设置')
const form = ref({
  username: '',
  nickname: '',
  mobile: '',
  avatar: '',
  calls: 0,
  tokens: 0
})
const fileList = ref([
  {
    url: '',
    message: '上传中...',
  }
]);

const products = ref([])
const vipMonthCalls = ref(0)
const vipMonthImgCalls = ref(0)

onMounted(() => {
  httpGet('/api/user/profile').then(res => {
    form.value = res.data
    fileList.value[0].url = form.value.avatar
  }).catch((e) => {
    console.log(e.message)
    showFailToast('获取用户信息失败')
  });

  // 获取产品列表
  httpGet("/api/product/list").then((res) => {
    products.value = res.data
  }).catch(e => {
    showFailToast("获取产品套餐失败：" + e.message)
  })

  httpGet("/api/admin/config/get?key=system").then(res => {
    vipMonthCalls.value = res.data['vip_month_calls']
    vipMonthImgCalls.value = res.data['vip_month_img_calls']
  }).catch(e => {
    showFailToast("获取系统配置失败：" + e.message)
  })
})

const afterRead = (file) => {
  file.status = 'uploading';
  file.message = '上传中...';
  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData();
      formData.append('file', result, result.name);
      // 执行上传操作
      httpPost('/api/upload', formData).then((res) => {
        form.value.avatar = res.data.url
        file.status = 'success'
        httpPost('/api/user/profile/update', form.value).then(() => {
          showSuccessToast('上传成功')
        }).catch(() => {
          showFailToast('上传失败')
        })
      }).catch((e) => {
        showNotify({type: 'danger', message: '上传失败：' + e.message})
      })
    },
    error(err) {
      console.log(err.message);
    },
  });
}

const showPasswordDialog = ref(false)
const pass = ref({
  old: "",
  new: "",
  renew: ""
})

const beforeClose = (action) => {
  new Promise((resolve) => {
    resolve(action === 'confirm');
  });
}


// 提交修改密码
const updatePass = () => {
  if (pass.value.old === '') {
    return showNotify({type: "danger", message: "请输入旧密码"})
  }
  if (!pass.value.new || pass.value.new.length < 8) {
    return showNotify({type: "danger", message: "密码的长度为8-16个字符"})
  }
  if (pass.value.renew !== pass.value.new) {
    return showNotify({type: "danger", message: "两次输入密码不一致"})
  }
  httpPost('/api/user/password', {
    old_pass: pass.value.old,
    password: pass.value.new,
    repass: pass.value.renew
  }).then(() => {
    showSuccessToast("更新成功！")
    showPasswordDialog.value = false
  }).catch((e) => {
    showFailToast('更新失败，' + e.message)
    showPasswordDialog.value = false
  })
}
</script>

<style lang="stylus">
.mobile-user-profile {
  .content {
    .van-field__label {
      width 100px
      text-align right
    }

    .modify-pass {
      padding 10px 15px
    }

    .product-list {
      padding 0 15px

      .item {
        border 1px solid #e5e5e5
        border-radius 10px
        margin-bottom 15px
        overflow hidden

        .title {
          padding 0 12px
          position relative

          .van-button {
            position absolute
            top -5px
            right 10px
          }
        }

        .price {
          font-size 18px
          color #f56c6c
        }
      }
    }
  }
}
</style>