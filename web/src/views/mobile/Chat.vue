<template>
  <div class="chat-mobile">
    <van-nav-bar :title="title"/>

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
  </div>
</template>

<script setup>
import {ref} from "vue";
import {httpGet} from "@/utils/http";
import {getLoginUser} from "@/utils/storage";
import {showFailToast} from "vant";

const title = ref("会话列表")
const chatName = ref("")
const chats = ref([])
const allChats = ref([])
const loading = ref(false)
const finished = ref(false)
const error = ref(false)
const user = getLoginUser()

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
  const roles = [];
  for (let i = 0; i < allChats.value.length; i++) {
    if (allChats.value[i].title.toLowerCase().indexOf(chatName.value.toLowerCase()) !== -1) {
      roles.push(allChats.value[i]);
    }
  }
  chats.value = roles;
}

</script>

<style scoped lang="stylus">

.chat-mobile {
  .content {
    padding: 0 10px;

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
}
</style>