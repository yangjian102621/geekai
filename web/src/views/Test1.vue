<template>
  <div>
   {{data}}
    <div>
      <canvas ref="clientCanvasRef" />
    </div>
    <el-button type="primary" @click="sendMessage">连接电话</el-button>
  </div>
</template>

<script setup>
import {nextTick, onMounted, ref} from "vue";
import Storage from "good-storage";

const data = ref('abc')
import { RealtimeClient } from '@openai/realtime-api-beta';

const client = new RealtimeClient({
  url: "wss://api.geekai.pro/v1/realtime",
  apiKey: "sk-Gc5cEzDzGQLIqxWA9d62089350F3454bB359C4A3Fa21B3E4",
  dangerouslyAllowAPIKeyInBrowser: true,
});

// Can set parameters ahead of connecting, either separately or all at once
client.updateSession({ instructions: 'You are a great, upbeat friend.' });
client.updateSession({ voice: 'alloy' });
client.updateSession({
  turn_detection: 'disabled', // or 'server_vad'
  input_audio_transcription: { model: 'whisper-1' },
});

// Set up event handling
client.on('conversation.updated', ({ item, delta }) => {
  console.info('conversation.updated', item, delta)
  switch (item.type) {
    case 'message':
      // system, user, or assistant message (item.role)
      localStorage.setItem("chat_data", JSON.stringify(Array.from(item.formatted.audio)))
        console.log("语言消息")
      break;
    case 'function_call':
      // always a function call from the model
      break;
    case 'function_call_output':
      // always a response from the user / application
      break;
  }
  if (delta) {
     console.info(delta.audio)
    //localStorage.setItem("chat_data", JSON.stringify(Array.from(delta.audio)))
    playPCM16(delta.audio, 24000);
    // Only one of the following will be populated for any given event
    // delta.audio = Int16Array, audio added
    // delta.transcript = string, transcript added
    // delta.arguments = string, function arguments added
  }
});

client.on('conversation.item.appended', ({ item }) => {
  if (item.role === 'assistant') {
    playPCM16(item.formatted.audio, 24000);
  }
});

const speaker = ref(null)
// 假设 PCM16 数据已经存储在一个 Int16Array 中
function playPCM16(pcm16Array, sampleRate = 44100) {
  // 创建 AudioContext
  const audioContext = new (window.AudioContext || window.webkitAudioContext)();

  // 将 Int16Array 转换为 Float32Array (Web Audio API 使用 Float32)
  let float32Array = new Float32Array(pcm16Array.length);
  for (let i = 0; i < pcm16Array.length; i++) {
    float32Array[i] = pcm16Array[i] / 32768; // Int16 转换为 Float32
  }

  // 创建 AudioBuffer
  const audioBuffer = audioContext.createBuffer(1, float32Array.length, sampleRate); // 单声道
  audioBuffer.getChannelData(0).set(float32Array); // 设置音频数据

  // 创建 AudioBufferSourceNode 并播放音频
  const source = audioContext.createBufferSource();
  source.buffer = audioBuffer;
  source.connect(audioContext.destination); // 连接到扬声器
  source.start(); // 播放
  speaker.value = source

}

onMounted(() => {
  // Connect to Realtime API
  // client.connect().then(res => {
  //   if (res) {
  //     console.log("连接成功!")
  //   }
  // }).catch(e => {
  //   console.log(e)
  // })


})

const sendMessage = () => {
  const data = JSON.parse(localStorage.getItem("chat_data"))
  playPCM16(data, 24000)
  setTimeout(() => {speaker.value.stop()}, 5000)
  // client.sendUserMessageContent([{ type: 'input_text', text: `你好，请用四川话给我讲一个笑话?` }]);
}

</script>
