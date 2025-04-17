<template>
  <div>
    <ThemeChange />
    <div
      @click="goBack"
      class="flex back animate__animated animate__pulse animate__infinite"
    >
      <el-icon><ArrowLeftBold /></el-icon
      >{{ title === "注册" ? "首页" : "返回" }}
    </div>
    <div class="title">{{ title }}</div>
    <div class="smTitle" v-if="title !== '重置密码'">
      {{ title === "登录" ? "没有账号？" : "已有账号？"
      }}<span @click="goPageFun" class="text-color-primary sign"
        >赶紧{{ title === "登录" ? "注册" : "登录" }}</span
      >
    </div>
    <slot></slot>
    <div class="flex orline" v-if="title !== '重置密码'">
      <div class="lineor"></div>
      <span>或</span>
      <div class="lineor"></div>
    </div>
  </div>
</template>

<script setup>
import { ArrowLeftBold } from "@element-plus/icons-vue";
import ThemeChange from "@/components/ThemeChange.vue";
import { defineProps } from "vue";
import { useRouter } from "vue-router";

const props = defineProps({
  title: {
    type: String,
    default: "登录"
  },
  smTitle: { type: String, default: "没有账号？" },
  goPage: {
    type: String,
    default: "/register"
  }
});

const router = useRouter();
const goBack = () => {
  if (props.title === "注册") {
    router.push("/");
  } else {
    router.go(-1);
  }
};
const goPageFun = () => {
  if (props.title === "登录") {
    router.push("/register");
  } else {
    router.push("/login");
  }
};
</script>

<style lang="stylus" scoped>
.back{
   color:var(--sm-txt)
   font-size: 14px;
   margin-bottom: 140px
   margin-top: 18px
   cursor: pointer
   .el-icon{
     margin-right: 6px
   }
 }
 .title{
   font-size: 36px
   margin-bottom: 16px
   color: var(--text-color)

 }
 .smTitle{
    color: var(--text-color)
   font-size: 14px;
   margin-bottom: 36px
 }
 .sign{
    text-decoration: underline;
    cursor :pointer
  }
  .orline{
    color:var(--text-secondary)
    span{
      font-size: 14px;
      margin: 0 10px

    }
    .lineor{

      width: 182px; height: 1px;
      background: var(--text-secondary)
    }
  }
</style>
