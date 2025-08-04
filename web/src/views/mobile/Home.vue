<template>
  <van-config-provider :theme="theme">
    <div class="mobile-home">
      <div class="page-content">
        <router-view :key="routerViewKey" v-slot="{ Component }">
          <transition name="move" mode="out-in">
            <component :is="Component"></component>
          </transition>
        </router-view>
      </div>

      <van-tabbar route v-model="active" :safe-area-inset-bottom="true">
        <van-tabbar-item to="/mobile/index" name="home" icon="home-o">
          <span>首页</span>
          <template #icon="props">
            <i class="iconfont icon-house" :class="{ 'active-icon': props.active }"></i>
          </template>
        </van-tabbar-item>
        <van-tabbar-item to="/mobile/chat" name="chat" icon="chat-o">
          <span>对话</span>
          <template #icon="props">
            <i class="iconfont icon-chat" :class="{ 'active-icon': props.active }"></i>
          </template>
        </van-tabbar-item>
        <van-tabbar-item to="/mobile/create" name="create" icon="plus">
          <span>创作</span>
          <template #icon="props">
            <i class="iconfont icon-mj" :class="{ 'active-icon': props.active }"></i>
          </template>
        </van-tabbar-item>
        <van-tabbar-item to="/mobile/discover" name="discover" icon="apps-o">
          <span>发现</span>
          <template #icon="props">
            <i class="iconfont icon-more" :class="{ 'active-icon': props.active }"></i>
          </template>
        </van-tabbar-item>
        <van-tabbar-item to="/mobile/profile" name="profile" icon="user-o">
          <span>我的</span>
          <template #icon="props">
            <i class="iconfont icon-user-circle" :class="{ 'active-icon': props.active }"></i>
          </template>
        </van-tabbar-item>
      </van-tabbar>
    </div>
  </van-config-provider>
</template>

<script setup>
import { useSharedStore } from '@/store/sharedata'
import { onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const active = ref('home')
const store = useSharedStore()
const theme = ref(store.theme)
const route = useRoute()
const router = useRouter()
const routerViewKey = ref(0)

// 监听路由变化，强制刷新组件
watch(
  () => route.path,
  () => {
    routerViewKey.value += 1
  }
)

watch(
  () => store.theme,
  (val) => {
    theme.value = val
  }
)

// 路由守卫
router.beforeEach((to, from, next) => {
  // 可以在这里添加路由权限检查等逻辑
  next()
})

onMounted(() => {
  // 组件挂载时的初始化逻辑
})
</script>

<style lang="scss">
@use '../../assets/iconfont/iconfont.css' as *;

.mobile-home {
  .page-content {
    padding-bottom: 60px;
  }

  .van-tabbar {
    box-shadow: 0 -2px 20px rgba(0, 0, 0, 0.1);

    .van-tabbar-item {
      .active-icon {
        color: var(--van-primary-color) !important;
        transform: scale(1.1);
        transition: all 0.3s ease;
      }

      &--active {
        .van-tabbar-item__text {
          color: var(--van-primary-color);
          font-weight: 600;
        }
      }
    }

    .iconfont {
      font-size: 20px;
      transition: all 0.3s ease;
    }
  }
}

// 黑色主题
.van-theme-dark body {
  background: #1c1c1e;
}

// 路由切换动画
.move-enter-active,
.move-leave-active {
  transition: all 0.3s ease;
}

.move-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.move-leave-to {
  opacity: 0;
  transform: translateX(-100%);
}

.move-enter-to,
.move-leave-from {
  opacity: 1;
  transform: translateX(0);
}
</style>
