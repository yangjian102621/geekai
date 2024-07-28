<template>
  <div class="page-song" :style="{ height: winHeight + 'px' }">
    <div class="inner">
      <h2 class="title">{{song.title}}</h2>
      <div class="row tags" v-if="song.tags">
        <span>{{song.tags}}</span>
      </div>

      <div class="row author">
        <span>
          <el-avatar :size="32" :src="song.user?.avatar" />
        </span>
        <span class="nickname">{{song.user?.nickname}}</span>
        <button class="btn btn-icon" @click="play">
          <i class="iconfont icon-play"></i> {{song.play_times}}
        </button>

        <el-tooltip effect="light" content="复制歌曲链接" placement="top">
          <button class="btn btn-icon copy-link" :data-clipboard-text="getShareURL(song)" >
            <i class="iconfont icon-share1"></i>
          </button>
        </el-tooltip>
      </div>

      <div class="row date">
        <span>{{dateFormat(song.created_at)}}</span>
        <span class="version">{{song.raw_data?.major_model_version}}</span>
      </div>

      <div class="row">
        <textarea class="prompt" maxlength="2000" rows="18" readonly>{{song.prompt}}</textarea>
      </div>
    </div>

    <div class="music-player" v-if="playList.length > 0">
      <music-player :songs="playList" ref="playerRef" @play="song.play_times += 1"/>
    </div>
  </div>
</template>

<script setup>
import {onMounted, onUnmounted, ref} from "vue"
import {useRouter} from "vue-router";
import {httpGet} from "@/utils/http";
import {showMessageError} from "@/utils/dialog";
import {dateFormat} from "@/utils/libs";
import Clipboard from "clipboard";
import {ElMessage} from "element-plus";
import MusicPlayer from "@/components/MusicPlayer.vue";

const router = useRouter()
const id = router.currentRoute.value.params.id
const song = ref({title:""})
const playList = ref([])
const playerRef = ref(null)

httpGet("/api/suno/detail",{song_id:id}).then(res => {
  song.value = res.data
  playList.value = [song.value]
  document.title = song.value?.title+ " | By "+song.value?.user.nickname+" | Suno音乐"
}).catch(e => {
  showMessageError("获取歌曲详情失败："+e.message)
})

const clipboard = ref(null)
onMounted(() => {
  clipboard.value = new Clipboard('.copy-link');
  clipboard.value.on('success', () => {
    ElMessage.success("复制歌曲链接成功！");
  })

  clipboard.value.on('error', () => {
    ElMessage.error('复制失败！');
  })
})

onUnmounted(() => {
  clipboard.value.destroy()
})

// 播放歌曲
const play = () => {
  playerRef.value.play()
}


const winHeight = ref(window.innerHeight-50)
const getShareURL = (item) => {
  return `${location.protocol}//${location.host}/song/${item.id}`
}
</script>

<style lang="stylus" scoped>
@import "@/assets/css/song.styl"
</style>
