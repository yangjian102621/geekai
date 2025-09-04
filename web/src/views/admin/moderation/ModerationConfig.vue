<template>
  <div class="settings container p-5">
    <el-tabs v-model="activeTab" type="border-card">
      <el-tab-pane name="gitee">
        <template #label>
          <div class="flex items-center">
            <i class="iconfont icon-gitee"></i>
            <span class="ml-2">模力方舟</span>
          </div>
        </template>

        <Alert type="info">
          如果你不知道怎么获取这些配置信息，请参考文档：
          <a href="https://ai.gitee.com/docs/organization/access-token" target="_blank"
            >模力方舟访问令牌配置</a
          >。
        </Alert>

        <el-form :model="configs.gitee" label-position="top">
          <el-form-item label="API密钥">
            <el-input v-model="configs.gitee.api_key" placeholder="请输入API密钥" />
          </el-form-item>
          <el-form-item label="模型">
            <el-select v-model="configs.gitee.model" placeholder="请选择模型">
              <el-option v-for="v in models" :value="v.value" :label="v.label" :key="v.value">
                {{ v.label }}
              </el-option>
            </el-select>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane name="baidu">
        <template #label>
          <div class="flex items-center">
            <i class="iconfont icon-baidu"></i>
            <span class="ml-2">百度</span>
          </div>
        </template>

        <Alert type="warning"> 百度文本审查服务暂未实现。 </Alert>

        <el-form :model="configs.baidu" label-position="top">
          <el-form-item label="AccessKey">
            <el-input v-model="configs.baidu.access_key" placeholder="请输入AccessKey" />
          </el-form-item>
          <el-form-item label="SecretKey">
            <el-input v-model="configs.baidu.secret_key" placeholder="请输入SecretKey" />
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane name="tencent">
        <template #label>
          <div class="flex items-center">
            <i class="iconfont icon-tencent"></i>
            <span class="ml-2">腾讯云</span>
          </div>
        </template>

        <Alert type="warning"> 腾讯云文本审查服务暂未实现。 </Alert>

        <el-form :model="configs.baidu" label-position="top">
          <el-form-item label="AccessKey">
            <el-input v-model="configs.baidu.access_key" placeholder="请输入AccessKey" />
          </el-form-item>
          <el-form-item label="SecretKey">
            <el-input v-model="configs.baidu.secret_key" placeholder="请输入SecretKey" />
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>

    <el-form :model="configs" label-position="top" class="py-5">
      <el-form-item label="启用模型引导提示词">
        <el-switch v-model="configs.enable_guide" />
      </el-form-item>

      <el-form-item v-if="configs.enable_guide">
        <template #label>
          <span class="mr-2">大模型引导提示词</span>
          <el-tooltip
            effect="dark"
            content="大模型引导提示词，用于引导大模型进行文本审查<br/>如果为空，则不使用大模型引导提示词"
            placement="right"
            raw-content
          >
            <i class="iconfont icon-info"></i>
          </el-tooltip>
        </template>
        <el-input
          v-model="configs.guide_prompt"
          type="textarea"
          :rows="3"
          placeholder="请输入大模型引导提示词"
        />
      </el-form-item>

      <el-form-item label="启用文本审查服务">
        <el-switch v-model="configs.enable" />
      </el-form-item>

      <el-form-item v-if="configs.enable">
        <template #label>
          <div class="flex items-center">
            <span class="mr-2">选择审查服务</span>
            <el-tooltip
              effect="dark"
              content="只有当文本审查启用时，选择审查服务才会生效"
              placement="right"
            >
              <i class="iconfont icon-info"></i>
            </el-tooltip>
          </div>
        </template>
        <el-radio-group v-model="configs.active" size="large">
          <el-radio value="gitee" border>模力方舟</el-radio>
          <el-radio value="baidu" border>百度</el-radio>
          <el-radio value="tencent" border>腾讯云</el-radio>
        </el-radio-group>
      </el-form-item>
    </el-form>

    <div class="flex justify-left">
      <el-button type="primary" @click="saveModerationConfig" :loading="loading"
        >提交保存</el-button
      >
    </div>

    <div class="mt-7">
      <el-card shadow="never" class="mb-6">
        <el-form :model="testForm" label-position="top">
          <el-form-item label="测试文本">
            <el-input
              v-model="testForm.text"
              type="textarea"
              :rows="4"
              placeholder="请输入要测试的文本内容"
              maxlength="1000"
              show-word-limit
            />
          </el-form-item>
          <el-form-item>
            <el-button
              type="success"
              @click="testModeration"
              :loading="testLoading"
              :disabled="!testForm.text.trim()"
            >
              提交测试
            </el-button>
          </el-form-item>
        </el-form>

        <!-- 测试结果显示 -->
        <div v-if="testResult" class="test-result">
          <div class="result-header mb-4">
            <div class="flex items-center py-2">
              <span class="text-base font-semibold mr-2">检测结果:</span>
              <el-tag :type="testResult.isAbnormal ? 'danger' : 'success'" size="large">
                <i
                  class="iconfont"
                  :class="testResult.isAbnormal ? 'icon-error' : 'icon-success'"
                ></i>
                <span class="text-sm ml-2">{{ testResult.isAbnormal ? '异常' : '正常' }}</span>
              </el-tag>
            </div>
            <p class="text-sm text-gray-500 mt-2">检测结果仅供参考</p>
          </div>

          <el-table :data="testResult.details" border class="result-table">
            <el-table-column prop="category" label="类别" width="120">
              <template #default="{ row }">
                <span class="font-medium">{{ row.category }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="描述" min-width="300">
              <template #default="{ row }">
                <span class="text-gray-700">{{ row.description }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="confidence" label="置信度" width="100" align="center">
              <template #default="{ row }">
                <span class="font-mono">{{ row.confidence }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="isCategory" label="是否为该类别" width="120" align="center">
              <template #default="{ row }">
                <el-tag :type="row.isCategory ? 'danger' : 'success'" size="small">
                  {{ row.isCategory ? '是' : '否' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import Alert from '@/components/ui/Alert.vue'
import { showMessageError } from '@/utils/dialog'
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { onMounted, ref, watch } from 'vue'

const loading = ref(false)
const activeTab = ref('gitee')
const configs = ref({
  enable: false,
  active: 'gitee',
  guide_prompt:
    '请拒绝输出任何有关色情，暴力相关内容，禁止输出跟中国政治相关的内容，比如政治敏感事件，国家领导人敏感信息等相关的内容。任何时刻都必须牢记这一原则。',
  gitee: {
    api_key: '',
    model: 'Security-semantic-filtering',
  },
  baidu: {
    access_key: '',
    secret_key: '',
  },
  tencent: {
    access_key: '',
    secret_key: '',
  },
})

// 测试相关数据
const testLoading = ref(false)
const testForm = ref({
  text: '',
})
const testResult = ref(null)

const models = ref([
  {
    label: '违规文本检测模型：限时免费',
    value: 'Security-semantic-filtering',
  },
  {
    label: '文本审核模型：0.0002元/条',
    value: 'moark-text-moderation',
  },
])

onMounted(async () => {
  try {
    const res = await httpGet('/api/admin/config/get?key=moderation')
    configs.value = Object.assign(configs.value, res.data)
  } catch (e) {
    // 使用默认值
    showMessageError('加载文本审查配置失败: ' + e.message)
  }
})

// 监听tab切换，清空测试结果
watch(activeTab, (newTab) => {
  if (newTab !== 'test') {
    testResult.value = null
    testForm.value.text = ''
  }
})

const saveModerationConfig = async () => {
  loading.value = true
  try {
    await httpPost('/api/admin/moderation/config', configs.value)
    ElMessage.success('保存成功')
  } catch (e) {
    ElMessage.error('保存失败：' + (e.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 测试文本审核服务
const testModeration = async () => {
  if (!testForm.value.text.trim()) {
    ElMessage.warning('请输入测试文本')
    return
  }

  // 检查是否启用了文本审查
  if (!configs.value.enable) {
    ElMessage.warning('请先启用文本审查服务')
    return
  }

  testLoading.value = true
  try {
    const res = await httpPost('/api/admin/moderation/test', {
      text: testForm.value.text.trim(),
      service: configs.value.active,
    })

    // 处理测试结果
    testResult.value = {
      isAbnormal: res.data.isAbnormal || false,
      details: res.data.details || [],
    }

    ElMessage.success('测试完成')

    // 清空输入框，提升用户体验
    testForm.value.text = ''
  } catch (e) {
    ElMessage.error('测试失败：' + (e.message || '未知错误'))
    // 清空之前的结果
    testResult.value = null
  } finally {
    testLoading.value = false
  }
}
</script>

<style lang="scss">
.settings {
  a {
    color: #409eff;
    &:hover {
      text-decoration: underline;
    }
  }
  .el-form-item__label {
    font-weight: 700;
  }
}

// 测试相关样式
.test-result {
  .result-header {
    .status-badge {
      display: inline-block;
      margin-left: 12px;

      .status-tag {
        font-size: 14px;
        padding: 8px 16px;

        .iconfont {
          margin-right: 6px;
        }
      }
    }
  }

  .result-table {
    .el-table__header {
      background-color: #f5f7fa;

      th {
        background-color: #f5f7fa;
        color: #606266;
        font-weight: 600;
      }
    }

    .el-table__row {
      &:hover {
        background-color: #f5f7fa;
      }
    }
  }
}
</style>
