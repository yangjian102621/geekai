<template>
  <el-container class="realtime-conversation" :style="{height: height}">
    <!-- connection animation -->
    <el-container class="connection-container" v-if="!isConnected">
      <div class="phone-container">
        <div class="signal"></div>
        <div class="signal"></div>
        <div class="signal"></div>
        <div class="phone"></div>
      </div>
      <div class="status-text">{{ connectingText }}</div>
      <audio ref="backgroundAudio" loop>
        <source src="/medias/calling.mp3" type="audio/mp3" />
        您的浏览器不支持音频元素。
      </audio>
      <audio ref="hangUpAudio">
        <source src="/medias/hang-up.mp3" type="audio/mp3" />
        您的浏览器不支持音频元素。
      </audio>
    </el-container>

    <!-- conversation body -->
    <div class="conversation-container" v-else>
      <div class="wave-container">
        <div class="wave-animation">
          <div v-for="i in 5" :key="i" class="wave-ellipse"></div>
        </div>
      </div>
      <!-- 其余部分保持不变 -->
      <div class="voice-indicators">
        <div class="voice-indicator left">
          <canvas ref="clientCanvasRef"></canvas>
        </div>
        <div class="voice-indicator right">
          <canvas ref="serverCanvasRef"></canvas>
        </div>
      </div>
      <div class="call-controls">
        <el-tooltip content="长按发送语音" placement="top" effect="light">
          <ripple-button>
            <button class="call-button answer" @mousedown="startRecording" @mouseup="stopRecording">
              <i class="iconfont icon-mic-bold"></i>
            </button>
          </ripple-button>
        </el-tooltip>
        <el-tooltip content="结束通话" placement="top" effect="light">
          <button class="call-button hangup" @click="hangUp">
            <i class="iconfont icon-hung-up"></i>
          </button>
        </el-tooltip>
      </div>
    </div>
  </el-container>

</template>

<script setup>
import RippleButton from "@/components/ui/RippleButton.vue";
import { ref, onMounted, onUnmounted } from 'vue';
import { RealtimeClient } from '@openai/realtime-api-beta';
import { WavRecorder, WavStreamPlayer } from '@/lib/wavtools/index.js';
import { instructions } from '@/utils/conversation_config.js';
import { WavRenderer } from '@/utils/wav_renderer';
import {showMessageError} from "@/utils/dialog";
import {getUserToken} from "@/store/session";

// eslint-disable-next-line no-unused-vars,no-undef
const props = defineProps({
  height: {
    type: String,
    default: '100vh'
  }
})
// eslint-disable-next-line no-undef
const emits = defineEmits(['close']);

/********************** connection animation code *************************/
const fullText = "正在接通中...";
const connectingText = ref("")
let index = 0;
const typeText = () => {
  if (index < fullText.length) {
    connectingText.value += fullText[index];
    index++;
    setTimeout(typeText, 200); // 每300毫秒显示一个字
  } else {
    setTimeout(() => {
      connectingText.value = '';
      index = 0;
      typeText();
    }, 1000); // 等待1秒后重新开始
  }
}
/*************************** end of code ****************************************/

/********************** conversation process code ***************************/
const leftVoiceActive = ref(false);
const rightVoiceActive = ref(false);

const animateVoice = () => {
  leftVoiceActive.value = Math.random() > 0.5;
  rightVoiceActive.value = Math.random() > 0.5;
};



const wavRecorder = ref(new WavRecorder({ sampleRate: 24000 }));
const wavStreamPlayer = ref(new WavStreamPlayer({ sampleRate: 24000 }));
let host = process.env.VUE_APP_WS_HOST
if (host === '') {
  if (location.protocol === 'https:') {
    host = 'wss://' + location.host;
  } else {
    host = 'ws://' + location.host;
  }
}
const client = ref(
    new RealtimeClient({
      url: `${host}/api/realtime`,
      apiKey: getUserToken(),
      dangerouslyAllowAPIKeyInBrowser: true,
    })
);
// // Set up client instructions and transcription
client.value.updateSession({
  instructions: instructions,
  turn_detection: null,
  input_audio_transcription: { model: 'whisper-1' },
  voice: 'alloy',
});

// set voice wave canvas
const clientCanvasRef = ref(null);
const serverCanvasRef = ref(null);
const isConnected = ref(false);
const isRecording = ref(false);
const backgroundAudio = ref(null);
const hangUpAudio = ref(null);
function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}
const connect = async () => {
  if (isConnected.value) {
    return
  }
  // 播放背景音乐
  if (backgroundAudio.value) {
    backgroundAudio.value.play().catch(error => {
      console.error('播放失败，可能是浏览器的自动播放策略导致的:', error);
    });
  }
  // 模拟拨号延时
  await sleep(3000)
  try {
    await client.value.connect();
    await wavRecorder.value.begin();
    await wavStreamPlayer.value.connect();
    console.log("对话连接成功！")
    if (!client.value.isConnected()) {
      return
    }

    isConnected.value = true;
    backgroundAudio.value?.pause()
    backgroundAudio.value.currentTime = 0
    client.value.sendUserMessageContent([
      {
        type: 'input_text',
        text: '你好，我是极客学长!',
      },
    ]);
    if (client.value.getTurnDetectionType() === 'server_vad') {
      await wavRecorder.value.record((data) => client.value.appendInputAudio(data.mono));
    }
  } catch (e) {
    console.error(e)
  }
};

// 开始语音输入
const startRecording = async () => {
  if (isRecording.value) {
    return
  }

  isRecording.value = true;
  try {
   const trackSampleOffset = await wavStreamPlayer.value.interrupt();
   if (trackSampleOffset?.trackId) {
     const { trackId, offset } = trackSampleOffset;
     client.value.cancelResponse(trackId, offset);
   }
   await wavRecorder.value.record((data) => client.value.appendInputAudio(data.mono));
  } catch (e) {
   console.error(e)
  }
};

// 结束语音输入
const stopRecording = async () => {
  try {
    isRecording.value = false;
    await wavRecorder.value.pause();
    client.value.createResponse();
  } catch (e) {
    console.error(e)
  }
};

// const changeTurnEndType = async (value) => {
//   if (value === 'none' && wavRecorder.value.getStatus() === 'recording') {
//     await wavRecorder.value.pause();
//   }
//   client.value.updateSession({
//     turn_detection: value === 'none' ? null : { type: 'server_vad' },
//   });
//   if (value === 'server_vad' && client.value.isConnected()) {
//     await wavRecorder.value.record((data) => client.value.appendInputAudio(data.mono));
//   }
//   canPushToTalk.value = value === 'none';
// };

// 初始化 WaveRecorder 组件和 RealtimeClient 事件处理
const initialize = async () => {
  // Set up render loops for the visualization canvas
  let isLoaded = true;
  const render = () => {
    if (isLoaded) {
      if (clientCanvasRef.value) {
        const canvas = clientCanvasRef.value;
        if (!canvas.width || !canvas.height) {
          canvas.width = canvas.offsetWidth;
          canvas.height = canvas.offsetHeight;
        }
        const ctx = canvas.getContext('2d');
        if (ctx) {
          ctx.clearRect(0, 0, canvas.width, canvas.height);
          const result = wavRecorder.value.recording
              ? wavRecorder.value.getFrequencies('voice')
              : { values: new Float32Array([0]) };
          WavRenderer.drawBars(canvas, ctx, result.values, '#0099ff', 10, 0, 8);
        }
      }
      if (serverCanvasRef.value) {
        const canvas = serverCanvasRef.value;
        if (!canvas.width || !canvas.height) {
          canvas.width = canvas.offsetWidth;
          canvas.height = canvas.offsetHeight;
        }
        const ctx = canvas.getContext('2d');
        if (ctx) {
          ctx.clearRect(0, 0, canvas.width, canvas.height);
          const result = wavStreamPlayer.value.analyser
              ?  wavStreamPlayer.value.getFrequencies('voice')
              : { values: new Float32Array([0]) };
          WavRenderer.drawBars(canvas, ctx, result.values, '#009900', 10, 0, 8);
        }
      }
      requestAnimationFrame(render);
    }
  };
  render();

  client.value.on('error', (event) => {
    showMessageError(event.error)
  });

  client.value.on('realtime.event', (re) => {
    if (re.event.type === 'error') {
      showMessageError(re.event.error)
    }
  });

  client.value.on('conversation.interrupted', async () => {
    const trackSampleOffset = await wavStreamPlayer.value.interrupt();
    if (trackSampleOffset?.trackId) {
      const { trackId, offset } = trackSampleOffset;
      client.value.cancelResponse(trackId, offset);
    }
  });

  client.value.on('conversation.updated', async ({ item, delta }) => {
    // console.log('item updated', item, delta)
    if (delta?.audio) {
      wavStreamPlayer.value.add16BitPCM(delta.audio, item.id);
    }
  });

}

const voiceInterval = ref(null);
onMounted(() => {
  initialize()
  // 启动聊天进行中的动画
  voiceInterval.value = setInterval(animateVoice, 200);
  typeText()
});

onUnmounted(() => {
  clearInterval(voiceInterval.value);
  client.value.reset();
});

// 挂断通话
const hangUp = async () => {
  try {
    isConnected.value = false
    // 停止播放拨号音乐
    if (backgroundAudio.value?.currentTime) {
      backgroundAudio.value?.pause()
      backgroundAudio.value.currentTime = 0
    }
    // 断开客户端的连接
    client.value.reset()
    // 中断语音输入和输出服务
    await wavRecorder.value.end()
    await wavStreamPlayer.value.interrupt()
  } catch (e) {
    console.error(e)
  } finally {
    // 播放挂断音乐
    hangUpAudio.value?.play()
    emits('close')
  }
};

// eslint-disable-next-line no-undef
defineExpose({ connect,hangUp });
</script>

<style scoped lang="stylus">

@import "@/assets/css/realtime.styl"

</style>