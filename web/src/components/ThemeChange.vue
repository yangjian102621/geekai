<template>
  <div class="theme-box" @click="toggleTheme" :class="size">
    <i class="iconfont" :class="themePage === 'light' ? 'icon-yueliang' : 'icon-taiyang'"></i>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useSharedStore } from "@/store/sharedata";

const props = defineProps({
  size: {
    type: String,
    default: "",
  },
});

// 定义主题状态，初始值从 localStorage 获取
const store = useSharedStore();
const themePage = ref(store.theme || "light");

// 切换主题函数
const toggleTheme = () => {
  themePage.value = themePage.value === "light" ? "dark" : "light";
  store.setTheme(themePage.value); // 保存主题
};
</script>

<style lang="stylus" scoped>
.theme-box{
  z-index :111
  position: fixed;
  right: 40px;
  bottom: 150px;
  cursor: pointer;
  border: 1px solid #ccc;
  border-radius: 50%;
  width 35px;
  height: 35px;
  line-height: 35px;
  text-align: center;
  // background-color: rgb(146, 147, 148);
  background: linear-gradient(135deg, rgba(134, 140, 255, 1) 0%, rgba(67, 24, 255, 1) 100%);

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

.theme-box.small {
  position: relative !important;
  right: initial;
  bottom: initial;
  width: 20px;
  height: 20px;
  line-height: 18px;

  .iconfont {
    font-size: 15px !important;
  }
}
</style>
