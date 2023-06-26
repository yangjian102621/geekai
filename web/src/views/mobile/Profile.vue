<template>
  <div class="mobile-user-profile container">
    <van-nav-bar :title="title"/>

    <div class="content">
      <van-form @submit="onSubmit">
        <van-cell-group inset>
          <van-field
              v-model="username"
              name="用户名"
              label="用户名"
              placeholder="用户名"
              :rules="[{ required: true, message: '请填写用户名' }]"
          />
          <van-field
              v-model="password"
              type="password"
              name="密码"
              label="密码"
              placeholder="密码"
              :rules="[{ required: true, message: '请填写密码' }]"
          />
        </van-cell-group>
        <div style="margin: 16px;">
          <van-button round block type="primary" native-type="submit">
            提交
          </van-button>
        </div>
      </van-form>
    </div>
  </div>
</template>

<script setup>
import {ref} from "vue";

const title = ref('用户设置')
</script>

<style scoped>

</style>