<template>
  <div class="page-suno" :style="{ height: winHeight + 'px' }">
    <div class="left-bar">
      <div class="bar-top">
        <el-tooltip effect="light" content="定义模式" placement="top">
          <black-switch  v-model:value="custom" size="large"  />
        </el-tooltip>
        <black-select v-model:value="data.model" :options="models" placeholder="请选择模型" style="width: 100px" />
      </div>

      <div class="params">
        <div class="pure-music">
          <span class="switch"><black-switch  v-model:value="instrumental" size="default"  /></span>
          <span class="text">纯音乐</span>
        </div>
        <div v-if="custom">
          <div class="item-group" v-if="!instrumental">
            <div class="label">
              <span class="text">歌词</span>
              <el-popover placement="right"
                          :width="200"
                          trigger="hover" content="自己写歌词或寻求 AI 的帮助。使用两节歌词（8 行）可获得最佳效果。">
                <template #reference>
                  <el-icon>
                    <InfoFilled/>
                  </el-icon>
                </template>
              </el-popover>
            </div>
            <div class="item">
              <black-input v-model:value="data.lyrics" type="textarea" :rows="10" placeholder="请在这里输入你自己写的歌词..."/>
            </div>
          </div>

          <div class="item-group">
            <div class="label">
              <span class="text">音乐风格</span>
              <el-popover placement="right"
                          :width="200"
                          trigger="hover" content="描述您想要的音乐风格（例如“原声流行音乐”）。Sunos 模特无法识别艺术家的名字，但能够理解音乐流派和氛围。">
                <template #reference>
                  <el-icon>
                    <InfoFilled/>
                  </el-icon>
                </template>
              </el-popover>
            </div>
            <div class="item">
              <black-input v-model:value="data.tags" type="textarea" :rows="3" placeholder="请输入音乐风格，多个风格之间用英文逗号隔开..."/>
            </div>
          </div>

          <div class="item-group">
            <div class="label">
              <span class="text">歌曲名称</span>
              <el-popover placement="right"
                          :width="200"
                          trigger="hover" content="给你的歌曲起一个标题，以便于分享、发现和组织。">
                <template #reference>
                  <el-icon>
                    <InfoFilled/>
                  </el-icon>
                </template>
              </el-popover>
            </div>
            <div class="item">
              <black-input v-model:value="data.title" type="textarea" :rows="2" placeholder="请输入歌曲名称..."/>
            </div>
          </div>
        </div>

        <div v-else>
          <div class="label">
            <span class="text">歌曲描述</span>
            <el-popover placement="right"
                        :width="200"
                        trigger="hover" content="描述您想要的音乐风格和主题（例如：关于假期的流行音乐）。请使用流派和氛围，而不是特定的艺术家和歌曲风格，AI无法识别。">
              <template #reference>
                <el-icon>
                  <InfoFilled/>
                </el-icon>
              </template>
            </el-popover>
          </div>
          <div class="item">
            <black-input v-model:value="data.lyrics" type="textarea" :rows="10" placeholder="例如：一首关于鸟人的摇滚歌曲..."/>
          </div>
        </div>


        <div class="item">
          <button class="create-btn" @click="create">
            <img src="/images/create-new.svg" alt=""/>
            <span>生成音乐</span>
          </button>
        </div>
      </div>
    </div>
    <div class="right-box" v-loading="loading" element-loading-background="rgba(100,100,100,0.3)">
      <div class="list-box" v-if="!noData">
        <div class="item" v-for="item in list">
          <div class="left">
            <div class="container">
              <el-image :src="item.thumb_img_url" fit="cover" />
              <div class="duration">{{duration(item.duration)}}</div>
              <button class="play" @click="play(item)">
                <img src="/images/play.svg" alt=""/>
              </button>
            </div>
          </div>
          <div class="center">
            <div class="title">
              <a href="/song/xxxxx">{{item.title}}</a>
              <span class="model">{{item.model}}</span>
            </div>
            <div class="tags">{{item.tags}}</div>
          </div>
          <div class="right">
            <div class="tools">
              <el-tooltip effect="light" content="以当前歌曲为素材继续创作" placement="top">
                <button class="btn">续写</button>
              </el-tooltip>

              <button class="btn btn-publish">
                <span class="text">发布</span>
                <black-switch v-model:value="item.publish" size="small" />
              </button>

              <el-tooltip effect="light" content="下载歌曲" placement="top">
                <a :href="item.audio_url" :download="item.title+'.mp3'">
                  <button class="btn btn-icon">
                    <i class="iconfont icon-download"></i>
                  </button>
                </a>
              </el-tooltip>

              <el-tooltip effect="light" content="复制歌曲链接" placement="top">
                <button class="btn btn-icon">
                  <i class="iconfont icon-share1"></i>
                </button>
              </el-tooltip>

              <el-tooltip effect="light" content="编辑" placement="top">
                <button class="btn btn-icon">
                  <i class="iconfont icon-edit"></i>
                </button>
              </el-tooltip>

              <el-tooltip effect="light" content="删除" placement="top">
                <button class="btn btn-icon">
                  <i class="iconfont icon-remove"></i>
                </button>
              </el-tooltip>
            </div>
          </div>
        </div>
      </div>
      <el-empty :image-size="100" description="没有任何作品，赶紧去创作吧！" v-else/>

      <div class="music-player" v-if="showPlayer">
        <music-player :songs="playList" ref="playerRef" @close="showPlayer = false" />
      </div>
    </div>
  </div>
</template>

<script setup>
import {nextTick, ref} from "vue"
import {InfoFilled} from "@element-plus/icons-vue";
import BlackSelect from "@/components/ui/BlackSelect.vue";
import BlackSwitch from "@/components/ui/BlackSwitch.vue";
import BlackInput from "@/components/ui/BlackInput.vue";
import MusicPlayer from "@/components/MusicPlayer.vue";
import {compact} from "lodash";

const winHeight = ref(window.innerHeight - 50)
const custom = ref(false)
const instrumental = ref(false)
const models = ref([
  {label: "v3.0", value: "chirp-v3-0"},
  {label: "v3.5", value:"chirp-v3-5"}
])
const data = ref({
  model: "chirp-v3-0",
  tags: "",
  lyrics: "",
  prompt: "",
  title: ""
})
const loading = ref(false)
const noData = ref(false)
const playList = ref([])
const playerRef = ref(null)
const showPlayer = ref(false)
const list = ref([
    {
      id: 1,
      title: "鸟人传说 (Birdman Legend)",
      model: "v3",
      tags: "uplifting pop",
      thumb_img_url: "https://cdn2.suno.ai/image_047796ce-7bf3-4051-a59c-66ce60448ff2.jpeg?width=100",
      audio_url: "/files/suno.mp3",
      publish: true,
      duration: 134,
    },
  {
    id: 1,
    title: "我是一个鸟人",
    model: "v3",
    tags: "摇滚",
    publish: false,
    thumb_img_url: "https://cdn2.suno.ai/image_e5d25fd0-06a5-4cd7-910c-4b62872cd503.jpeg?width=100",
    audio_url: "/files/test.mp3",
    duration: 194,
  },
  {
    id: 1,
    title: "鸟人传说 (Birdman Legend)",
    model: "v3",
    tags: "uplifting pop",
    publish: true,
    thumb_img_url: "https://cdn2.suno.ai/image_047796ce-7bf3-4051-a59c-66ce60448ff2.jpeg?width=100",
    audio_url: "/files/suno.mp3",
    duration: 138,
  },
  {
    id: 1,
    title: "我是一个鸟人",
    model: "v3",
    tags: "摇滚",
    publish: false,
    thumb_img_url: "https://cdn2.suno.ai/image_e5d25fd0-06a5-4cd7-910c-4b62872cd503.jpeg?width=100",
    audio_url: "/files/suno.mp3",
    duration: 144,
  },
])

const create = () => {
  console.log(data.value)
}

const play = (item) => {
  playList.value = [item]
  showPlayer.value = true
  nextTick(()=> playerRef.value.play())
}
// 格式化音频时长
const duration = (secs) => {
  const minutes =Math.floor(secs/60)
  const seconds = secs%60
  return `${minutes}:${seconds}`
}
</script>

<style lang="stylus" scoped>
@import "@/assets/css/suno.styl"
</style>
