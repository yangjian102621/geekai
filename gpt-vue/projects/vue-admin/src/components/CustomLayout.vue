<script lang="ts" setup>
import { IconDown, IconExport } from "@arco-design/web-vue/es/icon";
import { useAuthStore } from "@/stores/auth";
import Logo from "/images/logo.png";

import SystemMenu from "./SystemMenu.vue";
import PageWrapper from "./PageWrapper.vue";

const logoWidth = "200px";
const authStore = useAuthStore();
</script>
<template>
  <ALayout class="custom-layout">
    <ALayoutHeader class="custom-layout-header">
      <div class="logo">
        <img :src="Logo" alt="logo" />
        <span>ChatPlus 控制台</span>
      </div>
      <div class="action">
        <ADropdown>
          <ASpace align="center" :size="4">
            <span></span>
            <IconDown />
          </ASpace>
          <template #content>
            <ADoption value="changeOwnPwd">更改密码</ADoption>
          </template>
          <template #footer>
            <APopconfirm
              content="确认退出？"
              position="br"
              @ok="authStore.logout"
            >
              <ASpace align="center" class="logout-area">
                <IconExport size="16" />
                <span>退出</span>
              </ASpace>
            </APopconfirm>
          </template>
        </ADropdown>
      </div>
    </ALayoutHeader>
    <ALayout>
      <SystemMenu :width="logoWidth" />
      <ALayoutContent>
        <PageWrapper>
          <slot />
        </PageWrapper>
      </ALayoutContent>
    </ALayout>
  </ALayout>
</template>
<style lang="less" scoped>
.custom-layout {
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  &-header {
    display: flex;
    width: 100%;
    height: 60px;
    align-items: center;
    border-bottom: 1px solid var(--color-neutral-2);
    .logo {
      display: flex;
      width: v-bind("logoWidth");
      align-items: center;
      justify-content: center;
      gap: 12px;
      img {
        width: 30px;
        height: 30px;
      }
    }
    .action {
      display: flex;
      padding: 0 12px;
      flex: 1;
      justify-content: right;
      align-items: center;
    }
  }
}
.logout-area {
  padding: 8px 0;
  display: flex;
  width: 80px;
  align-items: center;
  justify-content: center;
}
</style>
