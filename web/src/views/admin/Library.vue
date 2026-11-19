<template>
  <div class="container library" v-loading="loading">
    <div class="library-heading">
      <div>
        <h2>公共资料库</h2>
        <p>上架后，ThinkBuddy 自动连接 WeKnora，所有用户都可以阅读和使用资料。</p>
      </div>
      <el-button @click="load">刷新</el-button>
    </div>
    <el-tabs v-model="tab">
      <el-tab-pane label="发布目录" name="resources">
        <el-empty v-if="!resources.length" description="先添加资料来源，再选择知识库发布">
          <el-button type="primary" @click="tab = 'sources'">管理资料来源</el-button>
        </el-empty>
        <el-table v-else :data="resources" row-key="id">
          <el-table-column prop="name" label="资料名称" min-width="160" />
          <el-table-column label="资料来源" min-width="150">
            <template #default="{ row }">{{ sourceName(row.source_id) }}</template>
          </el-table-column>
          <el-table-column prop="description" label="介绍" min-width="180" show-overflow-tooltip />
          <el-table-column prop="sort_num" label="排序" width="80" />
          <el-table-column label="发布状态" width="120">
            <template #default="{ row }">
              <el-switch :model-value="row.published" :loading="busy === 'publish-' + row.id" @change="setPublished(row, $event)" />
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180">
            <template #default="{ row }">
              <el-button size="small" @click="editResource(row)">编辑</el-button>
              <el-button size="small" type="danger" @click="remove('resources', row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
      <el-tab-pane label="资料来源" name="sources">
        <el-button type="primary" :icon="Plus" @click="editSource()">新增来源</el-button>
        <el-table :data="sources" row-key="id" class="source-table">
          <el-table-column prop="name" label="来源名称" min-width="150" />
          <el-table-column prop="base_url" label="服务器地址" min-width="200" show-overflow-tooltip />
          <el-table-column label="连接状态" width="100">
            <template #default="{ row }"><el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '已启用' : '已停用' }}</el-tag></template>
          </el-table-column>
          <el-table-column label="操作" min-width="270">
            <template #default="{ row }">
              <el-button size="small" type="primary" :loading="busy === 'test-' + row.id" @click="testSource(row)">测试并选择知识库</el-button>
              <el-button size="small" @click="editSource(row)">编辑</el-button>
              <el-button size="small" type="danger" @click="remove('sources', row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>
    <el-dialog v-model="sourceDialog" :title="sourceForm.id ? '编辑资料来源' : '新增资料来源'" width="580px" :close-on-click-modal="false">
      <el-form label-width="110px" @submit.prevent="saveSource">
        <el-form-item label="来源名称" required><el-input v-model="sourceForm.name" maxlength="120" /></el-form-item>
        <el-form-item label="服务器地址" required><el-input v-model="sourceForm.base_url" placeholder="https://knowledge.example.com" /></el-form-item>
        <el-form-item label="API Key" :required="!sourceForm.has_api_key">
          <el-input v-model="sourceForm.api_key" type="password" show-password autocomplete="new-password" :placeholder="sourceForm.has_api_key ? '已保存；留空保留原有密钥' : '填写该来源的读取凭据'" />
        </el-form-item>
        <el-form-item label="启用来源"><el-switch v-model="sourceForm.enabled" /></el-form-item>
        <p class="library-hint">此 API Key 同时用于后台读取和客户端直连，会随已发布知识库下发。请使用只允许读取赠送资料的 Key，服务器地址需能被用户设备访问。</p>
        <p class="library-hint">停用来源后停止继续发放，用户已添加的知识库保留。停止访问需在 WeKnora 中撤销 Key。</p>
      </el-form>
      <template #footer>
        <el-button @click="sourceDialog = false">取消</el-button>
        <el-button type="primary" :loading="busy === 'save-source'" @click="saveSource">保存来源</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="collectionDialog" :title="selectedSource.name + ' · 可用知识库'" width="900px">
      <el-alert title="连接成功。选择要发放的知识库，系统会自动记录知识库 ID。" type="success" :closable="false" />
      <el-table :data="collections" row-key="id">
        <el-table-column prop="name" label="知识库名称" min-width="180" />
        <el-table-column prop="id" label="知识库 ID" min-width="220" show-overflow-tooltip />
        <el-table-column prop="description" label="介绍" min-width="200" show-overflow-tooltip />
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button size="small" :disabled="alreadyPublished(row.id)" @click="publishCollection(row)">{{ alreadyPublished(row.id) ? '已添加' : '添加到目录' }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!collections.length" description="该来源尚无可访问的知识库，请检查读取权限和来源内容" />
    </el-dialog>

    <el-dialog v-model="resourceDialog" :title="resourceForm.id ? '编辑发布资料' : '发布资料'" width="560px" :close-on-click-modal="false">
      <el-form label-width="100px" @submit.prevent="saveResource">
        <el-form-item label="资料名称" required><el-input v-model="resourceForm.name" maxlength="120" /></el-form-item>
        <el-form-item label="资料来源">{{ sourceName(resourceForm.source_id) }}</el-form-item>
        <el-form-item label="知识库 ID">{{ resourceForm.remote_id }}</el-form-item>
        <p class="library-hint">自动使用资料来源的地址和 API Key，连同此知识库 ID 下发给 ThinkBuddy，无需重复配置。</p>
        <el-form-item label="介绍"><el-input v-model="resourceForm.description" type="textarea" :rows="4" maxlength="4000" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="resourceForm.sort_num" :min="0" :max="100000" /></el-form-item>
        <el-form-item label="立即上架"><el-switch v-model="resourceForm.published" /></el-form-item>
        <p class="library-hint">上架整个知识库后，后续新增的文档也会向用户提供。下架停止继续发放；停止已有连接的访问需在 WeKnora 中撤销对应 Key。</p>
      </el-form>
      <template #footer>
        <el-button @click="resourceDialog = false">取消</el-button>
        <el-button type="primary" :loading="busy === 'save-resource'" @click="saveResource">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { httpDelete, httpGet, httpPost } from '@/utils/http'

const sources = ref([])
const resources = ref([])
const collections = ref([])
const selectedSource = ref({})
const loading = ref(false)
const busy = ref('')
const tab = ref('resources')
const sourceDialog = ref(false)
const resourceDialog = ref(false)
const collectionDialog = ref(false)
const sourceForm = ref({})
const resourceForm = ref({})
const prefix = '/api/admin/library'
const sourceName = (id) => sources.value.find((source) => source.id === id)?.name || '来源不可用'
const alreadyPublished = (id) => resources.value.some((resource) => resource.source_id === selectedSource.value.id && resource.remote_id === id)
const showError = (error) => ElMessage.error(error.message || '操作失败，请重试')

async function load() {
  loading.value = true
  try {
    const [sourceResponse, resourceResponse] = await Promise.all([httpGet(prefix + '/sources'), httpGet(prefix + '/resources')])
    sources.value = sourceResponse.data || []
    resources.value = resourceResponse.data || []
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
}

function editSource(row = {}) {
  const { id, name = '', base_url = '', enabled = true, has_api_key = false } = row
  sourceForm.value = { id, name, base_url, enabled, has_api_key, api_key: '' }
  sourceDialog.value = true
}

async function saveSource() {
  if (busy.value) return
  if (!sourceForm.value.name?.trim() || !sourceForm.value.base_url?.trim()) {
    ElMessage.warning('请填写来源名称和服务器地址')
    return
  }
  busy.value = 'save-source'
  try {
    await httpPost(prefix + '/sources', sourceForm.value)
    sourceDialog.value = false
    sourceForm.value.api_key = ''
    ElMessage.success('资料来源已保存')
    await load()
  } catch (error) {
    showError(error)
  } finally {
    busy.value = ''
  }
}

async function testSource(row) {
  if (busy.value) return
  busy.value = 'test-' + row.id
  try {
    const response = await httpGet(prefix + '/sources/' + row.id + '/collections')
    selectedSource.value = row
    collections.value = response.data || []
    collectionDialog.value = true
  } catch (error) {
    showError(error)
  } finally {
    busy.value = ''
  }
}

function publishCollection(collection) {
  resourceForm.value = { source_id: selectedSource.value.id, remote_id: collection.id, name: collection.name, description: collection.description || '', sort_num: 0, published: true }
  resourceDialog.value = true
}

function editResource(row) {
  resourceForm.value = { ...row }
  resourceDialog.value = true
}

async function saveResource() {
  if (busy.value) return
  if (!resourceForm.value.name?.trim()) { ElMessage.warning('请填写资料名称'); return }
  busy.value = 'save-resource'
  try {
    await httpPost(prefix + '/resources', resourceForm.value)
    resourceDialog.value = false
    collectionDialog.value = false
    tab.value = 'resources'
    ElMessage.success('发布目录已更新')
    await load()
  } catch (error) {
    showError(error)
  } finally {
    busy.value = ''
  }
}

async function setPublished(row, published) {
  if (busy.value) return
  busy.value = 'publish-' + row.id
  try {
    await httpPost(prefix + '/resources', { ...row, published })
    row.published = published
    ElMessage.success(published ? '资料已上架' : '资料已下架')
  } catch (error) {
    showError(error)
  } finally {
    busy.value = ''
  }
}

async function remove(kind, row) {
  if (busy.value) return
  try {
    await ElMessageBox.confirm('确定删除“' + row.name + '”？已保存到用户设备的资料会保留。', '删除确认', { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' })
    busy.value = 'delete'
    await httpDelete(prefix + '/' + kind + '/' + row.id)
    ElMessage.success('已删除')
    await load()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') showError(error)
  } finally {
    busy.value = ''
  }
}

onMounted(load)
</script>

<style scoped>
.library-heading { display: flex; align-items: center; justify-content: space-between; gap: 16px; margin-bottom: 20px; }
.library-heading h2 { margin: 0 0 8px; font-size: 20px; }
.library-heading p, .library-hint { color: var(--el-text-color-secondary); font-size: 13px; line-height: 1.6; }
.source-table { margin-top: 16px; }
</style>
