<template>
  <div class="audio-chat-page">
    <el-button style="margin: 20px" type="primary" size="large" @click="connect">开始语音对话</el-button>

    <el-dialog v-model="showDialog" title="语音通话" :fullscreen="true">
      <el-container>
        <calling v-if="!connected" :height="dialogHeight+'px'" />
        <conversation v-else :height="dialogHeight+'px'" @hang-up="hangUp" />
      </el-container>
    </el-dialog>
  </div>
</template>

<script setup>
import {ref} from 'vue';
import { RealtimeClient } from '@openai/realtime-api-beta';
import Calling from "@/components/Calling.vue";
import Conversation from "@/components/RealtimeConversation.vue";
import {playPCM16} from "@/utils/wav_player";
import {showMessageError} from "@/utils/dialog";

const showDialog = ref(false);
const connected = ref(false);
const dialogHeight = ref(window.innerHeight - 75);

const recognition = ref(null)
if (!('webkitSpeechRecognition' in window)) {
  alert("你的浏览器不支持语音识别，请使用最新版本的 Chrome 浏览器。");
} else {
  recognition.value = new webkitSpeechRecognition();
  recognition.value.lang = 'zh-CN'; // 设置语言为简体中文
  recognition.value.continuous = false; // 设置为单句识别
  recognition.value.interimResults = false; // 不需要中间结果

  recognition.value.onresult = function(event) {
    const transcript = event.results[0][0].transcript;
    try {
      client.cancelResponse(chatId.value)
      speaker.value.stop()
    } catch (e) {
      console.warn(e)
    }
    console.log(`你说的是: ${transcript}`)
    console.log(client.isConnected())
   if (client.isConnected()){
     client.sendUserMessageContent([{ type: 'input_text', text: transcript }]);
   }
    //recognition.value.start()
  };

  recognition.value.onerror = function(event) {
    showMessageError("识别失败:", event.error)
  };
  recognition.value.onend = function() {
    console.log('语音识别结束，重新开始');
    recognition.value.start(); // 在结束时重新开始
  };

  //recognition.value.start()
}

const client = new RealtimeClient({
  url: "wss://api.geekai.pro/v1/realtime",
  apiKey: "sk-Gc5cEzDzGQLIqxWA9d62089350F3454bB359C4A3Fa21B3E4",
  dangerouslyAllowAPIKeyInBrowser: true,
});

// Can set parameters ahead of connecting, either separately or all at once
client.updateSession({ instructions: 'You are a great, upbeat friend.' });
client.updateSession({ voice: 'nova' });
client.updateSession({
  turn_detection: 'disabled', // or 'server_vad'
  input_audio_transcription: { model: 'whisper-1' },
});

const chatId = ref("")
const audioChunks = ref([])
// Set up event handling
client.on('conversation.updated', ({ item, delta }) => {
  chatId.value = item.id
  //console.info('conversation.updated', item, delta)
  switch (item.type) {
    case 'message':
      // system, user, or assistant message (item.role)
      localStorage.setItem("chat_data", JSON.stringify(Array.from(item.formatted.audio)))
      console.log(item)
      break;
    case 'function_call':
      // always a function call from the model
      break;
    case 'function_call_output':
      // always a response from the user / application
      break;
  }
  if (delta) {
    // console.info(delta.audio)
    if (delta.audio && delta.audio.length > 1) {
      audioChunks.value.push(delta.audio)
    }
    if (audioChunks.value.length === 1) {
      playAudio(0)
    }

    //localStorage.setItem("chat_data", JSON.stringify(Array.from(delta.audio)))
    // Only one of the following will be populated for any given event
    // delta.audio = Int16Array, audio added
    // delta.transcript = string, transcript added
    // delta.arguments = string, function arguments added
  }
});

const speaker = ref(null)
const playAudio = (index) => {
  if (index === 0 && speaker.value) {
    speaker.value.stop()
  }
  const data = audioChunks.value[index]
  console.log(data)
  if (index === audioChunks.value.length-1) {
    audioChunks.value = []
  }
  speaker.value = playPCM16(data, 24000);
  if (speaker.value !== null) {
    speaker.value.onended = () => {
      playAudio(index + 1)
    }
  }

}

client.on('conversation.interrupted', async () => {
  console.log('聊天中断')
});

client.on('conversation.item.appended', ({ item }) => {
  if (item.role === 'assistant') {
    // playPCM16(item.formatted.audio, 24000);
    // console.log(item)
  }
});

const connect = () => {
  showDialog.value = true
  client.connect().then(res => {
    if (res) {
      console.log("连接成功!")
      connected.value = true
      // const data = JSON.parse(localStorage.getItem("chat_data"))
      // playPCM16(data, 24000)
      client.sendUserMessageContent([{ type: 'input_text', text: `你好，我是老阳。` }]);
    }
  }).catch(e => {
    console.log(e)
  })
}

const hangUp = () => {
 try {
   client.cancelResponse(chatId.value)
   speaker.value.stop()
 } catch (e) {
   console.warn(e)
  }
  showDialog.value = false
  connected.value = false
}

</script>

<style scoped lang="stylus">
.audio-chat-page {
  display flex
  flex-flow column
  justify-content center
  align-items center

}

canvas {
  background-color: transparent;
}


</style>