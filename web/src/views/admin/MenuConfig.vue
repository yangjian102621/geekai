<template>
  <div class="menu-config form">
    <div class="container">
      <h3>菜单配置</h3>
      <Menu />
    </div>
  </div>
</template>

<script setup>
import Menu from '@/views/admin/Menu.vue'
</script>

<style scoped>
.menu-config {
  display: flex;
  justify-content: center;
}

.container {
  width: 100%;
  background-color: var(--el-bg-color);
  padding: 10px 20px 40px 20px;
}
</style>
