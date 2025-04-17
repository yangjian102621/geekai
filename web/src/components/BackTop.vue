<template>
  <button
    v-if="showButton"
    @click="scrollToTop"
    class="scroll-to-top"
    :style="{
      bottom: bottom + 'px',
      right: right + 'px',
      backgroundColor: bgColor
    }"
  >
    <el-icon><ArrowUpBold /></el-icon>
  </button>
</template>

<script>
import { ArrowUpBold } from "@element-plus/icons-vue";

export default {
  name: "BackTop",
  components: { ArrowUpBold },
  props: {
    bottom: {
      type: Number,
      default: 155
    },
    right: {
      type: Number,
      default: 30
    },
    bgColor: {
      type: String,
      default: "#b6aaf9"
    }
  },
  data() {
    return {
      showButton: false
    };
  },
  mounted() {
    this.checkScroll();
    window.addEventListener("resize", this.checkScroll);
    this.$el.parentElement.addEventListener("scroll", this.checkScroll);
  },
  beforeUnmount() {
    window.removeEventListener("resize", this.checkScroll);
    this.$el.parentElement.removeEventListener("scroll", this.checkScroll);
  },
  methods: {
    scrollToTop() {
      const container = this.$el.parentElement;
      container.scrollTo({
        top: 0,
        behavior: "smooth"
      });
    },
    checkScroll() {
      const container = this.$el.parentElement;
      this.showButton = container.scrollTop > 50;
    }
  }
};
</script>

<style scoped lang="stylus">
.scroll-to-top {
  position: fixed;
  color: white;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  outline: none;
  transition: opacity 0.3s;
  width 30px
  height 30px
  display flex
  justify-content center
  align-items center
  font-size 18px

  &:hover {
    opacity: 0.6;
  }
}
</style>
