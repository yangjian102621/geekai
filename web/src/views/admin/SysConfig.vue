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

                <el-form-item label="注册方式" prop="register_ways">
                  <el-checkbox-group v-model="system['register_ways']">
                    <el-checkbox value="mobile">手机注册</el-checkbox>
                    <el-checkbox value="email">邮箱注册</el-checkbox>
                    <el-checkbox value="username">用户名注册</el-checkbox>
                  </el-checkbox-group>
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
                <el-form-item label="默认AI模型" prop="default_models">
                  <template #default>
                    <div class="tip-input">
                      <el-select
                          v-model="system['default_models']"
                          multiple
                          :filterable="true"
                          placeholder="选择AI模型，多选"
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
              </el-tab-pane>
              <el-tab-pane label="众筹支付">
                <el-form-item label="启用众筹功能" prop="enabled_reward">
                  <div class="tip-input">
                    <el-switch v-model="system['enabled_reward']"/>
                    <div class="info">
                      <el-tooltip
                          effect="dark"
                          content="如果关闭次功能将不在用户菜单显示众筹二维码"
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

                <div v-if="system['enabled_reward']">
                  <el-form-item label="算力单价" prop="power_price">
                    <el-input v-model="system['power_price']"
                              placeholder="单位算力的价格，比如设置 0.1 表示捐赠1元钱可以得到10个单位算力"/>
                  </el-form-item>
                  <el-form-item label="收款二维码" prop="reward_img">
                    <el-input v-model="system['reward_img']" placeholder="众筹收款二维码地址">
                      <template #append>
                        <el-upload
                            :auto-upload="true"
                            :show-file-list="false"
                            @click="beforeUpload('reward_img')"
                            :http-request="uploadImg"
                        >
                          <el-icon class="uploader-icon">
                            <UploadFilled/>
                          </el-icon>
                        </el-upload>
                      </template>
                    </el-input>
                  </el-form-item>
                </div>

                <el-form-item label="订单超时时间" prop="order_pay_timeout">
                  <div class="tip-input">
                    <el-input v-model.number="system['order_pay_timeout']" placeholder="单位：秒"/>
                    <div class="info">
                      <el-tooltip
                          effect="dark"
                          content="系统会定期清理超时未支付的订单<br/>默认值：900秒"
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

      <el-tab-pane label="菜单配置" name="menu">
        <Menu/>
      </el-tab-pane>

      <el-tab-pane label="AI绘图配置" name="AIDrawing">
        <AIDrawing/>
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
    </el-tabs>
  </div>
</template>

<script setup>
import {onMounted, reactive, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import Compressor from "compressorjs";
import {ElMessage} from "element-plus";
import {InfoFilled, UploadFilled,Select,CloseBold} from "@element-plus/icons-vue";
import MdEditor from "md-editor-v3";
import 'md-editor-v3/lib/style.css';
import Menu from "@/views/admin/Menu.vue";
import {copyObj, dateFormat} from "@/utils/libs";
import AIDrawing from "@/views/admin/AIDrawing.vue";

const activeName = ref('basic')
const system = ref({models: []})
const configBak = ref({})
const loading = ref(true)
const systemFormRef = ref(null)
const models = ref([])
const openAIModels = ref([])
const notice = ref("")
const license = ref({is_active: false})
const menus = ref([])

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
    openAIModels.value = models.value.filter(v => v.platform === "OpenAI")
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
  httpGet("/api/admin/config/get/license").then(res => {
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
        system.value['power_price'] = parseFloat(system.value['power_price']) ?? 0
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
  httpPost("/api/admin/active", {license: licenseKey.value}).then(res => {
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