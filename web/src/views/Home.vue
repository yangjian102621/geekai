<template>
  <div class="home">
    <div class="navigator"></div>
    <div class="content">
      <router-view v-slot="{ Component }">
        <transition name="move" mode="out-in">
          <component :is="Component"></component>
        </transition>
      </router-view>
    </div>
  </div>
</template>

<script setup>

// const router = useRouter();
// checkSession().then(() => {
//   if (isMobile()) {
//     router.push("/mobile")
//   } else {
//     router.push("/chat")
//   }
// }).catch(() => {
//   router.push("/login")
// })

</script>

<style lang="stylus" scoped>
.home {
  display: flex;
  background-color: #343540;
  height 100vh
  width 100%

  .navigator {
    display flex
    width 30px
    padding 10px 6px
  }

  .content {
    width: 100%;
    height: 100vh;
    overflow-y: scroll;
    box-sizing: border-box;
  }

}
</style>
