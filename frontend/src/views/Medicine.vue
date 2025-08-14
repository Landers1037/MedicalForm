<template>
  <div class="medicine-container">
    <!-- 搜索和操作栏 -->
    <div class="search-section">
      <el-card class="search-card" shadow="never">
        <template #header>
          <div class="card-header">
            <div class="header-left">
              <span>药品管理</span>
            </div>
            <div class="header-right">
              <div class="toolbar-buttons">
                <el-button type="primary" plain @click="showAddDialog = true" :icon="Plus">
                  添加药品
                </el-button>
                <el-button type="success" plain @click="handleCrawlMedicine" :icon="Download" :loading="crawling">
                  {{ crawling ? '爬取中...' : '爬取药品' }}
                </el-button>
                <el-button plain @click="handleRefresh" :icon="Refresh">
                  刷新
                </el-button>
              </div>
              <div class="search-input">
                <el-input v-model="searchKeyword" placeholder="搜索药品名称或厂家" style="width: 300px" clearable
                  @input="handleSearch">
                  <template #prefix>
                    <el-icon>
                      <Search />
                    </el-icon>
                  </template>
                </el-input>
              </div>
            </div>
          </div>
        </template>
      </el-card>
    </div>

    <!-- 数据表格 -->
    <div class="table-section">
      <el-card class="table-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span>药品列表 (共 {{ filteredData.length }} 条记录)</span>
          </div>
        </template>
        <el-table :data="paginatedData" style="width: 100%" stripe v-loading="loading">
          <el-table-column prop="id" label="ID" width="80" sortable />
          <el-table-column prop="name" label="药品名称" width="200" sortable show-overflow-tooltip />
          <el-table-column prop="manufacturer" label="生产厂家" width="200" show-overflow-tooltip />
          <el-table-column prop="price" label="价格" width="120" sortable />
          <el-table-column prop="specification" label="规格" show-overflow-tooltip />
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="handleEdit(row)" :icon="Edit">
                编辑
              </el-button>
              <el-button type="danger" link @click="handleDelete(row)" :icon="Delete">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-wrapper">
          <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]"
            :total="filteredData.length" layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange" @current-change="handleCurrentChange" />
        </div>
      </el-card>
    </div>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="showAddDialog" :title="editingMedicine ? '编辑药品信息' : '添加药品信息'" width="600px"
      :close-on-click-modal="false" @close="handleDialogClose">
      <el-form ref="formRef" :model="medicineForm" :rules="formRules" label-width="100px">
        <el-form-item label="药品名称" prop="name">
          <el-input v-model="medicineForm.name" placeholder="请输入药品名称" />
        </el-form-item>
        <el-form-item label="生产厂家" prop="manufacturer">
          <el-input v-model="medicineForm.manufacturer" placeholder="请输入生产厂家" />
        </el-form-item>
        <el-form-item label="价格" prop="price">
          <el-input v-model="medicineForm.price" placeholder="请输入价格，如：10.5元" />
        </el-form-item>
        <el-form-item label="药品规格" prop="specification">
          <el-input v-model="medicineForm.specification" placeholder="请输入药品规格，如：100mg*10片" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button plain @click="showAddDialog = false">取消</el-button>
          <el-button plain type="primary" @click="handleSubmit" :loading="submitting">
            {{ editingMedicine ? '更新' : '添加' }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Refresh,
  Search,
  Edit,
  Delete,
  Download
} from '@element-plus/icons-vue'
import {
  GetAllMedicines,
  CreateMedicine,
  UpdateMedicine,
  DeleteMedicine,
  SearchMedicines,
  CrawlMedicineDataWithPages,
  GetMedicineCount
} from '../../wailsjs/go/main/App'

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const crawling = ref(false)
const showAddDialog = ref(false)
const editingMedicine = ref(null)
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const medicines = ref([])

// 表单引用
const formRef = ref()

// 药品表单
const medicineForm = ref({
  name: '',
  manufacturer: '',
  price: '',
  specification: ''
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入药品名称', trigger: 'blur' },
    { min: 1, max: 100, message: '药品名称长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  manufacturer: [
    { min: 0, max: 200, message: '生产厂家长度不能超过 200 个字符', trigger: 'blur' }
  ],
  price: [
    // 价格为非必填项
  ],
  specification: [
    { min: 0, max: 200, message: '药品规格长度不能超过 200 个字符', trigger: 'blur' }
  ]
}

// 计算属性
const filteredData = computed(() => {
  if (!medicines.value || !Array.isArray(medicines.value)) {
    return []
  }
  let data = [...medicines.value]

  // 关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    data = data.filter(item =>
      item.name.toLowerCase().includes(keyword) ||
      item.manufacturer.toLowerCase().includes(keyword)
    )
  }

  return data
})

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredData.value.slice(start, end)
})

// 方法
const loadMedicines = async () => {
  loading.value = true
  try {
    const result = await GetAllMedicines()
    medicines.value = result || []
  } catch (error) {
    ElMessage.error('获取药品数据失败：' + error.message)
    medicines.value = []
  } finally {
    loading.value = false
  }
}

const handleRefresh = async () => {
  await loadMedicines()
  ElMessage.success('数据刷新成功')
}

const handleSearch = async () => {
  currentPage.value = 1
  if (!searchKeyword.value.trim()) {
    await loadMedicines()
    return
  }
  
  try {
    loading.value = true
    const result = await SearchMedicines('name', searchKeyword.value)
    medicines.value = result
  } catch (error) {
    console.error('搜索药品失败:', error)
    ElMessage.error('搜索药品失败：' + error.message)
  } finally {
    loading.value = false
  }
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
}

const handleCurrentChange = (page) => {
  currentPage.value = page
}

const handleEdit = (row) => {
  editingMedicine.value = row
  medicineForm.value = { ...row }
  showAddDialog.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除药品 "${row.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await DeleteMedicine(row.id)
    ElMessage.success('删除成功')
    await loadMedicines()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败：' + error.message)
    }
  }
}

const handleDialogClose = () => {
  editingMedicine.value = null
  medicineForm.value = {
    name: '',
    manufacturer: '',
    price: '',
    specification: ''
  }
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    const submitData = {
      ...medicineForm.value,
      name: String(medicineForm.value.name || ''),
      manufacturer: String(medicineForm.value.manufacturer || ''),
      price: String(medicineForm.value.price || ''),
      specification: String(medicineForm.value.specification || '')
    }

    if (editingMedicine.value) {
      // 编辑模式
      submitData.id = editingMedicine.value.id
      await UpdateMedicine(submitData)
      ElMessage.success('药品信息更新成功')
    } else {
      // 添加模式
      await CreateMedicine(submitData)
      ElMessage.success('药品添加成功')
    }

    showAddDialog.value = false
    await loadMedicines()
  } catch (error) {
    ElMessage.error('操作失败：' + error.message)
  } finally {
    submitting.value = false
  }
}

// 爬取药品数据
const handleCrawlMedicine = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要从国产药品数据库爬取药品信息吗？这可能需要一些时间。',
      '确认爬取',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      }
    )

    crawling.value = true
    
    // 获取爬取前的药品数量
    const beforeCount = await GetMedicineCount()
    ElMessage.info(`开始爬取药品数据，当前数据库中有 ${beforeCount} 条药品记录`)
    
    // 爬取药品数据（默认爬取5页）
    await CrawlMedicineDataWithPages(5)
    
    // 获取爬取后的药品数量
    const afterCount = await GetMedicineCount()
    const newCount = afterCount - beforeCount
    
    ElMessage.success(`药品数据爬取完成！新增 ${newCount} 条药品记录`)
    
    // 刷新药品列表
    await loadMedicines()
    
  } catch (error) {
    if (error !== 'cancel') {
      console.error('爬取药品数据失败:', error)
      ElMessage.error('爬取药品数据失败：' + error.message)
    }
  } finally {
    crawling.value = false
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadMedicines()
})
</script>

<style scoped>
.medicine-container {
  padding: 16px;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.search-section {
  flex-shrink: 0;
}

.table-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.search-card,
.table-card {
  height: 100%;
}

.table-card {
  display: flex;
  flex-direction: column;
}

.table-card :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0;
}

.table-card .el-table {
  flex: 1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 15px;
}

.toolbar-buttons {
  display: flex;
  gap: 10px;
}

.search-input {
  display: flex;
  align-items: center;
}

.pagination-wrapper {
  padding: 20px;
  display: flex;
  justify-content: center;
  border-top: 1px solid var(--el-border-color-lighter);
  margin-top: auto;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>