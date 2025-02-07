<template>
  <div class="user-info flex-center-col" id="user-info">
    <el-form :model="user" label-width="80px" label-position="left">
      <el-row>
        <el-upload class="avatar-uploader" :auto-upload="true" :show-file-list="false" :http-request="afterRead" accept=".png,.jpg,.jpeg,.bmp">
          <el-tooltip content="点击上传头像" placement="top" v-if="user.avatar">
            <el-avatar :src="user.avatar" shape="circle" :size="100" />
          </el-tooltip>
          <el-icon v-else class="avatar-uploader-icon">
            <Plus />
          </el-icon>
        </el-upload>
      </el-row>
      <el-form-item label="昵称">
        <el-input v-model="user['nickname']" />
      </el-form-item>
      <el-form-item label="账号">
        <div class="flex">
          <span>{{ user.username }}</span>
          <el-tooltip class="box-item" content="您已经是 VIP 会员" placement="right">
            <span class="vip-icon"><el-image v-if="user.vip" :src="vipImg" class="rounded-full ml-1 size-5" /></span>
          </el-tooltip>
        </div>
      </el-form-item>
      <el-form-item label="剩余算力">
        <el-text type="warning">{{ user["power"] }}</el-text>
        <el-tag type="info" size="small" class="ml-2 cursor-pointer" @click="gotoLog">算力日志</el-tag>
        <el-tooltip :content="`每日签到可获得 ${systemConfig.daily_power} 算力`" placement="top" v-if="systemConfig.daily_power > 0">
          <el-button type="primary" size="small" @click="signIn" class="ml-2">签到</el-button>
        </el-tooltip>
      </el-form-item>
      <el-form-item label="会员到期时间" v-if="user['expired_time'] > 0">
        <el-tag type="danger">{{ dateFormat(user["expired_time"]) }}</el-tag>
      </el-form-item>

      <el-row class="opt-line">
        <el-button :dark="false" type="primary" @click="save">保存</el-button>
      </el-row>
    </el-form>
  </div>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { httpGet, httpPost } from "@/utils/http";
import { ElMessage } from "element-plus";
import { Plus } from "@element-plus/icons-vue";
import Compressor from "compressorjs";
import { dateFormat } from "@/utils/libs";
import { checkSession, getSystemInfo } from "@/store/cache";
import { useRouter } from "vue-router";
import { showMessageError, showMessageOK } from "@/utils/dialog";
const user = ref({
  vip: false,
  username: "演示数据",
  nickname: "演示数据",
  avatar: "/images/menu/member.png",
  mobile: "演示数据",
  power: 99999,
});

const vipImg = ref("/images/menu/member.png");
const systemConfig = ref({});
const router = useRouter();
const emits = defineEmits(["hide"]);
onMounted(() => {
  checkSession()
    .then(() => {
      // 获取最新用户信息
      httpGet("/api/user/profile")
        .then((res) => {
          user.value = res.data;
        })
        .catch((e) => {
          ElMessage.error("获取用户信息失败：" + e.message);
        });
    })
    .catch((e) => {
      console.log(e);
    });

  getSystemInfo().then((res) => {
    systemConfig.value = res.data;
  });
});

const afterRead = (file) => {
  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData();
      formData.append("file", result, result.name);
      // 执行上传操作
      httpPost("/api/upload", formData)
        .then((res) => {
          user.value.avatar = res.data.url;
          ElMessage.success({ message: "上传成功", duration: 500 });
        })
        .catch((e) => {
          ElMessage.error("图片上传失败:" + e.message);
        });
    },
    error(err) {
      console.log(err.message);
    },
  });
};

const save = () => {
  httpPost("/api/user/profile/update", user.value)
    .then(() => {
      ElMessage.success({ message: "更新成功", duration: 500 });
    })
    .catch((e) => {
      ElMessage.error("更新失败：" + e.message);
    });
};

const gotoLog = () => {
  router.push("/powerLog");
  emits("hide", false);
};

const signIn = () => {
  httpGet("/api/user/signin")
    .then(() => {
      showMessageOK("签到成功");
      user.value.power += systemConfig.value.daily_power;
    })
    .catch((e) => {
      showMessageError(e.message);
    });
};
</script>

<style lang="stylus" scoped>
.user-info {

  .el-row {
    justify-content center
    margin-bottom 10px
  }

  .vip-icon {
    position relative
    top 5px
  }

  .opt-line {

    .el-button {
      width 100%
    }
  }
}
</style>
