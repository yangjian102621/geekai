<template>
  <div class="system-config form" v-loading="loading">
    <div class="container">
      <el-form
        :model="system"
        label-width="150px"
        label-position="right"
        ref="systemFormRef"
        :rules="rules"
      >
        <el-tabs type="border-card">
          <el-tab-pane>
            <template #label>
              <i class="iconfont icon-token mr-1"></i>
              <span>秘钥配置</span>
            </template>
            <el-form-item label="网站标题" prop="title">
              <el-input v-model="system['title']" />
            </el-form-item>
          </el-tab-pane>

          <el-tab-pane>
            <template #label>
              <i class="iconfont icon-logout mr-1"></i>
              <span>算力配置</span>
            </template>
            <el-form-item label="注册赠送算力" prop="init_power">
              <el-input v-model.number="system['init_power']" placeholder="新用户注册赠送算力" />
            </el-form-item>

            <el-form-item>
              <template #label>
                <div class="label-title">
                  提示词算力
                  <el-tooltip
                    effect="dark"
                    content="生成AI绘图提示词，歌词，视频描述消耗的算力"
                    raw-content
                    placement="right"
                  >
                    <el-icon>
                      <InfoFilled />
                    </el-icon>
                  </el-tooltip>
                </div>
              </template>
              <el-input v-model.number="system['prompt_power']" placeholder="" />
            </el-form-item>
          </el-tab-pane>
        </el-tabs>

        <div style="padding: 10px">
          <el-form-item>
            <el-button type="primary" @click="save('system')">保存</el-button>
          </el-form-item>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { httpGet, httpPost } from '@/utils/http'
import { InfoFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import 'md-editor-v3/lib/style.css'
import { onMounted, ref } from 'vue'

const system = ref({ models: [] })
const loading = ref(true)

onMounted(() => {
  // 加载系统配置
  httpGet('/api/admin/config/get?key=system')
    .then((res) => {
      system.value = res.data
    })
    .catch((e) => {
      ElMessage.error('加载系统配置失败: ' + e.message)
    })
    .finally(() => {
      loading.value = false
    })
})

const save = function (key) {
  httpPost('/api/admin/config/update', {
    key: key,
    config: { content: notice.value, updated: true },
  })
    .then(() => {
      ElMessage.success('操作成功！')
    })
    .catch((e) => {
      ElMessage.error('操作失败：' + e.message)
    })
}
</script>

<style lang="stylus" scoped>
@import '../../../assets/css/admin/form.styl'
@import '../../../assets/css/main.styl'
.system-config {
  display flex
  justify-content center

  .sys-tabs {
    width 100%
    background-color var(--el-bg-color)
    padding 10px 20px 40px 20px
    //border: 1px solid var(--el-border-color);
  }
}
</style>
