<template>
  <div class="page-404">
    <div class="inner">
      <div class="code">404</div>
      <p class="title">抱歉，页面走丢了</p>
      <p class="desc">您访问的页面不存在或已被移动</p>
      <div class="actions">
        <button class="btn" @click="goHome">返回首页</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()
function goHome() {
  router.replace({ path: '/' })
}
</script>

<style lang="scss" scoped>
.page-404 {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: calc(env(safe-area-inset-top, 0px) + 24px) 16px
    calc(env(safe-area-inset-bottom, 0px) + 32px);
  background: radial-gradient(1200px 600px at 50% -10%, #2f3542 0%, #1e222a 60%, #12161c 100%);
  color: #ffffff;
  text-align: center;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;

  .inner {
    width: 100%;
    max-width: 720px;
    margin: 0 auto;
  }

  .code {
    font-size: clamp(96px, 18vw, 160px);
    font-weight: 800;
    line-height: 1;
    letter-spacing: 0.04em;
    background: linear-gradient(180deg, #ffffff 0%, #cbd5e1 60%, #94a3b8 100%);
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
    text-shadow: 0 12px 30px rgba(0, 0, 0, 0.35);
  }

  .title {
    margin-top: 16px;
    font-size: clamp(18px, 4.6vw, 28px);
    font-weight: 700;
  }

  .desc {
    margin-top: 8px;
    font-size: clamp(14px, 3.6vw, 16px);
    color: #cbd5e1;
  }

  .actions {
    margin-top: 24px;
  }

  .btn {
    appearance: none;
    border: 0;
    border-radius: 9999px;
    padding: 12px 20px;
    font-weight: 600;
    font-size: 14px;
    color: #ffffff;
    background: #4f46e5;
    box-shadow: 0 10px 20px rgba(79, 70, 229, 0.35), inset 0 0 0 1px rgba(255, 255, 255, 0.15);
    cursor: pointer;
    transition: transform 0.12s ease, box-shadow 0.2s ease, background 0.2s ease;
  }

  .btn:hover {
    background: #4338ca;
    box-shadow: 0 12px 24px rgba(79, 70, 229, 0.45), inset 0 0 0 1px rgba(255, 255, 255, 0.22);
  }

  .btn:active {
    transform: translateY(1px);
  }
}
</style>
