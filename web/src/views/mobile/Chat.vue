<template>
  <div class="container chat-mobile">
    <van-nav-bar
        :title="title"
        right-text="新建会话"
        @click-right="showPicker = true"
    />

    <div class="content">
      <van-search
          v-model="chatName"
          placeholder="请输入会话标题"
          input-align="center"
          @input="search"
      />

      <van-list
          v-model:loading="loading"
          v-model:error="error"
          error-text="请求失败，点击重新加载"
          :finished="finished"
          finished-text="没有更多了"
          @load="onLoad"
      >
        <van-cell v-for="item in chats" :key="item.id">
          <div class="chat-list-item">
            <van-image
                round
                :src="item.icon"
            />
            <van-text-ellipsis class="text" :content="item.title"/>
          </div>
        </van-cell>
      </van-list>
    </div>

    <van-popup v-model:show="showPicker" position="bottom">
      <van-picker
          title="选择模型和角色"
          :columns="columns"
          @cancel="showPicker = false"
          @confirm="newChat"
      >
        <template #option="item">
          <div class="picker-option">
            <van-image
                fit="cover"
                :src="item.icon"
                round
                v-if="item.icon"
            />
            <span>{{ item.text }}</span>
          </div>
        </template>
      </van-picker>
    </van-popup>
  </div>
</template>

<script setup>
import {ref} from "vue";
import {httpGet} from "@/utils/http";
import {getLoginUser} from "@/utils/storage";
import {showFailToast} from "vant";
import {checkSession} from "@/action/session";
import router from "@/router";

const title = ref("会话列表")
const chatName = ref("")
const chats = ref([])
const allChats = ref([])
const loading = ref(false)
const finished = ref(false)
const error = ref(false)
const user = getLoginUser()

const showPicker = ref(false)
const columns = ref([])

checkSession().then(() => {
  // 加载角色列表
  httpGet(`/api/role/list?user_id=${user.id}`).then((res) => {
    if (res.data) {
      const items = res.data
      const roles = []
      for (let i = 0; i < items.length; i++) {
        // console.log(items[i])
        roles.push({text: items[i].name, value: items[i].id, icon: items[i].icon})
      }
      columns.value[0] = roles
    }
  }).catch(() => {
    showFailToast("加载聊天角色失败")
  })

  // 加载系统配置
  httpGet('/api/admin/config/get?key=system').then(res => {
    if (res.data) {
      const items = res.data.models
      const models = []
      for (let i = 0; i < items.length; i++) {
        models.push({text: items[i].toUpperCase(), value: items[i]})
      }
      columns.value[1] = models
    }
  }).catch(() => {
    showFailToast("加载系统配置失败")
  })
}).catch(() => {
  router.push("/login")
})

const onLoad = () => {
  httpGet("/api/chat/list?user_id=" + user.id).then((res) => {
    if (res.data) {
      chats.value = res.data;
      allChats.value = res.data;
      finished.value = true
    }
    loading.value = false;
  }).catch(() => {
    error.value = true
    showFailToast("加载会话列表失败")
  })
};

const search = () => {
  if (chatName.value === '') {
    chats.value = allChats.value
    return
  }
  const items = [];
  for (let i = 0; i < allChats.value.length; i++) {
    if (allChats.value[i].title.toLowerCase().indexOf(chatName.value.toLowerCase()) !== -1) {
      items.push(allChats.value[i]);
    }
  }
  chats.value = items;
}

const newChat = (item) => {
  console.log(item.selectedValues)
  showPicker.value = false
}

</script>

<style scoped lang="stylus">

.chat-mobile {
  .content {
    .van-cell__value {
      .chat-list-item {
        display flex

        .van-image {
          width 30px
          height 30px
        }

        .text {
          margin-top 4px;
          margin-left 10px;
        }
      }
    }
  }

  .van-picker-column {
    .picker-option {
      display flex
      width 100%
      padding 0 10px

      .van-image {
        width 20px;
        height 20px;
        margin-right 5px
      }
    }
  }
}
</style>