<template>
  <div class="basic-config form" v-loading="loading">
    <div class="container">
      <el-form
        :model="system"
        label-position="top"
        class="py-3 px-5"
        ref="systemFormRef"
        :rules="rules"
      >
        <div class="basic-config-form">
          <el-form-item label="网站标题" prop="title">
            <el-input v-model="system['title']" />
          </el-form-item>
          <el-form-item label="控制台标题" prop="admin_title">
            <el-input v-model="system['admin_title']" />
          </el-form-item>
          <el-form-item label="网站Slogan" prop="slogan">
            <el-input v-model="system['slogan']" />
          </el-form-item>
          <el-form-item label="圆形 LOGO" prop="logo">
            <el-input v-model="system['logo']" placeholder="正方形或者圆形 Logo">
              <template #append>
                <el-upload
                  :auto-upload="true"
                  :show-file-list="false"
                  @click="beforeUpload('logo')"
                  :http-request="uploadImg"
                >
                  <el-icon class="uploader-icon">
                    <UploadFilled />
                  </el-icon>
                </el-upload>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="条形 LOGO" prop="logo">
            <el-input v-model="system['bar_logo']" placeholder="长方形 Logo">
              <template #append>
                <el-upload
                  :auto-upload="true"
                  :show-file-list="false"
                  @click="beforeUpload('bar_logo')"
                  :http-request="uploadImg"
                >
                  <el-icon class="uploader-icon">
                    <UploadFilled />
                  </el-icon>
                </el-upload>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="label-title">
                首页导航菜单
                <span class="text-xs text-gray-500">（被选中的菜单将会在首页导航栏显示）</span>
              </div>
            </template>
            <el-select
              v-model="system['index_navs']"
              multiple
              :filterable="true"
              placeholder="请选择菜单，多选"
              style="width: 100%"
            >
              <el-option v-for="item in menus" :key="item.id" :label="item.name" :value="item.id" />
            </el-select>
          </el-form-item>

          <el-form-item label="版权信息" prop="copyright">
            <el-input v-model="system['copyright']" placeholder="更改此选项需要获取 License 授权" />
          </el-form-item>

          <el-form-item label="默认昵称" prop="default_nickname">
            <el-input v-model="system['default_nickname']" placeholder="默认昵称" />
          </el-form-item>

          <el-form-item label="ICP 备案号" prop="icp">
            <el-input v-model="system['icp']" placeholder="请输入 ICP 备案号" />
          </el-form-item>

          <el-form-item>
            <template #label>
              <div class="label-title">
                开放注册
                <span class="text-xs text-gray-500">（关闭注册之后只能通过管理后台添加用户）</span>
              </div>
            </template>
            <el-switch v-model="system['enabled_register']" />
          </el-form-item>

          <el-form-item label="注册方式" prop="register_ways">
            <el-checkbox-group v-model="system['register_ways']">
              <el-checkbox value="mobile">手机注册</el-checkbox>
              <el-checkbox value="email">邮箱注册</el-checkbox>
              <el-checkbox value="username">用户名注册</el-checkbox>
            </el-checkbox-group>
          </el-form-item>

          <el-form-item label="邮件域名白名单" prop="email_white_list">
            <items-input v-model="system['email_white_list']" />
          </el-form-item>

          <el-form-item label="微信客服二维码" prop="wechat_card_url">
            <el-input v-model="system['wechat_card_url']" placeholder="微信客服二维码">
              <template #append>
                <el-upload
                  :auto-upload="true"
                  :show-file-list="false"
                  @click="beforeUpload('wechat_card_url')"
                  :http-request="uploadImg"
                >
                  <el-icon class="uploader-icon">
                    <UploadFilled />
                  </el-icon>
                </el-upload>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="label-title">
                系统辅助AI模型
                <span class="text-xs text-gray-500"
                  >（用来辅助用户生成提示词，翻译的AI模型，默认使用 gpt-4o-mini）</span
                >
              </div>
            </template>
            <el-select
              v-model.number="system['assistant_model_id']"
              :filterable="true"
              placeholder="选择一个系统辅助AI模型"
              style="width: 100%"
            >
              <el-option
                v-for="item in models"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="开启聊天上下文">
            <el-switch v-model="system['enable_context']" />
          </el-form-item>
          <el-form-item label="会话上下文深度">
            <div class="tip-input-line">
              <el-input-number v-model="system['context_deep']" :min="0" :max="10" />
              <div class="tip">
                会话上下文深度：在老会话中继续会话，默认加载多少条聊天记录作为上下文。如果设置为 0
                则不加载聊天记录，仅仅使用当前角色的上下文。该配置参数必须设置需要为偶数。
              </div>
            </div>
          </el-form-item>

          <el-form-item>
            <template #label>
              <div class="label-title">
                SD反向提示词
                <span class="text-xs text-gray-500">（Stable-Diffusion 绘画默认反向提示词）</span>
              </div>
            </template>
            <el-input type="textarea" :rows="2" v-model="system['sd_neg_prompt']" placeholder="" />
          </el-form-item>

          <el-form-item label="会员充值说明" prop="vip_info_text">
            <template #label>
              <div class="label-title">
                会员充值说明
                <span class="text-xs text-gray-500">（会员充值页面的充值说明文字）</span>
              </div>
            </template>
            <el-input type="textarea" :rows="2" v-model="system['vip_info_text']" placeholder="" />
          </el-form-item>

          <el-form-item label="MJ默认API模式" prop="mj_mode">
            <el-select v-model="system['mj_mode']" placeholder="请选择模式">
              <el-option
                v-for="item in mjModels"
                :value="item.value"
                :label="item.name"
                :key="item.value"
                >{{ item.name }}
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item prop="max_file_size">
            <template #label>
              <div class="label-title">
                上传文件限制
                <span class="text-xs text-gray-500">（最大上传文件大小，单位：MB）</span>
              </div>
            </template>
            <el-input
              v-model.number="system['max_file_size']"
              placeholder="最大上传文件大小，单位：MB"
            />
          </el-form-item>
        </div>

        <div style="padding: 10px">
          <el-form-item>
            <el-button type="primary" @click="save">保存</el-button>
          </el-form-item>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import ItemsInput from '@/components/ui/ItemsInput.vue'
import { httpGet, httpPost } from '@/utils/http'
import { copyObj } from '@/utils/libs'
import { UploadFilled } from '@element-plus/icons-vue'
import Compressor from 'compressorjs'
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'

const system = ref({ models: [] })
const configBak = ref({})
const loading = ref(true)
const systemFormRef = ref(null)
const models = ref([])
const menus = ref([])
const mjModels = ref([
  { name: '慢速（Relax）', value: 'relax' },
  { name: '快速（Fast）', value: 'fast' },
  { name: '急速（Turbo）', value: 'turbo' },
])

onMounted(() => {
  // 加载系统配置
  httpGet('/api/admin/config/get?key=system')
    .then((res) => {
      system.value = res.data
      configBak.value = copyObj(system.value)
    })
    .catch((e) => {
      ElMessage.error('加载系统配置失败: ' + e.message)
    })
    .finally(() => {
      loading.value = false
    })

  httpGet('/api/admin/model/list')
    .then((res) => {
      models.value = res.data
    })
    .catch((e) => {
      ElMessage.error('获取模型失败：' + e.message)
    })

  httpGet('/api/admin/menu/list')
    .then((res) => {
      menus.value = res.data
    })
    .catch((e) => {
      ElMessage.error('获取菜单失败：' + e.message)
    })
})

const rules = reactive({
  title: [{ required: true, message: '请输入网站标题', trigger: 'blur' }],
  admin_title: [{ required: true, message: '请输入控制台标题', trigger: 'blur' }],
})

const save = function () {
  systemFormRef.value.validate((valid) => {
    if (valid) {
      httpPost('/api/admin/config/update/base', system.value)
        .then(() => {
          ElMessage.success('操作成功！')
        })
        .catch((e) => {
          ElMessage.error('操作失败：' + e.message)
        })
    }
  })
}

const configKey = ref('')
const beforeUpload = (key) => {
  configKey.value = key
}

// 图片上传
const uploadImg = (file) => {
  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData()
      formData.append('file', result, result.name)
      // 执行上传操作
      httpPost('/api/admin/upload', formData)
        .then((res) => {
          system.value[configKey.value] = res.data.url
          ElMessage.success('上传成功')
        })
        .catch((e) => {
          ElMessage.error('上传失败:' + e.message)
        })
    },
    error(e) {
      ElMessage.error('上传失败:' + e.message)
    },
  })
}
</script>

<style lang="scss">
@use '@/assets/css/admin/form.scss' as *;
@use '@/assets/css/main.scss' as *;

.basic-config {
  display: flex;
  justify-content: center;
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
}
</style>
