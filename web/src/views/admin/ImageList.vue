<template>
  <div class="container image-page">
    <el-tabs v-model="activeName" @tab-change="handleChange">
      <el-tab-pane label="Midjourney" name="mj" v-loading="data.mj.loading">
        <div class="handle-box">
          <el-input v-model="data.mj.query.username" placeholder="用户名" class="handle-input mr10"
                    @keyup="search($event,'mj')" clearable />
          <el-input v-model="data.mj.query.prompt" placeholder="提示词" class="handle-input mr10"
                    @keyup="search($event,'mj')" clearable />
          <el-date-picker
              v-model="data.mj.query.created_at"
              type="daterange"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              style="margin-right: 10px;width: 200px; position: relative;top:3px;"
          />
          <el-button type="primary" :icon="Search" @click="fetchMjData">搜索</el-button>
        </div>

        <div v-if="data.mj.items.length > 0">
          <el-row>
            <el-table :data="data.mj.items" :row-key="row => row.id" table-layout="auto">
              <el-table-column prop="user_id" label="用户ID"/>
              <el-table-column label="任务类型">
                <template #default="scope">
                  <el-button :color="taskTypeTheme[scope.row.type].color" size="small" plain>{{taskTypeTheme[scope.row.type].text}}</el-button>
                </template>
              </el-table-column>
              <el-table-column prop="progress" label="任务进度">
                <template #default="scope">
                  <span v-if="scope.row.progress <= 100">{{scope.row.progress}}%</span>
                  <el-tag v-else type="danger">已失败</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="power" label="消耗算力"/>
              <el-table-column label="结果图片">
                <template #default="scope">
                  <el-button size="small" type="success" @click="showImage(scope.row.img_url)" v-if="scope.row.img_url !== ''" plain>预览图片</el-button>
                </template>
              </el-table-column>
              <el-table-column label="提示词">
                <template #default="scope">
                  <el-popover
                      placement="top-start"
                      title="绘画提示词"
                      :width="300"
                      trigger="hover"
                      :content="scope.row.prompt"
                  >
                    <template #reference>
                      <span>{{ substr(scope.row.prompt, 20) }}</span>
                    </template>
                  </el-popover>
                </template>
              </el-table-column>
              <el-table-column label="创建时间">
                <template #default="scope">
                  <span>{{ dateFormat(scope.row['created_at']) }}</span>
                </template>
              </el-table-column>
              <el-table-column label="失败原因">
                <template #default="scope">
                  <el-popover
                      placement="top-start"
                      title="失败原因"
                      :width="300"
                      trigger="hover"
                      :content="scope.row.err_msg"
                      v-if="scope.row.progress === 101"
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
                  <el-popconfirm title="确定要删除当前记录吗?" @confirm="remove(scope.row, 'mj')">
                    <template #reference>
                      <el-button size="small" type="danger">删除</el-button>
                    </template>
                  </el-popconfirm>
                </template>
              </el-table-column>
            </el-table>
          </el-row>

          <div class="pagination">
            <el-pagination v-if="data.mj.total > 0" background
                           layout="total,prev, pager, next"
                           :hide-on-single-page="true"
                           v-model:current-page="data.mj.page"
                           v-model:page-size="data.mj.pageSize"
                           @current-change="fetchMjData()"
                           :total="data.mj.total"/>

          </div>
        </div>
        <el-empty v-else />
      </el-tab-pane>
      <el-tab-pane label="Stable-Diffusion" name="sd" v-loading="data.sd.loading">
        <div class="handle-box">
          <el-input v-model="data.sd.query.username" placeholder="用户名" class="handle-input mr10"
                    @keyup="search($event, 'sd')" clearable />
          <el-input v-model="data.sd.query.prompt" placeholder="提示词" class="handle-input mr10"
                    @keyup="search($event, 'sd')" clearable />
          <el-date-picker
              v-model="data.sd.query.created_at"
              type="daterange"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              style="margin-right: 10px;width: 200px; position: relative;top:3px;"
          />
          <el-button type="primary" :icon="Search" @click="fetchSdData">搜索</el-button>
        </div>

        <div v-if="data.sd.items.length > 0">
          <el-row>
            <el-table :data="data.sd.items" :row-key="row => row.id" table-layout="auto">
              <el-table-column prop="user_id" label="用户ID"/>
              <el-table-column prop="progress" label="任务进度">
                <template #default="scope">
                  <span v-if="scope.row.progress <= 100">{{scope.row.progress}}%</span>
                  <el-tag v-else type="danger">已失败</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="power" label="消耗算力"/>
              <el-table-column label="结果图片">
                <template #default="scope">
                  <el-button size="small" type="success" @click="showImage(scope.row.img_url)" v-if="scope.row.img_url !== ''" plain>预览图片</el-button>
                </template>
              </el-table-column>
              <el-table-column label="提示词">
                <template #default="scope">
                  <el-popover
                      placement="top-start"
                      title="绘画提示词"
                      :width="300"
                      trigger="hover"
                      :content="scope.row.prompt"
                  >
                    <template #reference>
                      <span>{{ substr(scope.row.prompt, 20) }}</span>
                    </template>
                  </el-popover>
                </template>
              </el-table-column>
              <el-table-column label="创建时间">
                <template #default="scope">
                  <span>{{ dateFormat(scope.row['created_at']) }}</span>
                </template>
              </el-table-column>
              <el-table-column label="失败原因">
                <template #default="scope">
                  <el-popover
                      placement="top-start"
                      title="失败原因"
                      :width="300"
                      trigger="hover"
                      :content="scope.row.err_msg"
                      v-if="scope.row.progress === 101"
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
                  <el-popconfirm title="确定要删除当前记录吗?" @confirm="remove(scope.row, 'sd')">
                    <template #reference>
                      <el-button size="small" type="danger">删除</el-button>
                    </template>
                  </el-popconfirm>
                </template>
              </el-table-column>
            </el-table>
          </el-row>

          <div class="pagination">
            <el-pagination v-if="data.sd.total > 0" background
                           layout="total,prev, pager, next"
                           :hide-on-single-page="true"
                           v-model:current-page="data.sd.page"
                           v-model:page-size="data.sd.pageSize"
                           @current-change="fetchSdData()"
                           :total="data.sd.total"/>

          </div>
        </div>
        <el-empty v-else />
      </el-tab-pane>
      <el-tab-pane label="DALL-E" name="dall">
        <div class="handle-box">
          <el-input v-model="data.dall.query.username" placeholder="用户名" class="handle-input mr10"
                    @keyup="search($event,'dall')" clearable />
          <el-input v-model="data.dall.query.prompt" placeholder="提示词" class="handle-input mr10"
                    @keyup="search($event, 'dall')" clearable />
          <el-date-picker
              v-model="data.dall.query.created_at"
              type="daterange"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              style="margin-right: 10px;width: 200px; position: relative;top:3px;"
          />
          <el-button type="primary" :icon="Search" @click="fetchDallData">搜索</el-button>
        </div>

        <div v-if="data.dall.items.length > 0">
          <el-row>
            <el-table :data="data.dall.items" :row-key="row => row.id" table-layout="auto">
              <el-table-column prop="user_id" label="用户ID"/>
              <el-table-column prop="progress" label="任务进度">
                <template #default="scope">
                  <span v-if="scope.row.progress <= 100">{{scope.row.progress}}%</span>
                  <el-tag v-else type="danger">已失败</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="power" label="消耗算力"/>
              <el-table-column label="结果图片">
                <template #default="scope">
                  <el-button size="small" type="success" @click="showImage(scope.row.img_url)" v-if="scope.row.img_url !== ''" plain>预览图片</el-button>
                </template>
              </el-table-column>
              <el-table-column label="提示词">
                <template #default="scope">
                  <el-popover
                      placement="top-start"
                      title="绘画提示词"
                      :width="300"
                      trigger="hover"
                      :content="scope.row.prompt"
                  >
                    <template #reference>
                      <span>{{ substr(scope.row.prompt, 20) }}</span>
                    </template>
                  </el-popover>
                </template>
              </el-table-column>
              <el-table-column label="创建时间">
                <template #default="scope">
                  <span>{{ dateFormat(scope.row['created_at']) }}</span>
                </template>
              </el-table-column>
              <el-table-column label="失败原因">
                <template #default="scope">
                  <el-popover
                      placement="top-start"
                      title="失败原因"
                      :width="300"
                      trigger="hover"
                      :content="scope.row.err_msg"
                      v-if="scope.row.progress === 101"
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
                  <el-popconfirm title="确定要删除当前记录吗?" @confirm="remove(scope.row, 'dall')">
                    <template #reference>
                      <el-button size="small" type="danger">删除</el-button>
                    </template>
                  </el-popconfirm>
                </template>
              </el-table-column>
            </el-table>
          </el-row>

          <div class="pagination">
            <el-pagination v-if="data.dall.total > 0" background
                           layout="total,prev, pager, next"
                           :hide-on-single-page="true"
                           v-model:current-page="data.dall.page"
                           v-model:page-size="data.dall.pageSize"
                           @current-change="fetchDallData()"
                           :total="data.dall.total"/>

          </div>
        </div>
        <el-empty v-else />
      </el-tab-pane>
    </el-tabs>


    <el-dialog
        v-model="showImageDialog"
        title="图片预览"
    >
      <el-image
          :src="imgURL"
          :zoom-rate="1.2"
          :max-scale="7"
          :min-scale="0.2"
          :preview-src-list="[imgURL]"
          :initial-index="0"
          fit="cover"
      />
    </el-dialog>

  </div>
</template>

<script setup>
import {onMounted, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {dateFormat, substr} from "@/utils/libs";
import {Search} from "@element-plus/icons-vue";

// 变量定义
const data = ref({
  "mj": {
    items: [],
    query: {prompt: "", username: "", created_at: [], page: 1, page_size: 15},
    total: 0,
    page: 1,
    pageSize: 15,
    loading: true
  },
  "sd": {
    items: [],
    query: {prompt: "", username: "", created_at: [], page: 1, page_size: 15},
    total: 0,
    page: 1,
    pageSize: 15,
    loading: true
  },
  "dall": {
    items: [],
    query: {prompt: "", username: "", created_at: [], page: 1, page_size: 15},
    total: 0,
    page: 1,
    pageSize: 15,
    loading: true
  }
})
const activeName = ref("mj")
const taskTypeTheme = {
  image: {text: "绘图", color: "#2185d0"},
  upscale: {text: "放大", color: "#f2711c" },
  variation: {text: "变换", color: "#00b5ad"},
  blend: {text: "融图", color: "#21ba45"},
  swapFace: {text: "换脸", color: "#a333c8"}
}

onMounted(() => {
  fetchMjData()
})

const handleChange = (tab) => {
  switch (tab) {
    case "mj":
      fetchMjData()
      break
    case "sd":
      fetchSdData()
      break
    case "dall":
      fetchDallData()
      break
  }
}

// 搜索对话
const search = (evt,tab) => {
  if (evt.keyCode === 13) {
    handleChange(tab)
  }
}

// 获取数据
const fetchMjData = () => {
  const d = data.value.mj
  d.query.page = d.page
  d.query.page_size = d.pageSize
  httpPost('/api/admin/image/list/mj', d.query).then((res) => {
    if (res.data) {
      d.items = res.data.items
      d.total = res.data.total
      d.page = res.data.page
      d.pageSize = res.data.page_size
    }
    d.loading = false
  }).catch(e => {
    ElMessage.error("获取数据失败：" + e.message);
  })
}

const fetchSdData = () => {
  const d = data.value.sd
  d.query.page = d.page
  d.query.page_size = d.pageSize
  httpPost('/api/admin/image/list/sd', d.query).then((res) => {
    if (res.data) {
      d.items = res.data.items
      d.total = res.data.total
      d.page = res.data.page
      d.pageSize = res.data.page_size
    }
    d.loading = false
  }).catch(e => {
    ElMessage.error("获取数据失败：" + e.message);
  })
}

const fetchDallData = () => {
  const d = data.value.dall
  d.query.page = d.page
  d.query.page_size = d.pageSize
  httpPost('/api/admin/image/list/dall', d.query).then((res) => {
    if (res.data) {
      d.items = res.data.items
      d.total = res.data.total
      d.page = res.data.page
      d.pageSize = res.data.page_size
    }
    d.loading = false
  }).catch(e => {
    ElMessage.error("获取数据失败：" + e.message);
  })
}

const remove = function (row,tab) {
  httpGet(`/api/admin/image/remove?id=${row.id}&tab=${tab}`).then(() => {
    ElMessage.success("删除成功！")
    handleChange(tab)
  }).catch((e) => {
    ElMessage.error("删除失败：" + e.message)
  })
}

const showImageDialog = ref(false)
const imgURL = ref('')
const showImage = (url) => {
  showImageDialog.value = true
  imgURL.value = url
}
</script>

<style lang="stylus" scoped>
.image-page {
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
    padding 20px 0
    display flex
    justify-content right
  }
}
</style>