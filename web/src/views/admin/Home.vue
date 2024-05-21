<template>
  <div class="admin-home" v-if="isLogin">
    <admin-sidebar v-model:theme="theme"/>
    <div class="content-box" :class="{ 'content-collapse': sidebar.collapse }">
      <admin-header v-model:theme="theme" @changeTheme="changeTheme"/>
      <admin-tags v-model:theme="theme"/>
      <div :class="'content '+theme" :style="{height:contentHeight+'px'}">
        <router-view v-slot="{ Component }">
          <transition name="move" mode="out-in">
            <keep-alive :include="tags.nameList">
              <component :is="Component"></component>
            </keep-alive>
          </transition>
        </router-view>
      </div>
    </div>
  </div>
</template>
<script setup>
import {useSidebarStore} from '@/store/sidebar';
import {useTagsStore} from '@/store/tags';
import AdminHeader from "@/components/admin/AdminHeader.vue";
import AdminSidebar from "@/components/admin/AdminSidebar.vue";
import AdminTags from "@/components/admin/AdminTags.vue";
import {useRouter} from "vue-router";
import {checkAdminSession} from "@/action/session";
import {ref} from "vue";
import {getAdminTheme, setAdminTheme} from "@/store/system";

const sidebar = useSidebarStore();
const tags = useTagsStore();
const isLogin = ref(false)
const contentHeight = window.innerHeight - 80
const theme = ref(getAdminTheme())

// 获取会话信息
const router = useRouter();
checkAdminSession().then(() => {
  isLogin.value = true
}).catch(() => {
  router.replace('/admin/login')
})

const changeTheme = (value) => {
  if (value) {
    theme.value = 'dark'
  } else {
    theme.value = 'light'
  }
  setAdminTheme(theme.value)
}

</script>

<style scoped lang="stylus">
@import '@/assets/css/color-dark.styl';
@import '@/assets/css/main.styl';
@import '@/assets/iconfont/iconfont.css';
</style>
