<template>
  <div class="system-config form">
    <div class="container">
      <el-tabs type="border-card" v-model="activeName">
        <el-tab-pane label="微信公众号配置" name="gzh">
          <el-form
            label-width="150px"
            label-position="top"
            :model="gzhConfig"
            ref="wechatFormRef"
            :rules="rules"
          >
            <el-form-item label="微信公众号 AppID" prop="app_id">
              <el-input v-model="gzhConfig['app_id']" />
            </el-form-item>
            <el-form-item label="微信公众号 AppSecret" prop="secret">
              <el-input v-model="gzhConfig['secret']" />
            </el-form-item>
            <el-form-item label="微信公众号 Token" prop="token">
              <el-input v-model="gzhConfig['token']" />
            </el-form-item>
            <el-form-item label="微信公众号 EncodingAESKey" prop="encoding_aes_key">
              <el-input v-model="gzhConfig['encoding_aes_key']" />
            </el-form-item>
            <el-form-item label="启用微信公众号登录" prop="enabled">
              <el-switch v-model="gzhConfig['enabled']" />
            </el-form-item>
          </el-form>
          <div class="tab-actions">
            <el-button type="primary" @click="save">保存</el-button>
          </div>
        </el-tab-pane>

        <el-tab-pane label="自定义菜单" name="menu">
          <el-alert type="info" show-icon :closable="false" class="menu-hint">
            <template #title>
              菜单将同步到微信公众平台。跳转链接（view）的域名需在公众号后台配置为「业务域名」或通过授权页打开；小程序类型需填写 appid、pagepath
              与备用网页 url。
            </template>
          </el-alert>

          <div v-for="(item, i) in menuTops" :key="i" class="menu-top-card">
            <div class="menu-top-head">
              <span class="menu-top-label">一级菜单 {{ i + 1 }}</span>
              <el-button type="danger" link @click="removeTop(i)">删除</el-button>
            </div>
            <el-form label-position="top">
              <el-form-item label="菜单名称">
                <el-input v-model="item.name" placeholder="显示在公众号底部的名称" maxlength="16" show-word-limit />
              </el-form-item>
              <el-form-item label="类型">
                <el-radio-group v-model="item.isSubmenu">
                  <el-radio :label="false">直接响应（链接 / 点击）</el-radio>
                  <el-radio :label="true">展开子菜单（最多 5 个）</el-radio>
                </el-radio-group>
              </el-form-item>

              <template v-if="!item.isSubmenu">
                <el-form-item label="动作类型">
                  <el-select v-model="item.type" style="width: 220px">
                    <el-option label="跳转网页 view" value="view" />
                    <el-option label="点击推事件 click" value="click" />
                    <el-option label="打开小程序 miniprogram" value="miniprogram" />
                  </el-select>
                </el-form-item>
                <el-form-item v-if="item.type === 'view'" label="网页链接 URL">
                  <el-input v-model="item.url" placeholder="https:// 开头" />
                </el-form-item>
                <el-form-item v-if="item.type === 'click'" label="事件 Key">
                  <el-input v-model="item.key" placeholder="与服务器事件匹配，勿与现有 key 冲突" />
                </el-form-item>
                <template v-if="item.type === 'miniprogram'">
                  <el-form-item label="小程序 AppID">
                    <el-input v-model="item.appid" />
                  </el-form-item>
                  <el-form-item label="小程序页面路径 pagepath">
                    <el-input v-model="item.pagepath" placeholder="如 pages/index/index" />
                  </el-form-item>
                  <el-form-item label="备用网页 url（必填）">
                    <el-input v-model="item.url" placeholder="无法打开小程序时跳转" />
                  </el-form-item>
                </template>
              </template>

              <template v-else>
                <div class="sub-menu-toolbar">
                  <el-button
                    type="primary"
                    plain
                    size="small"
                    :disabled="(item.sub_button || []).length >= 5"
                    @click="addSub(i)"
                  >
                    添加子菜单
                  </el-button>
                </div>
                <div
                  v-for="(sub, j) in item.sub_button"
                  :key="j"
                  class="sub-menu-row"
                >
                  <el-divider content-position="left">子项 {{ j + 1 }}</el-divider>
                  <div class="sub-menu-head">
                    <el-button type="danger" link size="small" @click="removeSub(i, j)">删除</el-button>
                  </div>
                  <el-form-item label="子菜单名称">
                    <el-input v-model="sub.name" maxlength="60" show-word-limit />
                  </el-form-item>
                  <el-form-item label="动作类型">
                    <el-select v-model="sub.type" style="width: 220px">
                      <el-option label="跳转网页 view" value="view" />
                      <el-option label="点击推事件 click" value="click" />
                      <el-option label="打开小程序 miniprogram" value="miniprogram" />
                    </el-select>
                  </el-form-item>
                  <el-form-item v-if="sub.type === 'view'" label="网页链接 URL">
                    <el-input v-model="sub.url" />
                  </el-form-item>
                  <el-form-item v-if="sub.type === 'click'" label="事件 Key">
                    <el-input v-model="sub.key" />
                  </el-form-item>
                  <template v-if="sub.type === 'miniprogram'">
                    <el-form-item label="小程序 AppID">
                      <el-input v-model="sub.appid" />
                    </el-form-item>
                    <el-form-item label="小程序页面路径">
                      <el-input v-model="sub.pagepath" />
                    </el-form-item>
                    <el-form-item label="备用网页 url">
                      <el-input v-model="sub.url" />
                    </el-form-item>
                  </template>
                </div>
              </template>
            </el-form>
          </div>

          <el-button class="add-top-btn" :disabled="menuTops.length >= 3" @click="addTop">
            添加一级菜单（最多 3 个）
          </el-button>

          <div class="tab-actions menu-actions">
            <el-button :loading="menuSaving" @click="saveMenuDraft">保存草稿</el-button>
            <el-button type="primary" :loading="menuPublishing" @click="syncMenuToWechat">
              同步到微信
            </el-button>
            <el-button :loading="menuPulling" @click="pullMenuFromWechat">从微信拉取</el-button>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'

const activeName = ref('gzh')
const gzhConfig = ref({ app_id: '', secret: '' })
const wechatFormRef = ref(null)

const menuTops = ref([])
const menuSaving = ref(false)
const menuPublishing = ref(false)
const menuPulling = ref(false)

function emptySub() {
  return {
    name: '',
    type: 'view',
    url: '',
    key: '',
    appid: '',
    pagepath: '',
  }
}

function emptyTop() {
  return {
    name: '',
    isSubmenu: false,
    type: 'view',
    url: '',
    key: '',
    appid: '',
    pagepath: '',
    sub_button: [],
  }
}

function normalizeLeaf(s) {
  const t = s.type
  let type = 'view'
  if (t === 'click' || t === 'miniprogram') {
    type = t
  } else if (t && t !== 'view') {
    type = s.url ? 'view' : 'click'
  }
  return {
    name: s.name || '',
    type,
    url: s.url || '',
    key: s.key || '',
    appid: s.appid || '',
    pagepath: s.pagepath || '',
  }
}

function normalizeTop(b) {
  if (b.sub_button && b.sub_button.length) {
    return {
      name: b.name || '',
      isSubmenu: true,
      type: 'view',
      url: '',
      key: '',
      appid: '',
      pagepath: '',
      sub_button: b.sub_button.map(normalizeLeaf),
    }
  }
  const l = normalizeLeaf(b)
  return {
    name: b.name || '',
    isSubmenu: false,
    type: l.type,
    url: l.url,
    key: l.key,
    appid: l.appid,
    pagepath: l.pagepath,
    sub_button: [],
  }
}

function applyMenuFromApi(payload) {
  const buttons = payload?.button
  if (!buttons || !buttons.length) {
    menuTops.value = []
    return
  }
  menuTops.value = buttons.map(normalizeTop)
}

function leafToApi(s) {
  const o = { name: s.name, type: s.type }
  if (s.type === 'view') {
    o.url = s.url
  } else if (s.type === 'click') {
    o.key = s.key
  } else if (s.type === 'miniprogram') {
    o.appid = s.appid
    o.pagepath = s.pagepath
    o.url = s.url
  }
  return o
}

function menuTopsToButton() {
  return menuTops.value.map((b) => {
    if (b.isSubmenu) {
      return {
        name: b.name,
        sub_button: (b.sub_button || []).map(leafToApi),
      }
    }
    return leafToApi({
      name: b.name,
      type: b.type,
      url: b.url,
      key: b.key,
      appid: b.appid,
      pagepath: b.pagepath,
    })
  })
}

function addTop() {
  if (menuTops.value.length >= 3) {
    return
  }
  menuTops.value.push(emptyTop())
}

function removeTop(i) {
  menuTops.value.splice(i, 1)
}

function addSub(topIndex) {
  const item = menuTops.value[topIndex]
  if (!item.sub_button) {
    item.sub_button = []
  }
  if (item.sub_button.length >= 5) {
    return
  }
  item.sub_button.push(emptySub())
}

function removeSub(topIndex, subIndex) {
  menuTops.value[topIndex].sub_button.splice(subIndex, 1)
}

async function saveMenuDraft() {
  menuSaving.value = true
  try {
    const body = { button: menuTopsToButton() }
    await httpPost('/api/admin/config/update/wx_gzh_menu', body)
    ElMessage.success('菜单草稿已保存')
  } catch (e) {
    ElMessage.error('保存失败：' + (e.message || String(e)))
  } finally {
    menuSaving.value = false
  }
}

async function syncMenuToWechat() {
  menuPublishing.value = true
  try {
    const body = { button: menuTopsToButton() }
    await httpPost('/api/admin/config/update/wx_gzh_menu', body)
    await httpPost('/api/admin/config/wx_gzh/menu/publish', {})
    ElMessage.success('已保存并同步到微信')
  } catch (e) {
    ElMessage.error('同步失败：' + (e.message || String(e)))
  } finally {
    menuPublishing.value = false
  }
}

async function pullMenuFromWechat() {
  menuPulling.value = true
  try {
    const res = await httpGet('/api/admin/config/wx_gzh/menu/query')
    applyMenuFromApi(res.data)
    ElMessage.success('已从微信拉取并写入草稿')
  } catch (e) {
    ElMessage.error('拉取失败：' + (e.message || String(e)))
  } finally {
    menuPulling.value = false
  }
}

onMounted(() => {
  Promise.all([
    httpGet('/api/admin/config/get?key=wx_gzh'),
    httpGet('/api/admin/config/get?key=wx_gzh_menu'),
  ])
    .then(([r1, r2]) => {
      gzhConfig.value = r1.data || {}
      applyMenuFromApi(r2.data || {})
    })
    .catch((e) => {
      ElMessage.error('加载系统配置失败: ' + (e.message || String(e)))
    })
})

const rules = reactive({
  app_id: [{ required: true, message: '请输入公众号 AppID', trigger: 'blur' }],
  secret: [{ required: true, message: '请输入公众号 AppSecret', trigger: 'blur' }],
})
const save = function () {
  wechatFormRef.value.validate((valid) => {
    if (valid) {
      httpPost('/api/admin/config/update/wx_gzh', {
        app_id: gzhConfig.value.app_id,
        secret: gzhConfig.value.secret,
        token: gzhConfig.value.token,
        encoding_aes_key: gzhConfig.value.encoding_aes_key,
        enabled: gzhConfig.value.enabled,
      })
        .then(() => {
          ElMessage.success('操作成功！')
        })
        .catch((e) => {
          ElMessage.error('操作失败：' + (e.message || String(e)))
        })
    }
  })
}
</script>

<style lang="css" scoped>
@import '@/assets/css/admin/form.css';
@import '@/assets/css/main.css';

.system-config {
  display: flex;
  justify-content: center;
  .sys-tabs {
    width: 100%;
    background-color: var(--el-bg-color);
    padding: 10px 20px 40px 20px;
  }
}

.menu-hint {
  margin-bottom: 16px;
}

.menu-top-card {
  border: 1px solid var(--el-border-color);
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 16px;
  background: var(--el-fill-color-blank);
}

.menu-top-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.menu-top-label {
  font-weight: 600;
}

.sub-menu-toolbar {
  margin-bottom: 8px;
}

.sub-menu-row {
  padding-left: 8px;
  border-left: 3px solid var(--el-color-primary-light-5);
  margin-bottom: 12px;
}

.sub-menu-head {
  display: flex;
  justify-content: flex-end;
}

.add-top-btn {
  margin-bottom: 20px;
}

.tab-actions {
  padding: 10px 0 0;
}

.menu-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
</style>
