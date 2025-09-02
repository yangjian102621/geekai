<template>
  <div class="admin-threed-setting">
    <!-- 配置表单 -->
    <div class="settings-container">
      <!-- 配置选项卡 -->
      <el-card class="setting-card">
        <el-tabs v-model="activeTab" type="border-card">
          <!-- 腾讯混元3D配置 -->
          <el-tab-pane name="tencent">
            <template #label>
              <div class="tab-label">
                <i class="iconfont icon-tencent mr-2"></i>
                <span>腾讯混元3D</span>
              </div>
            </template>
            <div class="tab-content">
              <!-- 秘钥配置 -->
              <div class="config-section">
                <h4>秘钥配置</h4>
                <el-form :model="configs.tencent" label-width="140px" label-position="top">
                  <el-form-item label="SecretId">
                    <el-input
                      v-model="configs.tencent.secret_id"
                      placeholder="请输入腾讯云SecretId"
                      show-password
                    />
                  </el-form-item>

                  <el-form-item label="SecretKey">
                    <el-input
                      v-model="configs.tencent.secret_key"
                      placeholder="请输入腾讯云SecretKey"
                      show-password
                    />
                  </el-form-item>

                  <el-form-item label="地域(目前仅支持广州)">
                    <el-input v-model="configs.tencent.region" placeholder="请输入地域" />
                  </el-form-item>

                  <el-form-item label="启用状态">
                    <el-switch v-model="configs.tencent.enabled" />
                  </el-form-item>
                </el-form>
              </div>

              <!-- 模型配置 -->
              <div class="config-section">
                <h4>模型配置</h4>
                <div class="model-config">
                  <div class="model-header">
                    <span>支持的3D模型格式和算力消耗</span>
                    <el-button type="primary" plain @click="addTencentModel">添加模型</el-button>
                  </div>

                  <el-table
                    :data="configs.tencent.models"
                    border
                    style="width: 100%"
                    :max-height="400"
                    size="small"
                  >
                    <el-table-column prop="name" label="模型名称" min-width="180">
                      <template #default="{ row }">
                        <el-input v-model="row.name" placeholder="模型名称" />
                      </template>
                    </el-table-column>
                    <el-table-column prop="desc" label="模型描述" min-width="180">
                      <template #default="{ row }">
                        <el-input
                          v-model="row.desc"
                          placeholder="模型描述"
                          type="textarea"
                          :rows="3"
                        />
                      </template>
                    </el-table-column>
                    <el-table-column prop="power" label="算力消耗" min-width="120">
                      <template #default="{ row }">
                        <el-input-number v-model="row.power" :min="1" :max="1000" />
                      </template>
                    </el-table-column>
                    <el-table-column prop="formats" label="输出格式" min-width="200">
                      <template #default="{ row }">
                        <el-select
                          v-model="row.formats"
                          multiple
                          placeholder="选择输出格式"
                          style="width: 100%"
                          collapse-tags
                          collapse-tags-tooltip
                        >
                          <el-option
                            v-for="item in formatOptions"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                          />
                        </el-select>
                      </template>
                    </el-table-column>
                    <el-table-column label="操作" min-width="100" fixed="right">
                      <template #default="{ $index }">
                        <el-button size="small" type="danger" @click="removeTencentModel($index)">
                          删除
                        </el-button>
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
              </div>
            </div>
          </el-tab-pane>

          <!-- Gitee模力方舟配置 -->
          <el-tab-pane name="gitee">
            <template #label>
              <div class="tab-label">
                <i class="iconfont icon-gitee mr-2"></i>
                <span>Gitee模力方舟</span>
              </div>
            </template>
            <div class="tab-content">
              <Alert type="info">
                如果你不知道怎么获取这些配置信息，请参考文档：
                <a href="https://ai.gitee.com/docs/organization/access-token" target="_blank"
                  >模力方舟访问令牌配置</a
                >。
              </Alert>
              <!-- 秘钥配置 -->
              <div class="config-section mt-5">
                <h4>秘钥配置</h4>
                <el-form :model="configs.gitee" label-width="140px" label-position="top">
                  <el-form-item label="API密钥">
                    <el-input
                      v-model="configs.gitee.api_key"
                      placeholder="请输入Gitee API密钥"
                      show-password
                    />
                  </el-form-item>

                  <el-form-item label="启用状态">
                    <el-switch v-model="configs.gitee.enabled" />
                  </el-form-item>
                </el-form>
              </div>

              <!-- 模型配置 -->
              <div class="config-section">
                <h4>模型配置</h4>
                <div class="model-config">
                  <div class="model-header">
                    <span>支持的3D模型格式和算力消耗</span>
                    <el-button type="primary" plain @click="addGiteeModel">添加模型</el-button>
                  </div>

                  <el-table
                    :data="configs.gitee.models"
                    border
                    style="width: 100%"
                    :max-height="400"
                    size="small"
                  >
                    <el-table-column prop="name" label="模型名称" min-width="180">
                      <template #default="{ row }">
                        <el-input v-model="row.name" placeholder="模型名称" />
                      </template>
                    </el-table-column>
                    <el-table-column prop="desc" label="模型描述" min-width="180">
                      <template #default="{ row }">
                        <el-input
                          v-model="row.desc"
                          placeholder="模型描述"
                          type="textarea"
                          :rows="3"
                        />
                      </template>
                    </el-table-column>
                    <el-table-column prop="power" label="算力消耗" min-width="120">
                      <template #default="{ row }">
                        <el-input-number v-model="row.power" :min="1" :max="1000" />
                      </template>
                    </el-table-column>
                    <el-table-column prop="formats" label="输出格式" min-width="200">
                      <template #default="{ row }">
                        <el-select
                          v-model="row.formats"
                          multiple
                          placeholder="选择输出格式"
                          style="width: 100%"
                          collapse-tags
                          collapse-tags-tooltip
                        >
                          <el-option
                            v-for="item in formatOptions"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                          />
                        </el-select>
                      </template>
                    </el-table-column>
                    <el-table-column label="操作" min-width="100" fixed="right">
                      <template #default="{ $index }">
                        <el-button size="small" type="danger" @click="removeGiteeModel($index)">
                          删除
                        </el-button>
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
              </div>
            </div>
          </el-tab-pane>

          <div class="flex justify-center mb-5">
            <el-button type="primary" @click="saveConfig" :loading="loading">保存配置</el-button>
          </div>
        </el-tabs>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import Alert from '@/components/ui/Alert.vue'
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, ref } from 'vue'

// 响应式数据
const activeTab = ref('tencent')
const loading = ref(false)
const configs = ref({
  tencent: { region: 'ap-guangzhou', enabled: true, models: [] },
  gitee: { models: [] },
})

const formatOptions = ref([
  { label: 'OBJ', value: 'OBJ' },
  { label: 'GLB', value: 'GLB' },
  { label: 'STL', value: 'STL' },
  { label: 'USDZ', value: 'USDZ' },
  { label: 'FBX', value: 'FBX' },
  { label: 'MP4', value: 'MP4' },
])

// 方法
const loadConfig = async () => {
  try {
    const res = await httpGet('/api/admin/config/get?key=ai3d')
    configs.value = res.data
    const models = await httpGet('/api/admin/ai3d/models')
    if (!configs.value.tencent.models || configs.value.tencent.models.length === 0) {
      configs.value.tencent.models = models.data.tencent
    }
    if (!configs.value.gitee.models || configs.value.gitee.models.length === 0) {
      configs.value.gitee.models = models.data.gitee
    }
  } catch (error) {
    ElMessage.error('加载配置失败：' + error.message)
  }
}

const saveConfig = async () => {
  loading.value = true
  try {
    const response = await httpPost('/api/admin/ai3d/config', configs.value)
    ElMessage.success('所有配置保存成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('保存失败：' + error.message)
    }
  } finally {
    loading.value = false
  }
}

// 模型操作 - 腾讯
const addTencentModel = () => {
  configs.value.tencent.models.push({
    name: '',
    desc: '',
    power: 1,
    formats: [],
  })
}

const removeTencentModel = async (index) => {
  try {
    await ElMessageBox.confirm('确定要删除该模型吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    configs.value.tencent.models.splice(index, 1)
    ElMessage.success('删除成功')
  } catch (e) {
    // 用户取消
  }
}

// 模型操作 - Gitee
const addGiteeModel = () => {
  configs.value.gitee.models.push({
    name: '',
    desc: '',
    power: 1,
    formats: [],
  })
}

const removeGiteeModel = async (index) => {
  try {
    await ElMessageBox.confirm('确定要删除该模型吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    configs.value.gitee.models.splice(index, 1)
    ElMessage.success('删除成功')
  } catch (e) {
    // 用户取消
  }
}

// 生命周期
onMounted(() => {
  loadConfig()
})
</script>

<style lang="scss">
.admin-threed-setting {
  padding: 20px;

  a {
    color: #409eff;
    &:hover {
      text-decoration: underline;
    }
  }
  .el-form-item__label {
    font-weight: 700;
  }

  .settings-container {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .setting-card {
    .el-card__body {
      padding: 0;
    }
  }

  .tab-content {
    padding: 20px;
  }

  .config-section {
    margin-bottom: 30px;

    h4 {
      margin: 0 0 16px 0;
      font-size: 16px;
      font-weight: 600;
      padding-bottom: 8px;
      border-bottom: 2px solid #409eff;
    }

    .section-actions {
      margin-top: 16px;
      text-align: right;
    }
  }

  .model-config {
    .model-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;

      .el-button {
        font-weight: 500;
      }
    }
  }

  .el-form-item {
    margin-bottom: 20px;
  }

  :deep(.el-tabs__header) {
    margin: 0;
  }

  :deep(.el-tabs__content) {
    padding: 0;
  }
}
</style>
