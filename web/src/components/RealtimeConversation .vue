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
    setTimeout(typeText, 300); // 每300毫秒显示一个字
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
const client = ref(
    new RealtimeClient({
      url: "ws://localhost:5678/api/realtime",
      apiKey: "sk-Gc5cEzDzGQLIqxWA9d62089350F3454bB359C4A3Fa21B3E4",
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
// const eventsScrollRef = ref(null);
// const startTime = ref(new Date().toISOString());

// const items = ref([]);
// const realtimeEvents = ref([]);
// const expandedEvents = reactive({});
const isConnected = ref(false);
// const canPushToTalk = ref(true);
const isRecording = ref(false);
// const memoryKv = ref({});
// const coords = ref({ lat: 37.775593, lng: -122.418137 });
// const marker = ref(null);

// Methods
// const formatTime = (timestamp) => {
//   const t0 = new Date(startTime.value).valueOf();
//   const t1 = new Date(timestamp).valueOf();
//   const delta = t1 - t0;
//   const hs = Math.floor(delta / 10) % 100;
//   const s = Math.floor(delta / 1000) % 60;
//   const m = Math.floor(delta / 60_000) % 60;
//   const pad = (n) => {
//     let s = n + '';
//     while (s.length < 2) {
//       s = '0' + s;
//     }
//     return s;
//   };
//   return `${pad(m)}:${pad(s)}.${pad(hs)}`;
// };

const connect = async () => {
  // startTime.value = new Date().toISOString();
  // realtimeEvents.value = [];
  // items.value = client.value.conversation.getItems();
  if (isConnected.value) {
    return
  }

  try {
    await client.value.connect();
    await wavRecorder.value.begin();
    await wavStreamPlayer.value.connect();
    isConnected.value = true;
    console.log("对话连接成功！")
    client.value.sendUserMessageContent([
      {
        type: 'input_text',
        text: '你好，我是老阳!',
      },
    ]);

    if (client.value.getTurnDetectionType() === 'server_vad') {
      await wavRecorder.value.record((data) => client.value.appendInputAudio(data.mono));
    }
  } catch (e) {
    showMessageError(e.message)
  }
};

// const disconnectConversation = async () => {
//   isConnected.value = false;
//   // realtimeEvents.value = [];
//   // items.value = [];
//   // memoryKv.value = {};
//   // coords.value = { lat: 37.775593, lng: -122.418137 };
//   // marker.value = null;
//
//   client.value.disconnect();
//   await wavRecorder.value.end();
//   await wavStreamPlayer.value.interrupt();
// };

// const deleteConversationItem = async (id) => {
//   client.value.deleteItem(id);
// };

const startRecording = async () => {
  isRecording.value = true;
  const trackSampleOffset = await wavStreamPlayer.value.interrupt();
  if (trackSampleOffset?.trackId) {
    const { trackId, offset } = trackSampleOffset;
    client.value.cancelResponse(trackId, offset);
  }
  await wavRecorder.value.record((data) => client.value.appendInputAudio(data.mono));
};

const stopRecording = async () => {
  isRecording.value = false;
  await wavRecorder.value.pause();
  client.value.createResponse();
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
//
// const toggleEventDetails = (eventId) => {
//   if (expandedEvents[eventId]) {
//     delete expandedEvents[eventId];
//   } else {
//     expandedEvents[eventId] = true;
//   }
// };

// Lifecycle hooks and watchers
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


  // Set up client event listeners
  client.value.on('realtime.event', (realtimeEvent) => {
    // realtimeEvents.value = realtimeEvents.value.slice();
    // const lastEvent = realtimeEvents.value[realtimeEvents.value.length - 1];
    // if (lastEvent?.event.type === realtimeEvent.event.type) {
    //   lastEvent.count = (lastEvent.count || 0) + 1;
    //   realtimeEvents.value.splice(-1, 1, lastEvent);
    // } else {
    //   realtimeEvents.value.push(realtimeEvent);
    // }
    // console.log(realtimeEvent)
  });

  client.value.on('error', (event) => console.error(event));

  client.value.on('conversation.interrupted', async () => {
    const trackSampleOffset = await wavStreamPlayer.value.interrupt();
    if (trackSampleOffset?.trackId) {
      const { trackId, offset } = trackSampleOffset;
      client.value.cancelResponse(trackId, offset);
    }
  });

  client.value.on('conversation.updated', async ({ item, delta }) => {
    console.log('item updated', item, delta)
    if (delta?.audio) {
      wavStreamPlayer.value.add16BitPCM(delta.audio, item.id);
    }
    if (item.status === 'completed' && item.formatted.audio?.length) {
      const wavFile = await WavRecorder.decode(
          item.formatted.audio,
          24000,
          24000
      );
      item.formatted.file = wavFile;
    }
  });

}

// Watchers
// watch(realtimeEvents, () => {
//   if (eventsScrollRef.value) {
//     const eventsEl = eventsScrollRef.value;
//     eventsEl.scrollTop = eventsEl.scrollHeight;
//   }
// });

// watch(items, () => {
//   const conversationEls = document.querySelectorAll('[data-conversation-content]');
//   conversationEls.forEach((el) => {
//     el.scrollTop = el.scrollHeight;
//   });
// });

const voiceInterval = ref(null);
onMounted(() => {
  initialize()
  voiceInterval.value = setInterval(animateVoice, 500);
  typeText()
});

onUnmounted(() => {
  clearInterval(voiceInterval.value);
  client.value.reset();
});

const hangUp = async () => {
  emits('close')
  isConnected.value = false;
  client.value.disconnect();
  await wavRecorder.value.end();
  await wavStreamPlayer.value.interrupt();
};

// eslint-disable-next-line no-undef
defineExpose({ connect });
</script>

<style scoped lang="stylus">

@import "@/assets/css/realtime.styl"

</style>