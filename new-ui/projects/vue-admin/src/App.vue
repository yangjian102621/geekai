<template>
  <RouterView />
</template>
<style>
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-thumb {
  border-radius: 1px;
  box-shadow: inset 0 0 5px #0000000d;
  background: #d9d9d9;
}

::-webkit-scrollbar-track {
  box-shadow: inset 0 0 5px #0000000d;
  border-radius: 1px;
  background: #fafafa;
}

.public-bg {
  display: flex;
  width: 100vw;
  height: 100vh;
  align-items: center;
  justify-content: center;
  background: linear-gradient(133deg, #ffffff 0%, #dde8fe 100%);
  overflow: hidden;
}
</style>
