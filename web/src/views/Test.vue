<template>
  <div class="video-call-container">
    <div class="wave-animation">
      <div v-for="i in 5" :key="i" class="wave-ellipse"></div>
    </div>
    <!-- 其余部分保持不变 -->
    <div class="voice-indicators">
      <div class="voice-indicator left">
        <div v-for="i in 3" :key="i" class="bar"></div>
      </div>
      <div class="voice-indicator right">
        <canvas id="canvas" width="600" height="200"></canvas>
      </div>
    </div>
    <div class="call-controls">
      <button class="call-button hangup" @click="hangUp">
        <span class="icon">×</span>
      </button>
      <button class="call-button answer" @click="answer">
        <span class="icon">□</span>
      </button>
    </div>
  </div>
</template>

<script setup>
// Script 部分保持不变
import {ref, onMounted, onUnmounted} from 'vue';

const leftVoiceActive = ref(false);
const rightVoiceActive = ref(false);

const animateVoice = () => {
  leftVoiceActive.value = Math.random() > 0.5;
  rightVoiceActive.value = Math.random() > 0.5;
};

let voiceInterval;

onMounted(() => {
  voiceInterval = setInterval(animateVoice, 500);

  const flag = ref(false)
  const canvas = document.getElementById('canvas');
  const ctx = canvas.getContext('2d');

  async function setupAudioProcessing() {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      const audioContext = new (window.AudioContext || window.webkitAudioContext)();
      const analyser = audioContext.createAnalyser();
      const source = audioContext.createMediaStreamSource(stream);
      source.connect(analyser);
      analyser.fftSize = 256;
      const bufferLength = analyser.frequencyBinCount;
      const dataArray = new Uint8Array(bufferLength);

      function draw() {
        analyser.getByteFrequencyData(dataArray);

       if (!flag.value) {
         // ctx.clearRect(0, 0, canvas.width, canvas.height);
         // requestAnimationFrame(draw);
         return;
       }

        ctx.clearRect(0, 0, canvas.width, canvas.height);

        const barWidth = (canvas.width / bufferLength) * 5;
        let x = 0;

        for (let i = 0; i < bufferLength; i++) {
          const barHeight = dataArray[i] / 2;

          ctx.fillStyle = 'rgba(173, 216, 230, 0.7)'; // 淡蓝色
          ctx.fillRect(x, canvas.height - barHeight, barWidth, barHeight);

          x += barWidth + 2;
        }

        requestAnimationFrame(draw);
      }

      draw();
    } catch (err) {
      console.error('获取麦克风权限失败:', err);
    }
  }

  function setupSpeechRecognition() {
    const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
    if (!SpeechRecognition) {
      alert('您的浏览器不支持语音识别功能');
      return;
    }

    const recognition = new SpeechRecognition();
    recognition.continuous = true;
    recognition.interimResults = true;
    recognition.lang = 'zh-CN';

    recognition.onresult = (event) => {
      const transcript = Array.from(event.results)
          .map(result => result[0])
          .map(result => result.transcript)
          .join('');

      document.getElementById('transcript').textContent = '识别到的文本: ' + transcript;
      flag.value = true
    };

    recognition.onerror = (event) => {
      console.error('语音识别错误：', event.error);
      flag.value = false
    };

    recognition.start();
  }


  setupAudioProcessing();
  setupSpeechRecognition()
});

onUnmounted(() => {
  clearInterval(voiceInterval);
});

const hangUp = () => {
  console.log('Call hung up');
};

const answer = () => {
  console.log('Call answered');
};
</script>

<style scoped>
.video-call-container {
  background: linear-gradient(to right, #2c3e50, #4a5568, #6b46c1);
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  padding: 2rem;
}

canvas {
  background-color: transparent;
}

.wave-animation {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
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

/* 其余样式保持不变 */
.voice-indicators {
  display: flex;
  justify-content: space-between;
  width: 100%;
}

.voice-indicator {
  display: flex;
  align-items: flex-end;
}

.bar {
  width: 10px;
  height: 20px;
  background-color: #3498db;
  margin: 0 2px;
  transition: height 0.2s ease;
}

.voice-indicator.left .bar:nth-child(1) {
  height: 15px;
}

.voice-indicator.left .bar:nth-child(2) {
  height: 25px;
}

.voice-indicator.left .bar:nth-child(3) {
  height: 20px;
}

.voice-indicator.right .bar:nth-child(1) {
  height: 20px;
}

.voice-indicator.right .bar:nth-child(2) {
  height: 10px;
}

.voice-indicator.right .bar:nth-child(3) {
  height: 30px;
}

.call-controls {
  display: flex;
  justify-content: center;
  gap: 2rem;
}

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
</style>