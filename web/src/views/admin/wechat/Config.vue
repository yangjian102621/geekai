<template>
  <div class="system-config form" v-loading="loading">
    <el-tabs v-model="activeName" class="sys-tabs">
      <el-tab-pane label="公众号配置" name="basic">
        <div class="container">
          <el-form :model="wechat" label-width="150px" label-position="right" ref="wechatFormRef" :rules="rules">
            <el-tabs type="border-card">
              <el-tab-pane label="微信公众号配置">
                <el-form-item label="微信公众号app_id" prop="wechat_app_id">
                  <el-input v-model="wechat['wechat_app_id']" />
                </el-form-item>
                <el-form-item label="微信公众号secret" prop="wechat_secret">
                  <el-input v-model="wechat['wechat_secret']" />
                </el-form-item>
                <el-form-item label="微信公众号token" prop="wechat_token">
                  <el-input v-model="wechat['wechat_token']" />
                </el-form-item>
                <el-form-item label="微信公众号aes_key" prop="wechat_aes_key">
                  <el-input v-model="wechat['wechat_aes_key']" />
                </el-form-item>
                <el-form-item label="登录回调" prop="wechat_callback">
                  <el-input v-model="wechat['wechat_callback']" />
                </el-form-item>
              </el-tab-pane>
            </el-tabs>

            <div style="padding: 10px">
              <el-form-item>
                <el-button type="primary" @click="save('wechat')">保存</el-button>
              </el-form-item>
            </div>
          </el-form>
        </div>
      </el-tab-pane>

    </el-tabs>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from "vue";
import { httpGet, httpPost } from "@/utils/http";
import Compressor from "compressorjs";
import { ElMessage, ElMessageBox } from "element-plus";
import { CloseBold, InfoFilled, Select, UploadFilled } from "@element-plus/icons-vue";
import MdEditor from "md-editor-v3";
import "md-editor-v3/lib/style.css";
import Menu from "@/views/admin/Menu.vue";
import { copyObj, dateFormat } from "@/utils/libs";
import ItemsInput from "@/components/ui/ItemsInput.vue";
import { useSharedStore } from "@/store/sharedata";

const activeName = ref("basic");
const wechat = ref({ models: [] });
const wechatFormRef = ref(null);

const store = useSharedStore();

onMounted(() => {
  // 加载系统配置
  httpGet("/api/admin/config/get?key=wechat")
      .then((res) => {
        wechat.value = res.data;
      })
      .catch((e) => {
        ElMessage.error("加载系统配置失败: " + e.message);
      });

});


const rules = reactive({
  wechat_app_id: [{ required: true, message: "请输入微信公众号appid", trigger: "blur" }],
  wechat_secret: [{ required: true, message: "请输入微信公众号secret", trigger: "blur" }],
});
const save = function (key) {

  wechatFormRef.value.validate((valid) => {
    if (valid) {
      httpPost("/api/admin/config/update", { key: key, config: wechat.value})
          .then(() => {
            ElMessage.success("操作成功！");
          })
          .catch((e) => {
            ElMessage.error("操作失败：" + e.message);
          });
    }
  });

};

</script>

<style lang="stylus" scoped>
@import "@/assets/css/admin/form.styl"
@import "@/assets/css/main.styl"
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
