<template>
  <div class="container media-page">
    <el-tabs v-model="activeName" @tab-change="handleChange">
      <el-tab-pane v-for="media in mediaTypes" :key="media.name" :label="media.label" :name="media.name" v-loading="data[media.name].loading">
        <div class="handle-box">
          <el-input v-model="data[media.name].query.username" placeholder="用户名" class="handle-input mr10" @keyup="search($event, media.name)" clearable />
          <el-input v-model="data[media.name].query.prompt" placeholder="提示词" class="handle-input mr10" @keyup="search($event, media.name)" clearable />
          <el-date-picker
            v-model="data[media.name].query.created_at"
            type="daterange"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="margin-right: 10px; width: 200px; position: relative; top: 3px"
          />
          <el-button type="primary" :icon="Search" @click="media.fetchData">搜索</el-button>
        </div>

        <div v-if="data[media.name].items.length > 0">
          <el-row>
            <el-table :data="data[media.name].items" :row-key="(row) => row.id" table-layout="auto">
              <el-table-column prop="user_id" label="用户ID" />
              <el-table-column label="预览">
                <template #default="scope" v-if="media.previewComponent === 'MusicPreview'">
                  <div class="container" v-if="scope.row.cover_url">
                    <el-image :src="scope.row.cover_url" fit="cover" />
                    <div class="duration">
                      {{ formatTime(scope.row.duration) }}
                    </div>
                    <button class="play" @click="playMusic(scope.row)">
                      <img src="/images/play.svg" alt="" />
                    </button>
                  </div>
                  <el-image v-else src="/images/failed.jpg" style="height: 90px" fit="cover" />
                </template>
                <template #default="scope" v-if="media.previewComponent === 'VideoPreview'">
                  <div class="container">
                    <div v-if="scope.row.progress === 100">
                      <video class="video" :src="replaceImg(scope.row.video_url)" preload="auto" loop="loop" muted="muted">您的浏览器不支持视频播放</video>
                      <button class="play" @click="playVideo(scope.row)">
                        <img src="/images/play.svg" alt="" />
                      </button>
                    </div>
                    <el-image :src="scope.row.cover_url" fit="cover" v-else-if="scope.row.progress > 100" />
                    <generating message="正在生成视频" v-else />
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="progress" label="任务进度">
                <template #default="scope">
                  <span v-if="scope.row.progress <= 100">{{ scope.row.progress }}%</span>
                  <el-tag v-else type="danger">已失败</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="power" label="消耗算力" />
              <template v-if="media.previewComponent === 'MusicPreview'">
                <el-table-column prop="tags" label="风格" />
                <el-table-column prop="play_times" label="播放次数" />
                <el-table-column label="歌词">
                  <template #default="scope">
                    <el-button size="small" type="primary" plain @click="showLyric(scope.row)">查看歌词</el-button>
                  </template>
                </el-table-column>
              </template>
              <el-table-column label="提示词" v-if="media.previewComponent === 'VideoPreview'">
                <template #default="scope">
                  <el-popover placement="top-start" title="提示词" :width="300" trigger="hover" :content="scope.row.prompt">
                    <template #reference>
                      <span>{{ substr(scope.row.prompt, 20) }}</span>
                    </template>
                  </el-popover>
                </template>
              </el-table-column>
              <el-table-column label="创建时间">
                <template #default="scope">
                  <span>{{ dateFormat(scope.row["created_at"]) }}</span>
                </template>
              </el-table-column>
              <el-table-column label="失败原因">
                <template #default="scope">
                  <el-popover
                    v-if="scope.row.progress === 101"
                    placement="top-start"
                    title="失败原因"
                    :width="300"
                    trigger="hover"
                    :content="scope.row.err_msg"
                  >
                    <template #reference>
                      <el-text type="danger">{{ substr(scope.row.err_msg, 20) }}</el-text>
                    </template>
                  </el-popover>
                  <span v-else>无</span>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="180">
                <template #default="scope">
                  <el-popconfirm title="确定要删除当前记录吗?" @confirm="remove(scope.row, media.name)">
                    <template #reference>
                      <el-button size="small" type="danger">删除</el-button>
                    </template>
                  </el-popconfirm>
                </template>
              </el-table-column>
            </el-table>
          </el-row>

          <div class="pagination">
            <el-pagination
              v-if="data[media.name].total > 0"
              background
              layout="total,prev, pager, next"
              :hide-on-single-page="true"
              v-model:current-page="data[media.name].page"
              v-model:page-size="data[media.name].pageSize"
              @current-change="media.fetchData"
              :total="data[media.name].total"
            />
          </div>
        </div>
        <el-empty v-else />
      </el-tab-pane>
    </el-tabs>
    <el-dialog v-model="showVideoDialog" title="视频预览">
      <video style="width: 100%; max-height: 90vh" :src="currentVideoUrl" preload="auto" :autoplay="true" loop="loop" muted="muted">
        您的浏览器不支持视频播放
      </video>
    </el-dialog>

    <div class="music-player" v-if="showPlayer">
      <music-player :songs="playList" ref="playerRef" :show-close="true" @close="showPlayer = false" />
    </div>

    <el-dialog v-model="showLyricDialog" title="歌词">
      <div class="chat-line" v-html="lyrics"></div>
    </el-dialog>
  </div>
</template>

<script setup>
import { nextTick, onMounted, ref } from "vue";
import { httpGet, httpPost } from "@/utils/http";
import { ElMessage } from "element-plus";
import { dateFormat, formatTime, replaceImg, substr } from "@/utils/libs";
import { Search } from "@element-plus/icons-vue";
import MusicPlayer from "@/components/MusicPlayer.vue";
import Generating from "@/components/ui/Generating.vue";

const data = ref({
  suno: {
    items: [],
    query: { prompt: "", username: "", created_at: [], page: 1, page_size: 15 },
    total: 0,
    page: 1,
    pageSize: 10,
    loading: true,
  },
  luma: {
    items: [],
    query: { prompt: "", username: "", created_at: [], page: 1, page_size: 15 },
    total: 0,
    page: 1,
    pageSize: 10,
    loading: true,
  },
  keling: {
    items: [],
    query: { prompt: "", username: "", created_at: [], page: 1, page_size: 15 },
    total: 0,
    page: 1,
    pageSize: 10,
    loading: true,
  },
});

const mediaTypes = [
  {
    name: "suno",
    label: "Suno音乐",
    fetchData: () => fetchSunoData(),
    previewComponent: "MusicPreview",
  },
  {
    name: "luma",
    label: "Luma视频",
    fetchData: () => fetchLumaData(),
    previewComponent: "VideoPreview",
  },
  {
    name: "keling",
    label: "可灵视频",
    fetchData: () => fetchKelingData(),
    previewComponent: "VideoPreview",
  },
];
const activeName = ref("suno");
const playList = ref([]);
const playerRef = ref(null);
const showPlayer = ref(false);
const showLyricDialog = ref(false);
const lyrics = ref("");
const showVideoDialog = ref(false);
const currentVideoUrl = ref("");

onMounted(() => {
  fetchSunoData();
});

const handleChange = (tab) => {
  switch (tab) {
    case "suno":
      fetchSunoData();
      break;
    case "luma":
      fetchLumaData();
      break;
    case "keling":
      fetchKelingData();
      break;
  }
};

// 搜索对话
const search = (evt, tab) => {
  if (evt.keyCode === 13) {
    handleChange(tab);
  }
};

// 获取数据
const fetchSunoData = () => {
  const d = data.value.suno;
  d.query.page = d.page;
  d.query.page_size = d.pageSize;
  httpPost("/api/admin/media/suno", d.query)
    .then((res) => {
      if (res.data) {
        d.items = res.data.items;
        d.total = res.data.total;
        d.page = res.data.page;
        d.pageSize = res.data.page_size;
      }
      d.loading = false;
    })
    .catch((e) => {
      ElMessage.error("获取数据失败：" + e.message);
    });
};

const fetchLumaData = () => {
  const d = data.value.luma;
  d.query.page = d.page;
  d.query.page_size = d.pageSize;
  d.query.type = "luma";
  httpPost("/api/admin/media/videos", d.query)
    .then((res) => {
      if (res.data) {
        d.items = res.data.items;
        d.total = res.data.total;
        d.page = res.data.page;
        d.pageSize = res.data.page_size;
      }
      d.loading = false;
    })
    .catch((e) => {
      ElMessage.error("获取数据失败：" + e.message);
    });
};
const fetchKelingData = () => {
  const d = data.value.keling;
  d.query.page = d.page;
  d.query.page_size = d.pageSize;
  d.query.type = "keling";
  httpPost("/api/admin/media/videos", d.query)
    .then((res) => {
      if (res.data) {
        d.items = res.data.items;
        d.total = res.data.total;
        d.page = res.data.page;
        d.pageSize = res.data.page_size;
      }
      d.loading = false;
    })
    .catch((e) => {
      ElMessage.error("获取数据失败：" + e.message);
    });
};

const remove = function (row, tab) {
  httpGet(`/api/admin/media/remove?id=${row.id}&tab=${tab}`)
    .then(() => {
      ElMessage.success("删除成功！");
      handleChange(tab);
    })
    .catch((e) => {
      ElMessage.error("删除失败：" + e.message);
    })
    .finally(() => {
      nextTick(() => {
        // data.value[tab].page = 1;
        // data.value[tab].pageSize = 10;
        const mediaType = mediaTypes.find((type) => type.name === tab);
        if (mediaType && mediaType.fetchData) {
          mediaType.fetchData();
        }
      });
    });
};

const playMusic = (item) => {
  playList.value = [item];
  showPlayer.value = true;
  nextTick(() => playerRef.value.play());
};

const playVideo = (item) => {
  currentVideoUrl.value = replaceImg(item.video_url);
  showVideoDialog.value = true;
};

const md = require("markdown-it")({
  breaks: true,
  html: true,
  linkify: true,
});

const showLyric = (item) => {
  showLyricDialog.value = true;
  lyrics.value = md.render(item.prompt);
};
</script>

<style lang="stylus" scoped>
.media-page {
  .handle-box {
    margin-bottom 20px
    .handle-input {
      max-width 150px;
      margin-right 10px;
    }
  }

  .opt-box {
    padding-bottom: 10px;
    display flex;
    justify-content flex-end

    .el-icon {
      margin-right: 5px;
    }
  }

  .el-select {
    width: 100%
  }

  .pagination {
    margin-top 20px
    display flex
    justify-content center
  }

  .container {
    width 160px
    position relative

    .video{
      width 160px
      border-radius 5px
    }

    .el-image {
      width 160px
      height 90px
      border-radius 5px
    }

    .duration {
      position absolute
      bottom 6px
      right 0
      background-color rgba(255, 255, 255,.7)
      padding 0 3px
      font-family 'Input Sans'
      font-size 14px
      font-weight 700
      border-radius .125rem
    }

    .play {
      position absolute
      width: 100%
      height 100%
      top: 0;
      left: 50%;
      border none
      border-radius 5px
      background rgba(100, 100, 100, 0.3)
      cursor pointer
      color #ffffff
      opacity 0
      transform: translate(-50%, 0px);
      transition opacity 0.3s ease 0s
    }

    &:hover {
      .play {
        opacity 1
        //display block
      }
    }
  }


  .music-player {
    position absolute
    bottom 20px
    z-index 99999
    width 100%
  }
}
</style>
