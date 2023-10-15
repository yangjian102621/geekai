<template>
  <div class="page-apps custom-scroll">
    <div class="title">
      AI 助手应用中心
    </div>
    <div class="inner" :style="{height: listBoxHeight + 'px'}">
      <ItemList :items="list" v-if="list.length > 0" gap="20" width="250">
        <template #default="scope">
          <div class="app-item" :style="{width: scope.width+'px'}">
            <el-image :src="scope.item.icon" fit="cover" :style="{height: scope.width+'px'}"/>
            <div class="title">
              <span class="name">{{ scope.item.name }}</span>
              <div class="opt">
                <el-button size="small"
                           style="--el-color-primary:#009999"
                           @click="addRole(scope.item)">
                  <el-icon>
                    <Plus/>
                  </el-icon>
                  <span>添加应用</span>
                </el-button>
              </div>
            </div>
            <div class="hello-msg">{{ scope.item['hello_msg'] }}</div>
          </div>
        </template>
      </ItemList>
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue"
import {ElMessage} from "element-plus";
import {httpGet} from "@/utils/http";
import ItemList from "@/components/ItemList.vue";
import {Plus} from "@element-plus/icons-vue";

const listBoxHeight = window.innerHeight - 97
const list = ref([])
onMounted(() => {
  httpGet("/api/role/list?all=true").then((res) => {
    const data = res.data
    for (let i = 0; i < data.length; i++) {
      if (data[i].key === 'gpt') {
        continue
      }
      list.value.push(data[i])
    }
  }).catch(e => {
    ElMessage.error("获取应用失败：" + e.message)
  })
})

const addRole = (row) => {

}
</script>

<style lang="stylus">
@import "@/assets/css/chat-app.styl"
@import "@/assets/css/custom-scroll.styl"
</style>
