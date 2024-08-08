<template>
  <div>
    <div class="page-apps custom-scroll">
      <div class="inner" :style="{height: listBoxHeight + 'px'}">
        <ItemList :items="list" v-if="list.length > 0" :gap="15" :width="300">
          <template #default="scope">
            <div class="item">
              <div class="image">
                <el-image :src="scope.item.icon" fit="cover"/>
              </div>

              <div class="inner">
                <div class="info">
                  <div class="info-title">{{ scope.item.name }}</div>
                  <div class="info-text">{{ scope.item.hello_msg }}</div>
                </div>
                <div class="btn">
                  <el-button size="small" color="#21aa93" @click="useRole(scope.item)">使用</el-button>
                  <el-tooltip effect="light" content="从工作区移除" placement="top" v-if="hasRole(scope.item.key)">
                    <el-button size="small" type="danger" @click="updateRole(scope.item,'remove')">移除</el-button>
                  </el-tooltip>
                  <el-tooltip effect="light" content="添加到工作区" placement="top" v-else>
                    <el-button size="small" style="--el-color-primary:#009999" @click="updateRole(scope.item, 'add')">添加</el-button>
                  </el-tooltip>
                </div>
              </div>


            </div>
            <!--            <div class="app-item">-->
            <!--              <el-image :src="scope.item.icon" fit="cover"/>-->
            <!--              <div class="title">-->
            <!--                <span class="name">{{ scope.item.name }}</span>-->
            <!--                <div class="opt">-->
            <!--                  <div v-if="hasRole(scope.item.key)">-->
            <!--                    <el-button size="small" type="success" @click="useRole(scope.item)">使用</el-button>-->
            <!--                    <el-button size="small" type="danger" @click="updateRole(scope.item,'remove')">移除</el-button>-->
            <!--                  </div>-->
            <!--                  <el-button v-else size="small"-->
            <!--                             style="&#45;&#45;el-color-primary:#009999"-->
            <!--                             @click="updateRole(scope.item, 'add')">-->
            <!--                    <el-icon>-->
            <!--                      <Plus/>-->
            <!--                    </el-icon>-->
            <!--                    <span>添加应用</span>-->
            <!--                  </el-button>-->
            <!--                </div>-->
            <!--              </div>-->
            <!--              <div class="hello-msg" ref="elements">{{ scope.item.intro }}</div>-->
            <!--            </div>-->
          </template>
        </ItemList>
      </div>
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue"
import {ElMessage} from "element-plus";
import {httpGet, httpPost} from "@/utils/http";
import {checkSession} from "@/store/cache";
import {arrayContains, removeArrayItem, substr} from "@/utils/libs";
import {useRouter} from "vue-router";
import {useSharedStore} from "@/store/sharedata";
import ItemList from "@/components/ItemList.vue";
import {Plus} from "@element-plus/icons-vue";

const listBoxHeight = window.innerHeight - 87
const list = ref([])
const roles = ref([])
const store = useSharedStore();

onMounted(() => {
  httpGet("/api/role/list").then((res) => {
    const items = res.data
    // 处理 hello message
    for (let i = 0; i < items.length; i++) {
      items[i].intro = substr(items[i].hello_msg, 80)
    }
    list.value = items
  }).catch(e => {
    ElMessage.error("获取应用失败：" + e.message)
  })

  getRoles()
})

const getRoles = () => {
  checkSession().then(user => {
    roles.value = user.chat_roles
  }).catch(() => {
  })
}

const updateRole = (row, opt) => {
  checkSession().then(() => {
    const title = ref("")
    if (opt === "add") {
      title.value = "添加应用"
      const exists = arrayContains(roles.value, row.key)
      if (exists) {
        return
      }
      roles.value.push(row.key)
    } else {
      title.value = "移除应用"
      const exists = arrayContains(roles.value, row.key)
      if (!exists) {
        return
      }
      roles.value = removeArrayItem(roles.value, row.key)
    }
    httpPost("/api/role/update", {keys: roles.value}).then(() => {
      ElMessage.success({message: title.value + "成功！", duration: 1000})
    }).catch(e => {
      ElMessage.error(title.value + "失败：" + e.message)
    })
  }).catch(() => {
    store.setShowLoginDialog(true)
  })
}

const hasRole = (roleKey) => {
  return arrayContains(roles.value, roleKey, (v1, v2) => v1 === v2)
}

const router = useRouter()
const useRole = (role) => {
  router.push(`/chat?role_id=${role.id}`)
}
</script>

<style lang="stylus">
@import "@/assets/css/chat-app.styl"
@import "@/assets/css/custom-scroll.styl"
</style>
