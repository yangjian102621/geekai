<template>
 <div class="player">
   <div class="container">
     <div class="cover">
       <el-image :src="cover" fit="cover" />
     </div>
     <div class="info">
       <div class="title">{{title}}</div>
       <div class="style">
         <span class="tags">{{ tags }}</span>
         <span class="text-lightGray"> | </span>
         <span class="time">{{ formatTime(currentTime) }}<span class="split">/</span>{{ formatTime(duration) }}</span>
       </div>
     </div>

     <div class="controls-container">
       <div class="controls">
         <button @click="prevSong" class="control-btn">
           <i class="iconfont icon-prev"></i>
         </button>
         <button @click="togglePlay" class="control-btn">
           <i class="iconfont icon-play" v-if="!isPlaying"></i>
           <i class="iconfont icon-pause" v-else></i>
         </button>
         <button @click="nextSong" class="control-btn">
           <i class="iconfont icon-next"></i>
         </button>
       </div>
     </div>

     <div class="progress-bar" @click="setProgress" ref="progressBarRef">
       <div class="progress" :style="{ width: `${progressPercent}%` }"></div>
     </div>
     <audio ref="audio" @timeupdate="updateProgress" @ended="nextSong"></audio>

     <el-button class="close" type="info" :icon="Close" circle size="small" @click="emits('close')" />
   </div>
 </div>
</template>

<script setup>
import {ref, onMounted, watch} from 'vue';
import {showMessageError} from "@/utils/dialog";
import {Close} from "@element-plus/icons-vue";
import {formatTime} from "@/utils/libs";

const audio = ref(null);
const isPlaying = ref(false);
const songIndex = ref(0);
const currentTime = ref(0);
const duration = ref(100);
const progressPercent = ref(0);
const progressBarRef = ref(null)
const title = ref("")
const tags = ref("")
const cover = ref("")

const props = defineProps({
  songs: {
    type: Array,
    required: true,
    default: () => []
  },
});
// eslint-disable-next-line no-undef
const emits = defineEmits(['close']);

watch(() => props.songs, (newVal) => {
  console.log(newVal)
  loadSong(newVal[songIndex.value]);
});


const loadSong = (song) => {
  if (!song) {
    showMessageError("歌曲加载失败")
    return
  }
  title.value = song.title
  tags.value = song.tags
  cover.value = song.thumb_img_url
  audio.value.src = song.audio_url;
  audio.value.load();
  audio.value.onloadedmetadata = () => {
    duration.value = audio.value.duration;
  };
};

const togglePlay = () => {
  if (isPlaying.value) {
    audio.value.pause();
  } else {
    audio.value.play();
  }
  isPlaying.value = !isPlaying.value;
};

const play = () => {
  audio.value.play();
}

const prevSong = () => {
  songIndex.value = (songIndex.value - 1 + props.songs.length) % props.songs.length;
  loadSong(props.songs[songIndex.value]);
  audio.value.play();
  isPlaying.value = true;
};

const nextSong = () => {
  songIndex.value = (songIndex.value + 1) % props.songs.length;
  loadSong(props.songs[songIndex.value]);
  audio.value.play();
  isPlaying.value = true;
};

const updateProgress = () => {
  try {
    currentTime.value = audio.value.currentTime;
    progressPercent.value = (currentTime.value / duration.value) * 100;
  } catch (e) {
    console.error(e.message)
  }
};

const setProgress = (event) => {
  const totalWidth = progressBarRef.value.offsetWidth;
  const clickX = event.offsetX;
  const audioDuration = audio.value.duration;
  audio.value.currentTime = (clickX / totalWidth) * audioDuration;
};

// eslint-disable-next-line no-undef
defineExpose({
  play
});

onMounted(() => {
  loadSong(props.songs[songIndex.value]);
});
</script>

<style lang="stylus" scoped>

.player {
  display flex
  justify-content center
  width 100%

  .container {
    display flex
    background-color: #363030;
    border-radius: 10px;
    border 1px solid #544F4F;
    padding: 5px;
    width: 80%
    text-align: center;
    position relative
    overflow hidden


    .cover {
      .el-image {
        border-radius: 50%;
        width 50px
      }
    }

    .info {
      padding 0 10px
      min-width  300px
      display flex
      justify-content center
      align-items flex-start
      flex-flow column
      line-height 1.5

      .title {
        font-weight 700
        font-size 16px
      }

      .style {
        font-size 14px
        display flex
        color #e1e1e1
        .tags {
          font-weight 600
          white-space: nowrap; /* 防止文本换行 */
          overflow: hidden;    /* 隐藏溢出的文本 */
          text-overflow: ellipsis; /* 使用省略号表示溢出的文本 */
          max-width 200px
        }
        .text-lightGray {
          color: rgb(114 110 108);
          padding 0 3px
        }
        .time {
          font-family 'Input Sans'
          font-weight 700
          .split {
            font-size 12px
            position relative
            top -2px
            margin 0 1px 0 3px
          }
        }
      }
    }

    .controls-container {
      width 100%
      display flex
      flex-flow column
      justify-content center

      .controls {
        display: flex;
        justify-content: space-around;
        margin-bottom 10px
        .control-btn {
          background: none;
          border: none;
          color: #fff;
          cursor: pointer;
          background-color #363030
          border-radius 5px
          padding 6px

          .iconfont {
            font-size 20px
          }
          &:hover {
            background-color #5F5958
          }
        }

      }
    }

    .progress-bar {
      position absolute
      width 100%
      left 0
      bottom 0
      height: 8px;
      background-color: #555;
      cursor: pointer;

      .progress {
        height: 100%;
        background-color: #f50;
        border-radius: 5px;
        width: 0;
      }

    }

    .close {
      position absolute
      right 10px
      top 15px
    }
  }
}

</style>
