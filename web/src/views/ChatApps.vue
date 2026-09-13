<template>
  <div>
    <div class="page-apps custom-scroll">
      <div class="flex h-20">
        <ul class="scrollbar-type-nav">
          <li :class="{ active: typeId === '' }" @click="getAppList('')">全部分类</li>
          <li
            v-for="item in appTypes"
            :key="item.id"
            :class="{ active: typeId === item.id }"
            @click="getAppList(item.id)"
          >
            <div class="image" v-if="item.icon">
              <el-image :src="item.icon" fit="cover" />
            </div>
            {{ item.name }}
          </li>
        </ul>
      </div>

      <div class="app-list-container" :style="{ height: listBoxHeight + 'px' }">
        <ItemList :items="list" v-if="list.length > 0" :gap="15" :width="300">
          <template #default="scope">
            <div class="item">
              <div class="image">
                <el-image :src="scope.item.icon" fit="cover" />
              </div>

              <div class="inner">
                <div class="info">
                  <div class="info-title">{{ scope.item.name }}</div>
                  <div class="info-text">{{ scope.item.hello_msg }}</div>
                </div>
                <div class="btn">
                  <el-button size="small" class="sm-btn-theme" @click="useRole(scope.item)"
                    >使用</el-button
                  >
                </div>
              </div>
            </div>
          </template>
        </ItemList>
        <div v-else style="width: 100%">
          <el-empty description="暂无数据" :image="nodata" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import nodata from '@/assets/img/no-data.png'
import ItemList from '@/components/ItemList.vue'
import { httpGet } from '@/utils/http'
import { substr } from '@/utils/libs'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const listBoxHeight = window.innerHeight - 133

const typeId = ref('')
const appTypes = ref([])
const list = ref([])

onMounted(() => {
  getAppType()
  getAppList()
})

const getAppType = () => {
  httpGet('/api/app/type/list')
    .then((res) => {
      appTypes.value = res.data
    })
    .catch((e) => {
      ElMessage.error('获取分类失败：' + e.message)
    })
}

const getAppList = (tid = '') => {
  typeId.value = tid
  httpGet('/api/app/list', { tid })
    .then((res) => {
      const items = res.data
      // 处理 hello message
      for (let i = 0; i < items.length; i++) {
        items[i].intro = substr(items[i].hello_msg, 80)
      }
      list.value = items
    })
    .catch((e) => {
      ElMessage.error('获取应用失败：' + e.message)
    })
}

const router = useRouter()
const useRole = (role) => {
  router.push(`/chat?role_id=${role.id}`)
}
</script>

<style lang="scss" scoped>
@use '../assets/css/chat-app.scss' as *;
@use '../assets/css/custom-scroll.scss' as *;
</style>
