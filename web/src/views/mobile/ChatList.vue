<template>
  <div class="container mobile-chat-list" v-if="isLogin">
    <van-nav-bar
        :title="title"
        left-text="新建会话"
        @click-left="showPicker = true">
      <template #right>
        <van-icon name="delete-o" @click="clearAllChatHistory"/>
      </template>
    </van-nav-bar>

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
        <van-swipe-cell v-for="item in chats" :key="item.id">
          <van-cell @click="changeChat(item)">
            <div class="chat-list-item">
              <van-image
                  round
                  :src="item.icon"
              />
              <div class="van-ellipsis">{{ item.title }}</div>
            </div>
          </van-cell>
          <template #right>
            <van-button square type="primary" text="修改" @click="editChat(item)"/>
            <van-button square type="danger" text="删除" @click="removeChat(item)"/>
          </template>
        </van-swipe-cell>
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
import {getLoginUser} from "@/store/session";
import {showConfirmDialog, showFailToast, showSuccessToast, showToast} from "vant";
import {checkSession} from "@/action/session";
import router from "@/router";
import {setChatConfig} from "@/store/chat";
import {removeArrayItem, UUID} from "@/utils/libs";
import {ElMessage} from "element-plus";

const title = ref("会话列表")
const chatName = ref("")
const chats = ref([])
const allChats = ref([])
const loading = ref(false)
const finished = ref(false)
const error = ref(false)
const user = getLoginUser()
const isLogin = ref(false)
const roles = ref([])
const models = ref([])
const showPicker = ref(false)
const columns = ref([roles.value, models.value])

checkSession().then(() => {
  isLogin.value = true
  // 加载角色列表
  httpGet(`/api/role/list?user_id=${user.id}`).then((res) => {
    if (res.data) {
      const items = res.data
      for (let i = 0; i < items.length; i++) {
        // console.log(items[i])
        roles.value.push({text: items[i].name, value: items[i].id, icon: items[i].icon})
      }
    }
  }).catch(() => {
    showFailToast("加载聊天角色失败")
  })

  // 加载系统配置
  httpGet('/api/admin/config/get?key=system').then(res => {
    if (res.data) {
      const items = res.data.models
      for (let i = 0; i < items.length; i++) {
        models.value.push({text: items[i].toUpperCase(), value: items[i]})
      }
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

const clearAllChatHistory = () => {
  showConfirmDialog({
    title: '操作提示',
    message: '确定要删除所有的会话记录吗？'
  }).then(() => {
    httpGet("/api/chat/clear").then(() => {
      showSuccessToast('所有聊天记录已清空')
      chats.value = [];
    }).catch(e => {
      showFailToast("操作失败：" + e.message)
    })
  }).catch(() => {
    // on cancel
  })
}

const newChat = (item) => {
  showPicker.value = false
  const options = item.selectedOptions
  setChatConfig({
    role: {
      id: options[0].value,
      name: options[0].text,
      icon: options[0].icon
    },
    model: options[1].value,
    title: '新建会话',
    chatId: UUID()
  })
  router.push('/mobile/chat/session')
}

const changeChat = (chat) => {
  let role = {}
  for (let i = 0; i < roles.value.length; i++) {
    if (roles.value[i].value === chat.role_id) {
      role = roles.value[i]
      break
    }
  }
  setChatConfig({
    role: {
      id: chat.role_id,
      name: role.text,
      icon: role.icon
    },
    model: chat.model,
    title: chat.title,
    chatId: chat.chat_id
  })
  router.push('/mobile/chat/session')
}

const editChat = (item) => {
  showToast('修改会话标题')

}

const removeChat = (item) => {
  httpGet('/api/chat/remove?chat_id=' + item.chat_id).then(() => {
    chats.value = removeArrayItem(chats.value, item, function (e1, e2) {
      return e1.id === e2.id
    })
  }).catch(e => {
    showFailToast('操作失败：' + e.message);
  })

}

</script>

<style scoped lang="stylus">

.mobile-chat-list {
  .content {
    .van-cell__value {
      .chat-list-item {
        display flex

        .van-image {
          width 30px
          height 30px
        }

        .van-ellipsis {
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

  .van-nav-bar {
    .van-nav-bar__right {
      .van-icon {
        font-size 20px;
      }
    }
  }
}
</style>