<template>
  <div class="theme-box" @click="toggleTheme">
    <span class="iconfont icon-yueliang">{{
      themePage === "light" ? "&#xe679;" : "&#xe60b;"
    }}</span>
  </div>
</template>

<script setup>
import { onMounted, ref } from "vue";

// 定义主题状态，初始值从 localStorage 获取
const themePage = ref(localStorage.getItem("theme") || "light");

// 切换主题函数
const toggleTheme = () => {
  themePage.value = themePage.value === "light" ? "dark" : "light";
  document.documentElement.setAttribute("data-theme", themePage.value); // 设置 HTML 的 data-theme 属性
  localStorage.setItem("theme", themePage.value); // 保存主题到 localStorage
};

onMounted(() => {
  document.documentElement.setAttribute("data-theme", themePage.value);
});
</script>

<style lang="stylus" scoped>
@import '@/assets/iconfont/iconfont.css'
.theme-box{
  position: fixed;
  right: 40px;
  bottom: 262px;
  cursor: pointer;
  background-color: #fff;
  border: 1px solid #ccc;
  border-radius: 50%;
  width 35px;
  height: 35px;
  line-height: 35px;
  text-align: center;
  background-color: rgb(146, 147, 148);

  transition: all 0.3s ease;
  &:hover{
    transform: scale(1.1);
  }
  &:active{
    transform: scale(0.9);
  }
  .iconfont{
    font-size: 20px;
    color: yellow;
    transition: transform 0.3s ease;
  }
}
</style>
