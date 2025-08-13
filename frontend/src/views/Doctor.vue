<template>
  <div class="doctor-container">
    <!-- 页面标题和操作栏 -->
    <div class="header-section">
      <el-card class="header-card" shadow="never">
        <div class="header-content">
          <div class="title-section">
            <h2 class="page-title">
              <el-icon class="title-icon"><User /></el-icon>
              医师管理
            </h2>
            <p class="page-desc">管理医师信息，包括姓名、性别、级别等</p>
          </div>
          <div class="action-section">
            <el-button type="primary" plain :icon="Plus" @click="handleAdd">
              添加医师
            </el-button>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 搜索和筛选 -->
    <div class="search-section">
      <el-card class="search-card" shadow="never">
        <div class="search-content">
          <el-row :gutter="20">
            <el-col :span="8">
              <el-input
                v-model="searchKeyword"
                placeholder="请输入搜索关键词"
                :prefix-icon="Search"
                clearable
                @input="handleSearch"
              />
            </el-col>
            <el-col :span="6">
              <el-select
                v-model="searchType"
                placeholder="搜索类型"
                @change="handleSearch"
              >
                <el-option label="全部" value="all" />
                <el-option label="姓名" value="name" />
                <el-option label="级别" value="level" />
                <el-option label="电话" value="phone" />
              </el-select>
            </el-col>
            <el-col :span="6">
              <el-select
                v-model="levelFilter"
                placeholder="级别筛选"
                clearable
                @change="handleSearch"
              >
                <el-option label="全部级别" value="" />
                <el-option label="主治医师" value="主治医师" />
                <el-option label="助理医师" value="助理医师" />
                <el-option label="主任" value="主任" />
                <el-option label="副主任" value="副主任" />
              </el-select>
            </el-col>
            <el-col :span="4">
              <el-button type="primary" plain :icon="Refresh" @click="handleRefresh">
                刷新
              </el-button>
            </el-col>
          </el-row>
        </div>
      </el-card>
    </div>

    <!-- 医师列表 -->
    <div class="table-section">
      <el-card class="table-card" shadow="never">
        <el-table
          :data="paginatedData"
          style="width: 100%"
          stripe
          v-loading="loading"
          empty-text="暂无医师数据"
        >
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="name" label="姓名" width="120" />
          <el-table-column prop="sex" label="性别" width="80">
            <template #default="{ row }">
              <el-tag :type="row.sex === '男' ? 'primary' : 'danger'" size="small">
                {{ row.sex }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="level" label="级别" width="120">
            <template #default="{ row }">
              <el-tag type="success" size="small">
                {{ row.level }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="phone" label="电话" width="140" />
          <el-table-column prop="email" label="邮箱" width="180" show-overflow-tooltip />
          <el-table-column prop="remark" label="备注" show-overflow-tooltip />
          <el-table-column prop="created_at" label="创建时间" width="160">
            <template #default="{ row }">
              {{ formatDateTime(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="160" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link :icon="Edit" @click="handleEdit(row)">
                编辑
              </el-button>
              <el-button type="danger" link :icon="Delete" @click="handleDelete(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="filteredData.length"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </el-card>
    </div>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="showAddDialog"
      :title="editingDoctor ? '编辑医师信息' : '添加医师信息'"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form ref="formRef" :model="doctorForm" :rules="formRules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="doctorForm.name" placeholder="请输入姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="性别" prop="sex">
              <el-radio-group v-model="doctorForm.sex">
                <el-radio label="男">男</el-radio>
                <el-radio label="女">女</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="级别" prop="level">
              <el-select
                v-model="doctorForm.level"
                placeholder="请选择级别"
                filterable
                allow-create
                default-first-option
              >
                <el-option label="主治医师" value="主治医师" />
                <el-option label="助理医师" value="助理医师" />
                <el-option label="主任" value="主任" />
                <el-option label="副主任" value="副主任" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="电话" prop="phone">
              <el-input v-model="doctorForm.phone" placeholder="请输入电话" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="24">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="doctorForm.email" placeholder="请输入邮箱" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="24">
            <el-form-item label="备注">
              <el-input
                v-model="doctorForm.remark"
                type="textarea"
                :rows="3"
                placeholder="请输入备注信息"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddDialog = false" plain>取消</el-button>
          <el-button type="primary" plain @click="handleSubmit" :loading="submitting">
            {{ editingDoctor ? '更新' : '添加' }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  User,
  Plus,
  Search,
  Refresh,
  Edit,
  Delete
} from '@element-plus/icons-vue'
import { useDoctorStore } from '@/stores/doctor'

const doctorStore = useDoctorStore()

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const showAddDialog = ref(false)
const editingDoctor = ref(null)
const formRef = ref(null)
const isUnmounted = ref(false)

// 搜索相关
const searchKeyword = ref('')
const searchType = ref('all')
const levelFilter = ref('')

// 分页相关
const currentPage = ref(1)
const pageSize = ref(20)

// 表单数据
const doctorForm = reactive({
  name: '',
  sex: '男',
  level: '',
  phone: '',
  email: '',
  remark: ''
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '姓名长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  sex: [
    { required: true, message: '请选择性别', trigger: 'change' }
  ],
  level: [
    { required: true, message: '请选择或输入级别', trigger: 'change' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

// 计算属性
const filteredData = computed(() => {
  let data = doctorStore.doctors
  
  // 级别筛选
  if (levelFilter.value) {
    data = data.filter(doctor => doctor.level === levelFilter.value)
  }
  
  // 关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    data = data.filter(doctor => {
      switch (searchType.value) {
        case 'name':
          return doctor.name.toLowerCase().includes(keyword)
        case 'level':
          return doctor.level.toLowerCase().includes(keyword)
        case 'phone':
          return doctor.phone.includes(keyword)
        default:
          return (
            doctor.name.toLowerCase().includes(keyword) ||
            doctor.level.toLowerCase().includes(keyword) ||
            doctor.phone.includes(keyword) ||
            (doctor.email && doctor.email.toLowerCase().includes(keyword))
          )
      }
    })
  }
  
  return data
})

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredData.value.slice(start, end)
})

// 方法
const formatDateTime = (dateTime) => {
  if (!dateTime) return '-'
  return new Date(dateTime).toLocaleString('zh-CN')
}

const handleAdd = () => {
  editingDoctor.value = null
  resetForm()
  showAddDialog.value = true
}

const handleEdit = (doctor) => {
  editingDoctor.value = doctor
  Object.assign(doctorForm, {
    name: doctor.name,
    sex: doctor.sex,
    level: doctor.level,
    phone: doctor.phone,
    email: doctor.email || '',
    remark: doctor.remark || ''
  })
  showAddDialog.value = true
}

const handleDelete = async (doctor) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除医师 "${doctor.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await doctorStore.deleteDoctor(doctor.id)
    if (!isUnmounted.value) {
      ElMessage.success('删除成功')
    }
  } catch (error) {
    if (error !== 'cancel' && !isUnmounted.value) {
      ElMessage.error('删除失败：' + error.message)
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    if (editingDoctor.value) {
      // 更新医师
      const updateData = {
        id: editingDoctor.value.id,
        ...doctorForm
      }
      await doctorStore.updateDoctor(updateData)
      if (!isUnmounted.value) {
        ElMessage.success('更新成功')
      }
    } else {
      // 添加医师
      await doctorStore.addDoctor(doctorForm)
      if (!isUnmounted.value) {
        ElMessage.success('添加成功')
      }
    }
    
    showAddDialog.value = false
  } catch (error) {
    if (!isUnmounted.value) {
      ElMessage.error((editingDoctor.value ? '更新' : '添加') + '失败：' + error.message)
    }
  } finally {
    submitting.value = false
  }
}

const handleDialogClose = () => {
  resetForm()
  editingDoctor.value = null
}

const resetForm = () => {
  Object.assign(doctorForm, {
    name: '',
    sex: '男',
    level: '',
    phone: '',
    email: '',
    remark: ''
  })
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

const handleSearch = () => {
  currentPage.value = 1
}

const handleRefresh = async () => {
  loading.value = true
  try {
    await doctorStore.fetchAllDoctors()
    if (!isUnmounted.value) {
      ElMessage.success('数据刷新成功')
    }
  } catch (error) {
    if (!isUnmounted.value) {
      ElMessage.error('刷新失败：' + error.message)
    }
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

// 初始化
onMounted(async () => {
  loading.value = true
  try {
    await doctorStore.fetchAllDoctors()
  } catch (error) {
    if (!isUnmounted.value) {
      console.error('加载数据失败：', error)
      ElMessage({
        message: '加载数据失败：' + error.message,
        type: 'error',
        duration: 1500,
      })
    }
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  isUnmounted.value = true
})
</script>

<style scoped>
.doctor-container {
  padding: 20px;
  height: calc(100vh - 60px);
  overflow: auto;
}

.header-section,
.search-section,
.table-section {
  margin-bottom: 20px;
}

.header-card,
.search-card,
.table-card {
  border-radius: var(--el-border-radius-base);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-section {
  flex: 1;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 24px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 8px 0;
}

.title-icon {
  color: var(--el-color-primary);
}

.page-desc {
  color: var(--el-text-color-regular);
  margin: 0;
  font-size: 14px;
}

.action-section {
  flex-shrink: 0;
}

.search-content {
  padding: 0;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  padding: 20px 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .doctor-container {
    padding: 12px;
  }
  
  .header-content {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }
  
  .action-section {
    display: flex;
    justify-content: center;
  }
}
</style>