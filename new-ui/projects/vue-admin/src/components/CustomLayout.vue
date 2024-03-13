<script lang="ts" setup>
import { IconDown, IconExport } from "@arco-design/web-vue/es/icon";
import { useAuthStore } from "@/stores/auth";
import useState from "@/composables/useState";
import Logo from "/images/logo.png";
import avatar from "/images/user-info.jpg";
import donateImg from "/images/wechat-pay.png";

import SystemMenu from "./SystemMenu.vue";
import PageWrapper from "./PageWrapper.vue";

const logoWidth = "200px";
const authStore = useAuthStore();
const [visible, setVisible] = useState(false);
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
            <a-avatar class="user-avatar" :size="30">
              <img :src="avatar" />
            </a-avatar>
            <IconDown />
          </ASpace>
          <template #content>
            <a
              class="dropdown-link"
              href="https://github.com/yangjian102621/chatgpt-plus"
              target="_blank"
            >
              <ADoption value="1">
                <template #icon>
                  <icon-github />
                </template>
                <span>ChatPlus-AI 创作系统</span>
              </ADoption>
            </a>
            <ADoption value="2" @click="setVisible(true)">
              <template #icon>
                <icon-wechatpay />
              </template>
              <span>打赏作者</span>
            </ADoption>
          </template>
          <template #footer>
            <APopconfirm content="确认退出？" position="bl" @ok="authStore.logout">
              <AButton status="warning" class="logout-area">
                <ASpace align="center">
                  <IconExport size="16" />
                  <span>退出登录</span>
                </ASpace>
              </AButton>
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
  <a-modal
    v-model:visible="visible"
    class="donate-dialog"
    width="400px"
    title="请作者喝杯咖啡"
    :footer="false"
  >
    <a-alert :closable="false" :show-icon="false">
      如果你觉得这个项目对你有帮助，并且情况允许的话，可以请作者喝杯咖啡，非常感谢你的支持～
    </a-alert>
    <p>
      <a-image :src="donateImg" />
    </p>
  </a-modal>
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
.dropdown-link {
  text-decoration: none;
}
.donate-dialog {
  p {
    text-align: center;
  }
}
.logout-area {
  padding: 8px 0;
  display: flex;
  min-width: 80px;
  width: 100%;
  align-items: center;
  justify-content: center;
}
</style>
