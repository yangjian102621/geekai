<template>
  <div class="page-wrapper">
    <div class="page-content">
      <RouterView />
    </div>
  </div>
</template>
<style scoped>
.page-wrapper {
  height: calc(100vh - 60px);
  display: flex;
  flex-direction: column;
  background: #f7f8fa;
  overflow: hidden;
}
.page-content {
  margin: 12px;
  padding: 12px;
  flex: 1;
  background: #fff;
  overflow-y: auto;
}
</style>
