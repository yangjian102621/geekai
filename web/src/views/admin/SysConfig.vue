<template>
  <div class="system-config" v-loading="loading">

    <el-tabs v-model="activeName" class="sys-tabs">
      <el-tab-pane label="系统配置" name="basic">
        <div class="container">
          <el-form :model="system" label-width="150px" label-position="right" ref="systemFormRef" :rules="rules">
            <el-tabs type="border-card">
              <el-tab-pane label="基础配置">
                <el-form-item label="网站标题" prop="title">
                  <el-input v-model="system['title']"/>
                </el-form-item>
                <el-form-item label="控制台标题" prop="admin_title">
                  <el-input v-model="system['admin_title']"/>
                </el-form-item>
                <el-form-item label="网站Slogan" prop="slogan">
                  <el-input v-model="system['slogan']"/>
                </el-form-item>
                <el-form-item label="网站 LOGO" prop="logo">
                  <el-input v-model="system['logo']" placeholder="网站LOGO图片">
                    <template #append>
                      <el-upload
                          :auto-upload="true"
                          :show-file-list="false"
                          @click="beforeUpload('logo')"
                          :http-request="uploadImg"
                      >
                        <el-icon class="uploader-icon">
                          <UploadFilled/>
                        </el-icon>
                      </el-upload>
                    </template>
                  </el-input>
                </el-form-item>

                <el-form-item label="首页背景图" prop="logo">
                  <div class="tip-input">
                    <el-input v-model="system['index_bg_url']" placeholder="网站首页背景图片">
                      <template #append>
                        <el-upload
                            :auto-upload="true"
                            :show-file-list="false"
                            @click="beforeUpload('index_bg_url')"
                            :http-request="uploadImg"
                        >
                          <el-icon class="uploader-icon">
                            <UploadFilled/>
                          </el-icon>
                        </el-upload>
                      </template>
                    </el-input>
                    <el-button type="primary" @click="system.index_bg_url = 'https://api.dujin.org/bing/1920.php'">使用动态背景</el-button>
                    <el-button @click="system.index_bg_url = 'color'">使用纯色背景</el-button>
                  </div>
                </el-form-item>

                <el-form-item label="首页导航菜单" prop="index_navs">
                  <div class="tip-input">
                    <el-select
                        v-model="system['index_navs']"
                        multiple
                        :filterable="true"
                        placeholder="请选择菜单，多选"
                        style="width: 100%"
                    >
                      <el-option
                          v-for="item in menus"
                          :key="item.id"
                          :label="item.name"
                          :value="item.id"
                      />
                    </el-select>
                    <div class="info">
                      <el-tooltip
                          class="box-item"
                          effect="dark"
                          content="被选中的菜单将会在首页导航栏显示"
                          placement="right"
                      >
                        <el-icon>
                          <InfoFilled/>
                        </el-icon>
                      </el-tooltip>
                    </div>
                  </div>
                </el-form-item>

                <el-form-item label="版权信息" prop="copyright">
                  <el-input v-model="system['copyright']" placeholder="更改此选项需要获取 License 授权"/>
                </el-form-item>

                <el-form-item label="开放注册" prop="enabled_register">
                  <div class="tip-input">
                    <el-switch v-model="system['enabled_register']"/>
                    <div class="info">
                      <el-tooltip
                          effect="dark"
                          content="关闭注册之后只能通过管理后台添加用户"
                          raw-content
                          placement="right"
                      >
                        <el-icon>
                          <InfoFilled/>
                        </el-icon>
                      </el-tooltip>
                    </div>
                  </div>
                </el-form-item>

                <el-form-item label="启用验证码" prop="enabled_verify">
                  <div class="tip-input">
                    <el-switch v-model="system['enabled_verify']"/>
                    <div class="info">
                      <el-tooltip
                          effect="dark"
                          content="启用验证码之后，注册登录都会加载行为验证码，增加安全性。此功能需要购买验证码服务才会生效。"
                          raw-content
                          placement="right"
                      >
                        <el-icon>
                          <InfoFilled/>
                        </el-icon>
                      </el-tooltip>
                    </div>
                  </div>
                </el-form-item>

                <el-form-item label="注册方式" prop="register_ways">
                  <el-checkbox-group v-model="system['register_ways']">
                    <el-checkbox value="mobile">手机注册</el-checkbox>
                    <el-checkbox value="email">邮箱注册</el-checkbox>
                    <el-checkbox value="username">用户名注册</el-checkbox>
                  </el-checkbox-group>
                </el-form-item>

                <el-form-item label="邮件域名白名单" prop="register_ways">
                  <items-input v-model:value="system['email_white_list']"/>
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
                          <UploadFilled/>
                        </el-icon>
                      </el-upload>
                    </template>
                  </el-input>
                </el-form-item>
                <el-form-item label="默认翻译模型">
                  <template #default>
                    <div class="tip-input">
                      <el-select
                          v-model.number="system['translate_model_id']"
                          :filterable="true"
                          placeholder="选择一个默认模型来翻译提示词"
                          style="width: 100%"
                      >
                        <el-option
                            v-for="item in models"
                            :key="item.id"
                            :label="item.name"
                            :value="item.id"
                        />
                      </el-select>
                      <div class="info">
                        <el-tooltip
                            class="box-item"
                            effect="dark"
                            content="新用户注册默认开通的 AI 模型"
                            placement="right"
                        >
                          <el-icon>
                            <InfoFilled/>
                          </el-icon>
                        </el-tooltip>
                      </div>
                    </div>
                  </template>
                </el-form-item>

                <el-form-item label="开启聊天上下文">
                  <el-switch v-model="system['enable_context']"/>
                </el-form-item>
                <el-form-item label="会话上下文深度">
                  <div class="tip-input-line">
                    <el-input-number v-model="system['context_deep']" :min="0" :max="10"/>
                    <div class="tip">会话上下文深度：在老会话中继续会话，默认加载多少条聊天记录作为上下文。如果设置为
                      0
                      则不加载聊天记录，仅仅使用当前角色的上下文。该配置参数最好设置需要为偶数，否则将无法兼容百度的 API。
                    </div>
                  </div>
                </el-form-item>

                <el-form-item label="SD反向提示词" prop="sd_neg_prompt">
                  <div class="tip-input">
                    <el-input v-model="system['sd_neg_prompt']" placeholder=""/>
                    <div class="info">
                      <el-tooltip
                          effect="dark"
                          content="Stable-Diffusion 绘画默认反向提示词"
                          raw-content
                          placement="right"
                      >
                        <el-icon>
                          <InfoFilled/>
                        </el-icon>
                      </el-tooltip>
                    </div>
                  </div>
                </el-form-item>

                <el-form-item label="会员充值说明" prop="order_pay_timeout">
                  <div class="tip-input">
                    <el-input v-model="system['vip_info_text']" placeholder=""/>
                    <div class="info">
                      <el-tooltip
                          effect="dark"
                          content="会员充值页面的充值说明文字"
                          raw-content
                          placement="right"
                      >
                        <el-icon>
                          <InfoFilled/>
                        </el-icon>
                      </el-tooltip>
                    </div>
                  </div>
                </el-form-item>

                <el-form-item label="MJ默认API模式" prop="mj_mode">
                  <el-select v-model="system['mj_mode']" placeholder="请选择模式">
                    <el-option v-for="item in mjModels" :value="item.value" :label="item.name" :key="item.value">{{
                        item.name
                      }}
                    </el-option>
                  </el-select>
                </el-form-item>
              </el-tab-pane>

              <el-tab-pane label="算力配置">
                <el-form-item label="注册赠送算力" prop="init_power">
                  <el-input v-model.number="system['init_power']" placeholder="新用户注册赠送算力"/>
                </el-form-item>
                <el-form-item label="邀请赠送算力" prop="invite_power">
                  <el-input v-model.number="system['invite_power']" placeholder="邀请新用户注册赠送算力"/>
                </el-form-item>
                <el-form-item label="VIP每月赠送算力" prop="vip_month_power">
                  <el-input v-model.number="system['vip_month_power']" placeholder="VIP用户每月赠送算力"/>
                </el-form-item>
                <el-form-item label="每日赠送算力" prop="daily_power">
                  <div class="tip-input-line">
                    <el-input v-model.number="system['daily_power']" placeholder="默认值0"/>
                    <div class="tip">
                      如果设置0表示不赠送，用户享受完免费算力额度之后就不能再发起对话了。如果设置为N，则系统每天将算力值小于N的用户自动补充到N。注意，此功能要配合XXL-JOB启用。
                    </div>
                  </div>
                </el-form-item>
                <el-form-item label="MJ绘图算力" prop="mj_power">
                  <div class="tip-input">
                    <el-input v-model.number="system['mj_power']" placeholder=""/>
                    <div class="info">
                      <el-tooltip
                          effect="dark"
                          content="使用MidJourney画一张图消耗算力"
                          raw-content
                          placement="right"
                      >
                        <el-icon>
                          <InfoFilled/>
                        </el-icon>
                      </el-tooltip>
                    </div>
                  </div>
                </el-form-item>
                <el-form-item label="MJ操作算力" prop="mj_action_power">
                  <div class="tip-input">
                    <el-input v-model.number="system['mj_action_power']" placeholder=""/>
                    <div class="info">
                      <el-tooltip
                          effect="dark"
                          content="放大，变换，重绘操作一次消耗的算力"
                          raw-content
                          placement="right"
                      >
                        <el-icon>
                          <InfoFilled/>
                        </el-icon>
                      </el-tooltip>
                    </div>
                  </div>
                </el-form-item>
                <el-form-item label="Stable-Diffusion算力" prop="sd_power">
                  <el-input v-model.number="system['sd_power']" placeholder="使用Stable-Diffusion画一张图消耗算力"/>
                </el-form-item>
                <el-form-item label="DALL-E-3算力" prop="dall_power">
                  <el-input v-model.number="system['dall_power']" placeholder="使用DALL-E-3画一张图消耗算力"/>
                </el-form-item>
                <el-form-item label="Suno 算力" prop="suno_power">
                  <el-input v-model.number="system['suno_power']" placeholder="使用 Suno 生成一首音乐消耗算力"/>
                </el-form-item>
                <el-form-item label="Luma 算力" prop="luma_power">
                  <el-input v-model.number="system['luma_power']" placeholder="使用 Luma 生成一段视频消耗算力"/>
                </el-form-item>
              </el-tab-pane>
            </el-tabs>

            <div style="padding: 10px;">
              <el-form-item>
                <el-button type="primary" @click="save('system')">保存</el-button>
              </el-form-item>
            </div>
          </el-form>
        </div>
      </el-tab-pane>
      <el-tab-pane label="公告配置" name="notice">
        <md-editor class="mgb20" v-model="notice" @on-upload-img="onUploadImg"/>
        <el-form-item>
          <div style="padding-top: 10px;margin-left: 150px;">
            <el-button type="primary" @click="save('notice')">保存</el-button>
          </div>
        </el-form-item>
      </el-tab-pane>
      <el-tab-pane label="思维导图" name="mark_map">
        <md-editor class="mgb20" v-model="system['mark_map_text']" @on-upload-img="onUploadImg"/>
        <el-form-item>
          <div style="padding-top: 10px;margin-left: 150px;">
            <el-button type="primary" @click="save('system')">保存</el-button>
          </div>
        </el-form-item>
      </el-tab-pane>
      <el-tab-pane label="菜单配置" name="menu">
        <Menu/>
      </el-tab-pane>

      <el-tab-pane label="授权激活" name="license">
        <div class="container">
          <el-descriptions
              v-if="license.is_active"
              class="margin-top"
              title="已授权信息"
              :column="1"
              border
          >
            <el-descriptions-item>
              <template #label>
                <div class="cell-item">License Key</div>
              </template>
              {{ license.key }}
            </el-descriptions-item>
            <el-descriptions-item>
              <template #label>
                <div class="cell-item">机器码</div>
              </template>
              {{ license.machine_id }}
            </el-descriptions-item>
            <el-descriptions-item>
              <template #label>
                <div class="cell-item">到期时间</div>
              </template>
              {{ dateFormat(license.expired_at) }}
            </el-descriptions-item>
            <el-descriptions-item>
              <template #label>
                <div class="cell-item">用户人数</div>
              </template>
              {{ license.configs?.user_num }}
            </el-descriptions-item>
            <el-descriptions-item>
              <template #label>
                <div class="cell-item">去版权</div>
              </template>
              <el-icon class="selected" v-if="license.configs?.de_copy"><Select /></el-icon>
              <el-icon class="closed" v-else><CloseBold /></el-icon>
              <span class="text">去版权之后前端页面将不会显示版权信息和源码地址</span>
            </el-descriptions-item>
          </el-descriptions>

          <h3>激活后可获得以下权限：</h3>
          <ol class="active-info">
            <li>1、使用任意第三方中转 API KEY，而不用局限于 GeekAI 推荐的白名单列表</li>
            <li>2、可以在相关页面去除 GeekAI 的版权信息，或者修改为自己的版权信息</li>
          </ol>

          <el-form :model="system" label-width="150px" label-position="right">
            <el-form-item label="许可授权码" prop="license">
              <el-input v-model="licenseKey"/>
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="active">立即激活</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>

      <el-tab-pane label="修复数据" name="fixData">
        <div class="container">
<!--          <p class="text">有些版本升级的时候更新了数据库的结构，比如字段名字改了，需要把之前的字段的值转移到其他字段，这些无法通过简单的 SQL 语句可以实现的，需要手动写程序修正数据。</p>-->

<!--          <p class="text">当前版本 v4.1.4 需要修正用户数据，增加了 mobile 和 email 字段，需要把之前用手机号或者邮箱注册的用户的 username 字段数据初始化到 mobile 或者 email 字段。另外，需要把订单的支付渠道从名字称修正为 key。</p>-->

<!--          <el-text type="danger">请注意：在修复数据前，请先备份好数据库，以免数据丢失！</el-text>-->

          <p><el-button type="primary" @click="fixData">立即修复</el-button></p>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import {onMounted, reactive, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import Compressor from "compressorjs";
import {ElMessage, ElMessageBox} from "element-plus";
import {InfoFilled, UploadFilled,Select,CloseBold} from "@element-plus/icons-vue";
import MdEditor from "md-editor-v3";
import 'md-editor-v3/lib/style.css';
import Menu from "@/views/admin/Menu.vue";
import {copyObj, dateFormat} from "@/utils/libs";
import ItemsInput from "@/components/ui/ItemsInput.vue";

const activeName = ref('basic')
const system = ref({models: []})
const configBak = ref({})
const loading = ref(true)
const systemFormRef = ref(null)
const models = ref([])
const notice = ref("")
const license = ref({is_active: false})
const menus = ref([])
const mjModels = ref([
  {name: "慢速（Relax）", value: "relax"},
  {name: "快速（Fast）", value: "fast"},
  {name: "急速（Turbo）", value: "turbo"},
])

onMounted(() => {
  // 加载系统配置
  httpGet('/api/admin/config/get?key=system').then(res => {
    system.value = res.data
    configBak.value = copyObj(system.value)
  }).catch(e => {
    ElMessage.error("加载系统配置失败: " + e.message)
  })
  // 加载聊天配置
  httpGet('/api/admin/config/get?key=notice').then(res => {
    notice.value = res.data['content']
  }).catch(e => {
    ElMessage.error("公告信息失败: " + e.message)
  })

  httpGet('/api/admin/model/list').then(res => {
    models.value = res.data
    loading.value = false
  }).catch(e => {
    ElMessage.error("获取模型失败：" + e.message)
  })

  httpGet('/api/admin/menu/list').then(res => {
    menus.value = res.data
  }).catch(e => {
    ElMessage.error("获取模型失败：" + e.message)
  })

  fetchLicense()
})

const fetchLicense = () => {
  httpGet("/api/admin/config/license").then(res => {
    license.value = res.data
  }).catch(e => {
    ElMessage.error("获取 License 失败：" + e.message)
  })
}

const rules = reactive({
  title: [{required: true, message: '请输入网站标题', trigger: 'blur',}],
  admin_title: [{required: true, message: '请输入控制台标题', trigger: 'blur',}],
  init_chat_calls: [{required: true, message: '请输入赠送对话次数', trigger: 'blur'}],
  user_img_calls: [{required: true, message: '请输入赠送绘图次数', trigger: 'blur'}],
})
const save = function (key) {
  if (key === 'system') {
    systemFormRef.value.validate((valid) => {
      if (valid) {
        httpPost('/api/admin/config/update', {key: key, config: system.value, config_bak: configBak.value}).then(() => {
          ElMessage.success("操作成功！")
        }).catch(e => {
          ElMessage.error("操作失败：" + e.message)
        })
      }
    })
  } else if (key === 'notice') {
    httpPost('/api/admin/config/update', {key: key, config: {content: notice.value, updated: true}}).then(() => {
      ElMessage.success("操作成功！")
    }).catch(e => {
      ElMessage.error("操作失败：" + e.message)
    })
  }
}

// 激活授权
const licenseKey = ref("")
const active = () => {
  if (licenseKey.value === "") {
    return ElMessage.error("请输入授权码")
  }
  httpPost("/api/admin/config/active", {license: licenseKey.value}).then(res => {
    ElMessage.success("授权成功，机器编码为：" + res.data)
    fetchLicense()
  }).catch(e => {
    ElMessage.error(e.message)
  })
}

const configKey = ref("")
const beforeUpload = (key) => {
  configKey.value = key
}

// 图片上传
const uploadImg = (file) => {
  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData();
      formData.append('file', result, result.name);
      // 执行上传操作
      httpPost('/api/admin/upload', formData).then((res) => {
        system.value[configKey.value] = res.data.url
        ElMessage.success('上传成功')
      }).catch((e) => {
        ElMessage.error('上传失败:' + e.message)
      })
    },
    error(e) {
      ElMessage.error('上传失败:' + e.message)
    },
  });
};

// 编辑期文件上传处理
const onUploadImg = (files, callback) => {
  Promise.all(
      files.map((file) => {
        return new Promise((rev, rej) => {
          const formData = new FormData();
          formData.append('file', file, file.name);
          // 执行上传操作
          httpPost('/api/admin/upload', formData).then((res) => rev(res)).catch((error) => rej(error));
        });
      })
  ).then(res => {
    ElMessage.success({message: "上传成功", duration: 500})
    callback(res.map((item) => item.data.url));
  }).catch(e => {
    ElMessage.error('图片上传失败:' + e.message)
  })
};

const fixData = () => {

  ElMessageBox.confirm(
      '在修复数据前，请先备份好数据库，以免数据丢失！是否继续操作?',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  ).then(() => {
    loading.value = true
    httpGet("/api/admin/config/fixData").then(() => {
      ElMessage.success("数据修复成功")
      loading.value = false
    }).catch(e => {
      loading.value = false
      ElMessage.error("数据修复失败：" + e.message)
    })
  })
}


</script>

<style lang="stylus" scoped>
@import "@/assets/css/admin/form.styl"
.system-config {
  display flex
  justify-content center

  .sys-tabs {
    width 100%
    background-color var(--el-bg-color)
    padding 10px 20px 40px 20px
    //border: 1px solid var(--el-border-color);

    .container {
      .el-form {
        .el-form-item__content {

          .tip-text {
            padding-left 10px;
          }

          .el-icon {
            font-size 16px
            cursor pointer
          }

          .uploader-icon {
            font-size 24px
            position relative
            top 3px
          }

          .tip-input-line {
            .tip {
              margin-top 10px
              color #c1c1c1
              font-size 12px;
              line-height 1.5;
            }
          }
        }

        .el-input {
          width 100%
        }
      }


      .text {
        font-size 14px
      }

      .active-info {
        line-height 1.5
        padding 10px 0 30px 0
      }
      .el-descriptions {
        margin-bottom 20px
        .el-icon {
          font-size 18px
        }
        .selected {
          color #0bc15f
        }

        .closed {
          color #da0d54
        }
        .text {
          margin-left 10px
          font-size 12px
          color #999999
          position: relative;
          top -5px
        }
      }

      .el-alert {
        margin-bottom 15px;
      }
    }

  }
}
</style>