<script lang="ts" setup>
import { IconDown, IconExport } from "@arco-design/web-vue/es/icon";

import SystemMenu from "./SystemMenu.vue";
import PageWrapper from "./PageWrapper.vue";

const logoWidth = "200px";
</script>
<template>
  <ALayout class="custom-layout">
    <ALayoutHeader class="custom-layout-header">
      <div class="logo"></div>
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
            <APopconfirm content="确认退出？" position="br">
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
      display: block;
      width: v-bind("logoWidth");
      text-align: center;
      cursor: pointer;
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
