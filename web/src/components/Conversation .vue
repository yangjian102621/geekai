<template>
  <!--语音通话组件-->
  <div class="video-call-container" :style="{height: height}">
    <div class="wave-container">
      <div class="wave-animation">
        <div v-for="i in 5" :key="i" class="wave-ellipse"></div>
      </div>
    </div>
    <!-- 其余部分保持不变 -->
    <div class="voice-indicators">
      <div class="voice-indicator left">
        <canvas ref="canvasClientRef" width="600" height="200"></canvas>
      </div>
      <div class="voice-indicator right">
        <canvas ref="canvasServerRef" width="600" height="200"></canvas>
      </div>
    </div>
    <div class="call-controls">
      <button class="call-button hangup" @click="hangUp">
        <i class="iconfont icon-hung-up"></i>
      </button>
    </div>
  </div>
</template>

<script setup>
import {onMounted, onUnmounted, ref} from "vue";

const leftVoiceActive = ref(false);
const rightVoiceActive = ref(false);
const props = defineProps({
  height: {
    type: String,
    default: '100vh'
  }
})
const emits = defineEmits(['hangUp']);

const animateVoice = () => {
  leftVoiceActive.value = Math.random() > 0.5;
  rightVoiceActive.value = Math.random() > 0.5;
};

let voiceInterval;

onMounted(() => {
  voiceInterval = setInterval(animateVoice, 500);
});

onUnmounted(() => {
  clearInterval(voiceInterval);
});

const hangUp = () => {
  console.log('Call hung up');
  emits('hangUp')
};

</script>

<style scoped lang="stylus">

.video-call-container {
  background: linear-gradient(to right, #2c3e50, #4a5568, #6b46c1);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  padding: 0;
  width 100vw

  .wave-container {
    padding 3rem
    .wave-animation {
      display: flex;
      justify-content: center;
      align-items: center;
      gap: 10px;
    }
  }


  .wave-ellipse {
    width: 40px;
    height: 40px;
    background-color: white;
    border-radius: 20px;
    animation: wave 0.8s infinite ease-in-out;
  }

  .wave-ellipse:nth-child(odd) {
    height: 60px;
  }

  .wave-ellipse:nth-child(even) {
    height: 80px;
  }

  @keyframes wave {
    0%, 100% {
      transform: scaleY(0.8);
    }
    50% {
      transform: scaleY(1.2);
    }
  }

  .wave-ellipse:nth-child(2) {
    animation-delay: 0.1s;
  }

  .wave-ellipse:nth-child(3) {
    animation-delay: 0.2s;
  }

  .wave-ellipse:nth-child(4) {
    animation-delay: 0.3s;
  }

  .wave-ellipse:nth-child(5) {
    animation-delay: 0.4s;
  }

  .call-controls {
    display: flex;
    justify-content: center;
    gap: 3rem;
    padding 3rem

    .call-button {
      width: 60px;
      height: 60px;
      border-radius: 50%;
      border: none;
      display: flex;
      justify-content: center;
      align-items: center;
      font-size: 24px;
      color: white;
      cursor: pointer;

      .iconfont {
        font-size 24px
      }
    }
    .hangup {
      background-color: #e74c3c;
    }

    .answer {
      background-color: #2ecc71;
    }

    .icon {
      font-size: 28px;
    }
  }

}

</style>