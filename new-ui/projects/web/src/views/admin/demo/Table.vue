<template>
  <div>
    <div class="container">
      <div class="handle-box">
        <el-select v-model="query.address" placeholder="模型" class="handle-select mr10">
          <el-option key="1" label="GPT-3.5" value="GPT-3.5"></el-option>
          <el-option key="2" label="GPT-4.0" value="GPT-4.0"></el-option>
          <el-option key="2" label="GPT-5.0" value="GPT-5.0"></el-option>
        </el-select>
        <el-input v-model="query.name" placeholder="姓名" class="handle-input mr10"></el-input>
        <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
        <el-button type="primary" :icon="Plus">新增</el-button>
      </div>
      <el-table :data="tableData" border class="table" style="width: 100%" header-cell-class-name="table-header">
        <el-table-column prop="id" fixed label="ID" width="55" align="center"></el-table-column>
        <el-table-column prop="name" label="姓名"></el-table-column>
        <el-table-column label="头像(查看大图)" align="center">
          <template #default="scope">
            <el-image
                class="table-td-thumb"
                :src="scope.row.thumb"
                :z-index="10"
                :preview-src-list="[scope.row.thumb]"
                preview-teleported
            >
            </el-image>
          </template>
        </el-table-column>
        <el-table-column prop="info" label="简介"></el-table-column>
        <el-table-column label="状态" align="center">
          <template #default="scope">
            <el-tag
                :type="scope.row.state === '启用' ? 'success' : scope.row.state === '禁用' ? 'danger' : ''"
            >
              {{ scope.row.state }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="date" label="注册时间"></el-table-column>
        <el-table-column label="操作" fixed="right" align="center">
          <template #default="scope">
            <el-button text :icon="Edit" @click="handleEdit(scope.$index, scope.row)">
              编辑
            </el-button>
            <el-button text :icon="Delete" class="red" @click="handleDelete(scope.$index)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
            background
            layout="total, prev, pager, next"
            :current-page="query.pageIndex"
            :page-size="query.pageSize"
            :total="pageTotal"
            @current-change="handlePageChange"
        ></el-pagination>
      </div>
    </div>

    <!-- 编辑弹出框 -->
    <el-dialog title="编辑" v-model="editVisible" width="30%">
      <el-form label-width="70px">
        <el-form-item label="姓名">
          <el-input v-model="form.name"></el-input>
        </el-form-item>
        <el-form-item label="简介">
          <el-input v-model="form.info"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
				<span class="dialog-footer">
					<el-button @click="editVisible = false">取 消</el-button>
					<el-button type="primary" @click="saveEdit">确 定</el-button>
				</span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {reactive, ref} from 'vue';
import {ElMessage, ElMessageBox} from 'element-plus';
import {Delete, Edit, Plus, Search} from '@element-plus/icons-vue';

const query = reactive({
  name: '',
  pageIndex: 1,
  pageSize: 10
});
const tableData = ref([]);
const pageTotal = ref(0);
// 获取表格数据
const getData = () => {
  tableData.value = [{
    "id": 1,
    "name": "孔子",
    "info": "有朋自远方来，不亦说乎？",
    "state": "禁用",
    "date": "2023-06-21",
    "thumb": "/images/avatar/kong_zi.jpg"
  },
    {
      "id": 2,
      "name": "乔布斯",
      "info": "活着就是为了改变世界！难道还有其他原因吗？",
      "state": "禁用",
      "date": "2023-06-21",
      "thumb": "/images/avatar/steve_jobs.jpg"
    },
    {
      "id": 3,
      "name": "马斯克",
      "info": "梦想要远大，如果你的梦想没有吓到你，说明你做得不对。",
      "state": "启用",
      "date": "2023-06-21",
      "thumb": "/images/avatar/elon_musk.jpg"
    },
    {
      "id": 4,
      "name": "鲁迅",
      "info": "自由之歌，永不过时，横眉冷对千夫指，俯首甘为孺子牛。",
      "state": "启用",
      "date": "2023-06-21",
      "thumb": "/images/avatar/lu_xun.jpg"
    }
  ]
  pageTotal.value = 5
};
getData();

// 查询操作
const handleSearch = () => {
  query.pageIndex = 1;
  getData();
};
// 分页导航
const handlePageChange = (val) => {
  query.pageIndex = val;
  getData();
};

// 删除操作
const handleDelete = (index) => {
  // 二次确认删除
  ElMessageBox.confirm('确定要删除吗？', '提示', {
    type: 'warning'
  })
      .then(() => {
        ElMessage.success('删除成功');
        tableData.value.splice(index, 1);
      })
      .catch(() => {
      });
};

// 表格编辑时弹窗和保存
const editVisible = ref(false);
let form = reactive({
  name: '',
  info: ''
});
let idx = -1;
const handleEdit = (index, row) => {
  idx = index;
  form.name = row.name;
  form.info = row.info;
  editVisible.value = true;
};
const saveEdit = () => {
  editVisible.value = false;
  ElMessage.success(`修改第 ${idx + 1} 行成功`);
  tableData.value[idx].name = form.name;
  tableData.value[idx].info = form.info;
};
</script>

<style scoped>
.handle-box {
  margin-bottom: 20px;
}

.handle-select {
  width: 120px;
}

.handle-input {
  width: 300px;
}

.table {
  width: 100%;
  font-size: 14px;
}

.red {
  color: #F56C6C;
}

.mr10 {
  margin-right: 10px;
}

.table-td-thumb {
  display: block;
  margin: auto;
  width: 40px;
  height: 40px;
}
</style>
