<template>
  <button
      class="ripple-button"
      @mousedown="startRipples"
      @mouseup="stopRipples"
      @mouseleave="stopRipples"
  >
    <slot></slot>
    <span
        v-for="ripple in ripples"
        :key="ripple.id"
        class="ripple"
        :style="getRippleStyle(ripple)"
    ></span>
  </button>
</template>

<script setup>
import { ref } from 'vue';

const ripples = ref([]);
let rippleCount = 0;
let animationId;

const startRipples = (event) => {
  const button = event.currentTarget;
  const rect = button.getBoundingClientRect();
  const size = Math.max(rect.width, rect.height);
  // const x = event.clientX - rect.left;
  // const y = event.clientY - rect.top;
  const x = rect.right - rect.left - size/2;
  const y = rect.bottom - rect.top - size/2;

  const createRipple = () => {
    ripples.value.push({
      id: rippleCount++,
      x,
      y,
      size: 0,
      opacity: 0.5
    });

    if (ripples.value.length > 3) {
      ripples.value.shift();
    }
  };

  const animate = () => {
    ripples.value.forEach(ripple => {
      ripple.size += 2;
      ripple.opacity -= 0.01;
    });

    ripples.value = ripples.value.filter(ripple => ripple.opacity > 0);

    if (ripples.value.length < 3) {
      createRipple();
    }

    animationId = requestAnimationFrame(animate);
  };

  createRipple();
  animate();
};

const stopRipples = () => {
  cancelAnimationFrame(animationId);
  ripples.value = [];
};

const getRippleStyle = (ripple) => ({
  left: `${ripple.x}px`,
  top: `${ripple.y}px`,
  width: `${ripple.size}px`,
  height: `${ripple.size}px`,
  opacity: ripple.opacity
});
</script>

<style scoped lang="stylus">
.ripple-button {
  position: relative;
  overflow: hidden;
  border: none;
  background none;
  color: white;
  cursor: pointer;
  border-radius: 50%;
  outline: none;
  margin 0
  padding 0
}

.ripple {
  position: absolute;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.7);
  transform: translate(-50%, -50%);
  pointer-events: none;
}
</style>