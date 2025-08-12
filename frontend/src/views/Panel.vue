<template>
  <div class="panel-container">
    <!-- 操作栏 -->
    <div class="toolbar">
      <el-card class="toolbar-card" shadow="never">
        <div class="toolbar-content">
          <div class="toolbar-left">
            <el-button type="primary" @click="showAddDialog = true" :icon="Plus">
              添加就诊记录
            </el-button>
            <el-button @click="handleRefresh" :icon="Refresh">
              刷新
            </el-button>
            <el-button type="warning" @click="handleBackup" :icon="Download">
              备份数据
            </el-button>
          </div>
          <div class="toolbar-right">
            <el-input v-model="searchForm.keyword" placeholder="搜索姓名、电话或诊断" style="width: 300px" clearable
              @input="handleSearch">
              <template #prefix>
                <el-icon>
                  <Search />
                </el-icon>
              </template>
            </el-input>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 高级搜索 -->
    <div class="search-section">
      <el-card class="search-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span>高级搜索</span>
            <el-button type="primary" link @click="searchExpanded = !searchExpanded">
              {{ searchExpanded ? '收起' : '展开' }}
            </el-button>
          </div>
        </template>
        <el-collapse-transition>
          <div v-show="searchExpanded" class="search-content">
            <el-form :model="searchForm" inline>
              <el-form-item label="性别">
                <el-select v-model="searchForm.sex" placeholder="选择性别" clearable style="width: 120px">
                  <el-option label="男" value="男" />
                  <el-option label="女" value="女" />
                </el-select>
              </el-form-item>
              <el-form-item label="日期范围">
                <el-date-picker v-model="searchForm.dateRange" type="daterange" range-separator="至"
                  start-placeholder="开始日期" end-placeholder="结束日期" format="YYYY-MM-DD" value-format="YYYY-MM-DD"
                  @change="handleSearch" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleSearch" :icon="Search">
                  搜索
                </el-button>
                <el-button @click="handleResetSearch" :icon="RefreshLeft">
                  重置
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-collapse-transition>
      </el-card>
    </div>

    <!-- 数据表格 -->
    <div class="table-section">
      <el-card class="table-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span>患者列表 (共 {{ filteredData.length }} 条记录)</span>
          </div>
        </template>
        <el-table :data="paginatedData" style="width: 100%" stripe v-loading="loading"
          @selection-change="handleSelectionChange">
          <el-table-column type="selection" width="55" />
          <el-table-column prop="id" label="ID" width="80" sortable>
            <template #default="{ row }">
              <el-button type="primary" link @click="handleViewDetail(row)">
                {{ row.id }}
              </el-button>
            </template>
          </el-table-column>
          <el-table-column prop="name" label="姓名" width="120" sortable />
          <el-table-column prop="sex" label="性别" width="80">
            <template #default="{ row }">
              <el-tag :type="row.sex === '男' ? 'primary' : 'danger'" size="small">
                {{ row.sex }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="age" label="年龄" width="80" sortable />
          <el-table-column prop="phone" label="电话" width="140" />
          <el-table-column prop="contact" label="联系方式" width="150" show-overflow-tooltip />
          <el-table-column prop="date" label="登记日期" width="120" sortable />
          <el-table-column prop="doc" label="诊断医师" width="120" show-overflow-tooltip />
          <el-table-column prop="detail" label="诊断" show-overflow-tooltip />
          <el-table-column prop="solution" label="治疗方案" show-overflow-tooltip />
          <el-table-column label="操作" width="150" fixed="right">
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
    <el-dialog v-model="showAddDialog" :title="editingPatient ? '编辑患者就诊记录' : '添加患者就诊记录'" width="680px"
      @close="handleDialogClose">
      <el-form ref="formRef" :model="patientForm" :rules="formRules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="patientForm.name" placeholder="请输入姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="性别" prop="sex">
              <el-radio-group v-model="patientForm.sex">
                <el-radio label="男">男</el-radio>
                <el-radio label="女">女</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="年龄" prop="age">
              <el-input-number v-model="patientForm.age" :min="0" :max="150" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="电话" prop="phone">
              <el-input v-model="patientForm.phone" placeholder="请输入电话号码" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="联系方式">
          <el-input v-model="patientForm.contact" placeholder="其他联系方式，QQ, 微信等" />
        </el-form-item>
        <el-form-item label="登记日期" prop="date">
          <el-date-picker v-model="patientForm.date" type="date" placeholder="选择日期" format="YYYY-MM-DD"
            value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="诊断医师" prop="doc">
          <el-select v-model="patientForm.doc" placeholder="请选择诊断医师" style="width: 100%" clearable>
            <el-option v-for="doctor in doctorStore.doctors" :key="doctor.id"
              :label="`${doctor.name} (${doctor.level})`" :value="doctor.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="过敏史">
          <el-input v-model="patientForm.AllergyHistory" type="textarea" :rows="3" placeholder="请输入过敏史信息" />
        </el-form-item>
        <el-form-item label="诊断" prop="detail">
          <el-input v-model="patientForm.detail" type="textarea" :rows="3" placeholder="请输入诊断信息" />
        </el-form-item>
        <el-form-item label="医嘱">
          <el-input v-model="patientForm.medical_advice" type="textarea" :rows="3" placeholder="请输入医嘱信息" />
        </el-form-item>
        <el-form-item label="治疗方案" prop="solution">
          <el-input v-model="patientForm.solution" type="textarea" :rows="3" placeholder="请输入治疗方案" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddDialog = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            {{ editingPatient ? '更新' : '添加' }}
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 患者详情查看对话框 -->
    <el-dialog v-model="showDetailDialog" title="患者详细信息" width="680px" :close-on-click-modal="false">
      <div v-if="viewingPatient" class="patient-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID">{{ viewingPatient.id }}</el-descriptions-item>
          <el-descriptions-item label="姓名">{{ viewingPatient.name }}</el-descriptions-item>
          <el-descriptions-item label="性别">
            <el-tag :type="viewingPatient.sex === '男' ? 'primary' : 'danger'" size="small">
              {{ viewingPatient.sex }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="年龄">{{ viewingPatient.age }} 岁</el-descriptions-item>
          <el-descriptions-item label="电话">{{ viewingPatient.phone }}</el-descriptions-item>
          <el-descriptions-item label="联系方式">{{ viewingPatient.contact || '无' }}</el-descriptions-item>
          <el-descriptions-item label="登记日期">{{ viewingPatient.date }}</el-descriptions-item>
          <el-descriptions-item label="诊断医师">{{ viewingPatient.doc }}</el-descriptions-item>
          <el-descriptions-item label="过敏史" :span="2">
            <div class="detail-text">{{ viewingPatient.AllergyHistory || '无' }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="诊断" :span="2">
            <div class="detail-text">{{ viewingPatient.detail }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="医嘱" :span="2">
            <div class="detail-text">{{ viewingPatient.medical_advice || '无' }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="治疗方案" :span="2">
            <div class="detail-text">{{ viewingPatient.solution || '无' }}</div>
          </el-descriptions-item>
        </el-descriptions>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showDetailDialog = false">关闭</el-button>
          <el-button type="primary" @click="handleEditFromDetail">编辑</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useDataStore } from '@/stores/data'
import { useDoctorStore } from '@/stores/doctor'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Refresh,
  Download,
  Search,
  RefreshLeft,
  Edit,
  Delete
} from '@element-plus/icons-vue'

const dataStore = useDataStore()
const doctorStore = useDoctorStore()

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const showAddDialog = ref(false)
const showDetailDialog = ref(false)
const editingPatient = ref(null)
const viewingPatient = ref(null)
const searchExpanded = ref(false)
const selectedRows = ref([])
const currentPage = ref(1)
const pageSize = ref(20)

// 表单引用
const formRef = ref()

// 搜索表单
const searchForm = ref({
  keyword: '',
  sex: '',
  dateRange: []
})

// 患者表单
const patientForm = ref({
  name: '',
  sex: '男',
  age: null,
  phone: '',
  contact: '',
  date: new Date().toISOString().split('T')[0],
  doc: '',
  AllergyHistory: '',
  detail: '',
  medical_advice: '',
  solution: ''
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
  age: [
    { type: 'number', min: 0, max: 150, message: '年龄必须在 0 到 150 之间', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  date: [
    { required: true, message: '请选择登记日期', trigger: 'change' }
  ],
  doc: [
    { required: true, message: '请选择诊断医师', trigger: 'change' }
  ],
  detail: [
    { required: true, message: '请输入诊断信息', trigger: 'blur' }
  ]
}

// 计算属性
const filteredData = computed(() => {
  // 确保dataStore.patientItems是数组
  if (!dataStore.patientItems || !Array.isArray(dataStore.patientItems)) {
    return []
  }
  let data = [...dataStore.patientItems]

  // 关键词搜索
  if (searchForm.value.keyword) {
    const keyword = searchForm.value.keyword.toLowerCase()
    data = data.filter(item =>
      item.name.toLowerCase().includes(keyword) ||
      item.phone.includes(keyword) ||
      item.detail.toLowerCase().includes(keyword)
    )
  }

  // 性别筛选
  if (searchForm.value.sex) {
    data = data.filter(item => item.sex === searchForm.value.sex)
  }

  // 日期范围筛选
  if (searchForm.value.dateRange && searchForm.value.dateRange.length === 2) {
    const [startDate, endDate] = searchForm.value.dateRange
    data = data.filter(item => item.date >= startDate && item.date <= endDate)
  }

  return data
})

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredData.value.slice(start, end)
})

// 方法
const handleRefresh = async () => {
  if (isUnmounted.value) return
  loading.value = true
  try {
    await dataStore.fetchAllPatients()
    if (!isUnmounted.value) {
      ElMessage.success('数据刷新成功')
    }
  } catch (error) {
    if (!isUnmounted.value) {
      ElMessage.error('数据刷新失败：' + error.message)
    }
  } finally {
    if (!isUnmounted.value) {
      loading.value = false
    }
  }
}

const handleBackup = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要备份数据库吗？',
      '确认备份',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await dataStore.backupDatabase(false)
    ElMessage.success('数据备份成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('数据备份失败：' + error.message)
    }
  }
}

const handleSearch = () => {
  currentPage.value = 1
}

const handleResetSearch = () => {
  searchForm.value = {
    keyword: '',
    sex: '',
    dateRange: []
  }
  currentPage.value = 1
}

const handleSelectionChange = (selection) => {
  selectedRows.value = selection
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
}

const handleCurrentChange = (page) => {
  currentPage.value = page
}

const handleEdit = (row) => {
  editingPatient.value = row
  patientForm.value = { ...row }
  showAddDialog.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除患者 "${row.name}" 的信息吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await dataStore.deletePatient(row.id)
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败：' + error.message)
    }
  }
}

const handleDialogClose = () => {
  editingPatient.value = null
  patientForm.value = {
    name: '',
    sex: '男',
    age: null,
    phone: '',
    contact: '',
    date: new Date().toISOString().split('T')[0],
    doc: '',
    AllergyHistory: '',
    detail: '',
    medical_advice: '',
    solution: ''
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

    // 确保字段类型正确，防止前端数据类型转换导致后端反序列化失败
    const submitData = {
      ...patientForm.value,
      name: String(patientForm.value.name || ''),
      phone: String(patientForm.value.phone || ''),
      age: String(patientForm.value.age || ''),
      contact: String(patientForm.value.contact || ''),
      doc: String(patientForm.value.doc || ''),
      detail: String(patientForm.value.detail || ''),
      solution: String(patientForm.value.solution || ''),
      medical_advice: String(patientForm.value.medical_advice || ''),
      AllergyHistory: String(patientForm.value.AllergyHistory || '')
    }

    if (editingPatient.value) {
      await dataStore.updatePatient(editingPatient.value.id, submitData)
      ElMessage.success('更新成功')
    } else {
      await dataStore.addPatient(submitData)
      ElMessage.success('添加成功')
    }

    showAddDialog.value = false
    handleDialogClose()
  } catch (error) {
    if (error.message) {
      ElMessage.error('操作失败：' + error.message)
    }
  } finally {
    submitting.value = false
  }
}

// 组件是否已卸载
const isUnmounted = ref(false)

// 初始化
onMounted(async () => {
  loading.value = true
  try {
    await dataStore.fetchAllPatients()
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
    if (!isUnmounted.value) {
      loading.value = false
    }
  }
})

// 查看患者详情
const handleViewDetail = (patient) => {
  viewingPatient.value = patient
  showDetailDialog.value = true
}

// 从详情对话框编辑患者
const handleEditFromDetail = () => {
  if (viewingPatient.value) {
    showDetailDialog.value = false
    handleEdit(viewingPatient.value)
  }
}

// 组件卸载时清理
onUnmounted(() => {
  isUnmounted.value = true
})
</script>

<style scoped>
.panel-container {
  padding: 20px;
  height: calc(100vh - 60px);
  overflow: auto;
}

.toolbar {
  margin-bottom: 16px;
}

.toolbar-card {
  border: 1px solid var(--el-border-color-lighter);
}

.toolbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 0;
}

.toolbar-left {
  display: flex;
  gap: 12px;
}

.search-section {
  margin-bottom: 16px;
}

.search-card {
  border: 1px solid var(--el-border-color-lighter);
}

.search-content {
  padding: 16px 0;
}

.table-section {
  margin-bottom: 20px;
}

.table-card {
  border: 1px solid var(--el-border-color-lighter);
}

.patient-detail {
  padding: 16px 0;
}

.detail-text {
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.5;
  max-height: 120px;
  overflow-y: auto;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
}

.search-section .search-card :deep(.el-card__body) {
  padding: 0;
}

.search-card .search-content {
  padding-left: 10px;
  padding-right: 10px;
  padding-bottom: 0;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .panel-container {
    padding: 12px;
  }

  .toolbar-content {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .toolbar-left {
    justify-content: center;
  }

  .search-content .el-form {
    flex-direction: column;
  }

  .search-content .el-form-item {
    margin-bottom: 16px;
  }
}
</style>