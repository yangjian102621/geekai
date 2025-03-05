<template>
  <div data-component="ConsolePage">
    <div class="content-top">
      <div class="content-title">
        <img src="/openai-logomark.svg" alt="OpenAI Logo" />
        <span>realtime console</span>
      </div>

    </div>
    <div class="content-main">
      <div class="content-logs">
        <div class="content-block events">
          <div class="visualization">
            <div class="visualization-entry client">
              <canvas ref="clientCanvasRef" />
            </div>
            <div class="visualization-entry server">
              <canvas ref="serverCanvasRef" />
            </div>
          </div>
          <div class="content-block-title">events</div>
          <div class="content-block-body" ref="eventsScrollRef">
            <template v-if="!realtimeEvents.length">
              awaiting connection...
            </template>
            <template v-else>
              <div v-for="(realtimeEvent, i) in realtimeEvents" :key="realtimeEvent.event.event_id" class="event">
                <div class="event-timestamp">
                  {{ formatTime(realtimeEvent.time) }}
                </div>
                <div class="event-details">
                  <div
                      class="event-summary"
                      @click="toggleEventDetails(realtimeEvent.event.event_id)"
                  >
                    <div
                        :class="[
                        'event-source',
                        realtimeEvent.event.type === 'error'
                          ? 'error'
                          : realtimeEvent.source,
                      ]"
                    >
                      <component :is="realtimeEvent.source === 'client' ? ArrowUp : ArrowDown" />
                      <span>
                        {{ realtimeEvent.event.type === 'error'
                          ? 'error!'
                          : realtimeEvent.source }}
                      </span>
                    </div>
                    <div class="event-type">
                      {{ realtimeEvent.event.type }}
                      {{ realtimeEvent.count ? `(${realtimeEvent.count})` : '' }}
                    </div>
                  </div>
                  <div
                      v-if="expandedEvents[realtimeEvent.event.event_id]"
                      class="event-payload"
                  >
                    {{ JSON.stringify(realtimeEvent.event, null, 2) }}
                  </div>
                </div>
              </div>
            </template>
          </div>
        </div>
        <div class="content-block conversation">
          <div class="content-block-title">conversation</div>
          <div class="content-block-body" data-conversation-content>
            <template v-if="!items.length">
              awaiting connection...
            </template>
            <template v-else>
              <div
                  v-for="(conversationItem, i) in items"
                  :key="conversationItem.id"
                  class="conversation-item"
              >
                <div :class="['speaker', conversationItem.role || '']">
                  <div>
                    {{
                      (conversationItem.role || conversationItem.type).replaceAll(
                          '_',
                          ' '
                      )
                    }}
                  </div>
                  <div class="close" @click="deleteConversationItem(conversationItem.id)">
                    <X />
                  </div>
                </div>
                <div class="speaker-content">
                  <!-- tool response -->
                  <div v-if="conversationItem.type === 'function_call_output'">
                    {{ conversationItem.formatted.output }}
                  </div>
                  <!-- tool call -->
                  <div v-if="conversationItem.formatted.tool">
                    {{ conversationItem.formatted.tool.name }}(
                    {{ conversationItem.formatted.tool.arguments }})
                  </div>
                  <div
                      v-if="
                      !conversationItem.formatted.tool &&
                      conversationItem.role === 'user'
                    "
                  >
                    {{
                      conversationItem.formatted.transcript ||
                      (conversationItem.formatted.audio?.length
                          ? '(awaiting transcript)'
                          : conversationItem.formatted.text || '(item sent)')
                    }}
                  </div>
                  <div
                      v-if="
                      !conversationItem.formatted.tool &&
                      conversationItem.role === 'assistant'
                    "
                  >
                    {{
                      conversationItem.formatted.transcript ||
                      conversationItem.formatted.text ||
                      '(truncated)'
                    }}
                  </div>
                  <audio
                      v-if="conversationItem.formatted.file"
                      :src="conversationItem.formatted.file.url"
                      controls
                  />
                </div>
              </div>
            </template>
          </div>
        </div>
        <div class="content-actions" style="position:absolute; top: 0; left: 0">
          <el-button
              :type="isConnected ? '' : 'primary'"
              @click="connectConversation"
          >
            {{isConnected ? '断开连接' : '连接对话'}}
          </el-button>

          <el-button @mousedown="startRecording" @mouseup="stopRecording">开始讲话</el-button>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, watch } from 'vue';
import { RealtimeClient } from '@openai/realtime-api-beta';
import { WavRecorder, WavStreamPlayer } from '@/lib/wavtools/index.js';
import { instructions } from '@/utils/conversation_config.js';
import { WavRenderer } from '@/utils/wav_renderer';

// Constants
const LOCAL_RELAY_SERVER_URL = process.env.REACT_APP_LOCAL_RELAY_SERVER_URL || '';

// Reactive state
const apiKey = ref(
    LOCAL_RELAY_SERVER_URL
        ? ''
        : localStorage.getItem('tmp::voice_api_key') || prompt('OpenAI API Key') || ''
);
const wavRecorder = ref(new WavRecorder({ sampleRate: 24000 }));
const wavStreamPlayer = ref(new WavStreamPlayer({ sampleRate: 24000 }));
const client = ref(
    new RealtimeClient({
      url: "ws://localhost:5678/api/realtime",
      apiKey: "sk-Gc5cEzDzGQLIqxWA9d62089350F3454bB359C4A3Fa21B3E4",
      dangerouslyAllowAPIKeyInBrowser: true,
    })
);

const clientCanvasRef = ref(null);
const serverCanvasRef = ref(null);
const eventsScrollRef = ref(null);
const startTime = ref(new Date().toISOString());

const items = ref([]);
const realtimeEvents = ref([]);
const expandedEvents = reactive({});
const isConnected = ref(false);
const canPushToTalk = ref(true);
const isRecording = ref(false);
const memoryKv = ref({});
const coords = ref({ lat: 37.775593, lng: -122.418137 });
const marker = ref(null);

// Methods
const formatTime = (timestamp) => {
  const t0 = new Date(startTime.value).valueOf();
  const t1 = new Date(timestamp).valueOf();
  const delta = t1 - t0;
  const hs = Math.floor(delta / 10) % 100;
  const s = Math.floor(delta / 1000) % 60;
  const m = Math.floor(delta / 60_000) % 60;
  const pad = (n) => {
    let s = n + '';
    while (s.length < 2) {
      s = '0' + s;
    }
    return s;
  };
  return `${pad(m)}:${pad(s)}.${pad(hs)}`;
};

const connectConversation = async () => {
  alert(123)
  startTime.value = new Date().toISOString();
  isConnected.value = true;
  realtimeEvents.value = [];
  items.value = client.value.conversation.getItems();

  await wavRecorder.value.begin();
  await wavStreamPlayer.value.connect();
  await client.value.connect();
  client.value.sendUserMessageContent([
    {
      type: 'input_text',
      text: '你好，我是老阳!',
    },
  ]);

  if (client.value.getTurnDetectionType() === 'server_vad') {
    await wavRecorder.value.record((data) => client.value.appendInputAudio(data.mono));
  }
};

const disconnectConversation = async () => {
  isConnected.value = false;
  realtimeEvents.value = [];
  items.value = [];
  memoryKv.value = {};
  coords.value = { lat: 37.775593, lng: -122.418137 };
  marker.value = null;

  client.value.disconnect();
  await wavRecorder.value.end();
  await wavStreamPlayer.value.interrupt();
};

const deleteConversationItem = async (id) => {
  client.value.deleteItem(id);
};

const startRecording = async () => {
  isRecording.value = true;
  const trackSampleOffset = await wavStreamPlayer.value.interrupt();
  if (trackSampleOffset?.trackId) {
    const { trackId, offset } = trackSampleOffset;
    await client.value.cancelResponse(trackId, offset);
  }
  await wavRecorder.value.record((data) => client.value.appendInputAudio(data.mono));
};

const stopRecording = async () => {
  isRecording.value = false;
  await wavRecorder.value.pause();
  client.value.createResponse();
};

const changeTurnEndType = async (value) => {
  if (value === 'none' && wavRecorder.value.getStatus() === 'recording') {
    await wavRecorder.value.pause();
  }
  client.value.updateSession({
    turn_detection: value === 'none' ? null : { type: 'server_vad' },
  });
  if (value === 'server_vad' && client.value.isConnected()) {
    await wavRecorder.value.record((data) => client.value.appendInputAudio(data.mono));
  }
  canPushToTalk.value = value === 'none';
};

const toggleEventDetails = (eventId) => {
  if (expandedEvents[eventId]) {
    delete expandedEvents[eventId];
  } else {
    expandedEvents[eventId] = true;
  }
};

// Lifecycle hooks and watchers
onMounted(() => {
  if (apiKey.value !== '') {
    localStorage.setItem('tmp::voice_api_key', apiKey.value);
  }

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
    realtimeEvents.value = realtimeEvents.value.slice();
    const lastEvent = realtimeEvents.value[realtimeEvents.value.length - 1];
    if (lastEvent?.event.type === realtimeEvent.event.type) {
      lastEvent.count = (lastEvent.count || 0) + 1;
      realtimeEvents.value.splice(-1, 1, lastEvent);
    } else {
      realtimeEvents.value.push(realtimeEvent);
    }
  });

  client.value.on('error', (event) => console.error(event));

  client.value.on('conversation.interrupted', async () => {
    const trackSampleOffset = await wavStreamPlayer.value.interrupt();
    if (trackSampleOffset?.trackId) {
      const { trackId, offset } = trackSampleOffset;
      await client.value.cancelResponse(trackId, offset);
    }
  });

  client.value.on('conversation.updated', async ({ item, delta }) => {
    items.value = client.value.conversation.getItems();
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

  // Set up client instructions and tools
  client.value.updateSession({ instructions: instructions });
  client.value.updateSession({ input_audio_transcription: { model: 'whisper-1' } });

  client.value.addTool(
      {
        name: 'set_memory',
        description: 'Saves important data about the user into memory.',
        parameters: {
          type: 'object',
          properties: {
            key: {
              type: 'string',
              description:
                  'The key of the memory value. Always use lowercase and underscores, no other characters.',
            },
            value: {
              type: 'string',
              description: 'Value can be anything represented as a string',
            },
          },
          required: ['key', 'value'],
        },
      },
      async ({ key, value }) => {
        memoryKv.value = { ...memoryKv.value, [key]: value };
        return { ok: true };
      }
  );

  client.value.addTool(
      {
        name: 'get_weather',
        description:
            'Retrieves the weather for a given lat, lng coordinate pair. Specify a label for the location.',
        parameters: {
          type: 'object',
          properties: {
            lat: {
              type: 'number',
              description: 'Latitude',
            },
            lng: {
              type: 'number',
              description: 'Longitude',
            },
            location: {
              type: 'string',
              description: 'Name of the location',
            },
          },
          required: ['lat', 'lng', 'location'],
        },
      },
      async ({ lat, lng, location }) => {
        marker.value = { lat, lng, location };
        coords.value = { lat, lng, location };
        const result = await fetch(
            `https://api.open-meteo.com/v1/forecast?latitude=${lat}&longitude=${lng}&current=temperature_2m,wind_speed_10m`
        );
        const json = await result.json();
        const temperature = {
          value: json.current.temperature_2m,
          units: json.current_units.temperature_2m,
        };
        const wind_speed = {
          value: json.current.wind_speed_10m,
          units: json.current_units.wind_speed_10m,
        };
        marker.value = { lat, lng, location, temperature, wind_speed };
        return json;
      }
  );

  items.value = client.value.conversation.getItems();
});

onUnmounted(() => {
  client.value.reset();
});

// Watchers
watch(realtimeEvents, () => {
  if (eventsScrollRef.value) {
    const eventsEl = eventsScrollRef.value;
    eventsEl.scrollTop = eventsEl.scrollHeight;
  }
});

watch(items, () => {
  const conversationEls = document.querySelectorAll('[data-conversation-content]');
  conversationEls.forEach((el) => {
    el.scrollTop = el.scrollHeight;
  });
});
</script>

<style scoped>
/* You can add your component-specific styles here */
/* If you're using SCSS, you might want to import your existing SCSS file */
/* @import './ConsolePage.scss'; */
</style>