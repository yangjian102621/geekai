<template>
  <el-dialog
      v-model="$props.show"
      :close-on-click-modal="false"
      :show-close="false"
      title="聊天配置"
  >
    <div class="user-info">
      <el-input v-model="user['api_key']" placeholder="填写你 OpenAI 的 API KEY">
        <template #prepend>API KEY</template>
      </el-input>

      <el-descriptions
          class="margin-top"
          title="账户信息"
          :column="2"
          border
      >

        <el-descriptions-item>
          <template #label>
            <div class="cell-item">
              <el-icon>
                <UserFilled/>
              </el-icon>
              账户
            </div>
          </template>
          {{ user.name }}
        </el-descriptions-item>

        <el-descriptions-item>
          <template #label>
            <div class="cell-item">
              <el-icon>
                <List/>
              </el-icon>
              聊天记录
            </div>
          </template>
          <el-tag v-if="user['enable_history']" type="success">已开通</el-tag>
          <el-tag v-else type="info">未开通</el-tag>
        </el-descriptions-item>

        <el-descriptions-item>
          <template #label>
            <div class="cell-item">
              <el-icon>
                <Histogram/>
              </el-icon>
              总调用次数
            </div>
          </template>
          {{ user['max_calls'] }}
        </el-descriptions-item>

        <el-descriptions-item>
          <template #label>
            <div class="cell-item">
              <el-icon>
                <Histogram/>
              </el-icon>
              剩余点数
            </div>
          </template>
          {{ user["remaining_calls"] }}
        </el-descriptions-item>

        <el-descriptions-item>
          <template #label>
            <div class="cell-item">
              <el-icon>
                <Timer/>
              </el-icon>
              激活时间
            </div>
          </template>
          {{ user['active_time'] }}
        </el-descriptions-item>

        <el-descriptions-item>
          <template #label>
            <div class="cell-item">
              <el-icon>
                <Watch/>
              </el-icon>
              到期时间
            </div>
          </template>
          {{ user['expired_time'] }}
        </el-descriptions-item>

      </el-descriptions>
    </div>
    <el-row class="row-center">
      <span>其他功能正在开发中，有什么使用建议可以通过下面的方式联系作者。</span>
    </el-row>
    <el-row>
      <el-col :span="12">
        <el-image :src="wechatGroup"></el-image>
      </el-col>

      <el-col :span="12">
        <el-image :src="wechatCard"></el-image>
      </el-col>
    </el-row>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="close">关闭</el-button>
        <el-button type="primary" @click="save">
          保存
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script>
import {defineComponent} from "vue"
import {
  List, Timer, Watch,
  UserFilled,
  Histogram
} from '@element-plus/icons-vue'
import {getLoginUser} from "@/utils/storage";
import {dateFormat} from "@/utils/libs";

export default defineComponent({
  name: 'ConfigDialog',
  components: {Watch, Timer, UserFilled, List, Histogram},
  props: {
    show: {
      type: Boolean,
      default: true
    }
  },
  data() {
    return {
      user: {},
      wechatGroup: "https://img.r9it.com/chatgpt/wechat-group.jpeg",
      wechatCard: "https://img.r9it.com/chatgpt/wechat-card.jpeg"
    }
  },
  mounted() {
    // 获取用户信息
    const data = getLoginUser();
    this.user = data.user;
    this.user['active_time'] = dateFormat(this.user['active_time']);
    this.user['expired_time'] = dateFormat(this.user['expired_time']);
  },
  methods: {
    save: function () {
      this.$emit('update:show', false);
    },
    close: function () {
      this.$emit('update:show', false);
    }
  }
})
</script>

<style lang="stylus">
.el-dialog {
  --el-dialog-width 90%;
  max-width 800px;

  .el-dialog__body {
    padding-top 10px;

    .user-info {
      .margin-top {
        margin-top 20px;
      }

      .el-icon {
        top 2px;
      }

      margin-bottom 15px;
    }
  }
}
</style>