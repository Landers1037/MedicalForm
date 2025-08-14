<template>
  <div class="panel-container">
    <!-- 搜索和操作栏 -->
    <div class="search-section">
      <el-card class="search-card" shadow="never">
        <template #header>
          <div class="card-header">
            <div class="header-left">
              <span>高级搜索</span>
              <el-button type="primary" link @click="searchExpanded = !searchExpanded">
                {{ searchExpanded ? '收起' : '展开' }}
              </el-button>
            </div>
            <div class="header-right">
              <div class="toolbar-buttons">
                <el-button plain type="primary" @click="showAddDialog = true" :icon="Plus">
                  添加就诊记录
                </el-button>
                <el-button plain @click="handleRefresh" :icon="Refresh">
                  刷新
                </el-button>
              </div>
              <div class="search-input">
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
          </div>
        </template>

        <!-- 展开状态下显示的高级搜索 -->
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
                <el-button plain type="primary" @click="handleSearch" :icon="Search">
                  搜索
                </el-button>
                <el-button plain @click="handleResetSearch" :icon="RefreshLeft">
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
          <el-table-column v-if="visibleColumns.includes('id')" prop="id" label="ID" width="80" sortable>
            <template #default="{ row }">
              <el-button type="primary" link @click="handleViewDetail(row)">
                {{ row.id }}
              </el-button>
            </template>
          </el-table-column>
          <el-table-column v-if="visibleColumns.includes('name')" prop="name" label="姓名" width="120" sortable />
          <el-table-column v-if="visibleColumns.includes('sex')" prop="sex" label="性别" width="80" sortable>
            <template #default="{ row }">
              <el-tag :type="row.sex === '男' ? 'primary' : 'danger'" size="small">
                {{ row.sex }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column v-if="visibleColumns.includes('age')" prop="age" label="年龄" width="80" sortable />
          <el-table-column v-if="visibleColumns.includes('phone')" prop="phone" label="电话" width="140" />
          <el-table-column v-if="visibleColumns.includes('contact')" prop="contact" label="联系方式" width="150"
            show-overflow-tooltip />
          <el-table-column v-if="visibleColumns.includes('address')" prop="address" label="家庭住址" width="150"
            show-overflow-tooltip />
          <el-table-column v-if="visibleColumns.includes('parent')" prop="parent" label="监护人" width="100"
            show-overflow-tooltip />
          <el-table-column v-if="visibleColumns.includes('work')" prop="work" label="工作" width="100"
            show-overflow-tooltip />
          <el-table-column v-if="visibleColumns.includes('date')" prop="date" label="登记日期" width="120" sortable />
          <el-table-column v-if="visibleColumns.includes('ill_time')" prop="ill_time" label="患病时间" width="120"
            show-overflow-tooltip />
          <el-table-column v-if="visibleColumns.includes('doc')" prop="doc" label="诊断医师" width="120"
            show-overflow-tooltip />
          <el-table-column v-if="visibleColumns.includes('detail')" prop="detail" label="诊断"
            show-overflow-tooltip />
          <el-table-column v-if="visibleColumns.includes('solution')" prop="solution" label="治疗方案"
            show-overflow-tooltip />
          <el-table-column v-if="visibleColumns.includes('medical_advice')" prop="medical_advice" label="医嘱"
            show-overflow-tooltip />
          <el-table-column v-if="visibleColumns.includes('fee')" prop="fee" label="费用" width="100" sortable />
          <el-table-column v-if="visibleColumns.includes('money')" prop="money" label="费用备注" width="120"
            show-overflow-tooltip />
          <el-table-column label="操作" width="220" fixed="right">
            <template #header>
              <span>操作</span>
              <el-button type="primary" link @click="showColumnConfig = true" :icon="Setting" size="small"
                style="margin-left: 8px;">
                设置
              </el-button>
            </template>
            <template #default="{ row }">
              <el-button type="primary" link @click="handleEdit(row)" :icon="Edit">
                编辑
              </el-button>
              <el-button type="success" link @click="handleCopy(row)" :icon="CopyDocument">
                复制
              </el-button>
              <el-button type="danger" link @click="handleDelete(row)" :icon="Delete">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-wrapper">
          <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 30, 40, 50, 100]"
            :total="filteredData.length" layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange" @current-change="handleCurrentChange" />
        </div>
      </el-card>
    </div>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="showAddDialog" :title="editingPatient ? '编辑患者就诊记录' : '添加患者就诊记录'" :width="dialogWidth"
      :close-on-click-modal="false" @close="handleDialogClose">
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
        <el-form-item label="家庭住址">
          <el-input v-model="patientForm.address" placeholder="请输入家庭住址" />
        </el-form-item>
        <el-row :gutter="8">
          <el-col :span="12">
            <el-form-item label="监护人">
              <el-input v-model="patientForm.parent" placeholder="请输入监护人姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="工作">
              <div style="display: flex; gap: 8px;">
                <el-input v-model="patientForm.work" placeholder="请输入工作" style="flex: 1;width: 200px;" />
                <el-select v-model="selectedWorkOption" placeholder="选快速择" style="width: 100px;" @change="fillWorkOption">
                  <el-option label="学生" value="学生" />
                  <el-option label="职工" value="职工" />
                  <el-option label="自由职业" value="自由职业" />
                  <el-option label="退休" value="退休" />
                  <el-option label="农民" value="农民" />
                  <el-option label="个体户" value="个体户" />
                  <el-option label="公务员" value="公务员" />
                  <el-option label="教师" value="教师" />
                  <el-option label="医生" value="医生" />
                  <el-option label="工人" value="工人" />
                </el-select>
              </div>
            </el-form-item>
          </el-col>
        </el-row>
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
        <el-form-item label="患病时间">
          <el-input v-model="patientForm.ill_time" placeholder="请输入患病时间，如：3天、1周、2个月等" />
        </el-form-item>
        <el-form-item label="主诉及病史">
          <el-input v-model="patientForm.AllergyHistory" type="textarea" :rows="3" placeholder="请输入病史信息" />
        </el-form-item>
        <el-form-item label="诊断" prop="detail">
          <div class="template-input-group">
            <el-select v-model="selectedDiagnosisTemplate" placeholder="从模板中填充" clearable size="small"
              style="width: 200px; margin-bottom: 8px;" popper-class="template-select-item"
              @change="fillDiagnosisTemplate">
              <el-option v-for="template in templateStore.diagnosisTemplates" :key="template.id" :label="template.name"
                :value="template.id" />
            </el-select>
          </div>
          <el-input v-model="patientForm.detail" type="textarea" :rows="3" placeholder="请输入诊断信息" />
        </el-form-item>
        <el-form-item label="检查及医嘱">
          <div class="template-input-group">
            <el-select v-model="selectedAdviceTemplate" placeholder="从模板中填充" clearable size="small"
              style="width: 200px; margin-bottom: 8px;" popper-class="template-select-item"
              @change="fillAdviceTemplate">
              <el-option v-for="template in templateStore.adviceTemplates" :key="template.id" :label="template.name"
                :value="template.id" />
            </el-select>
          </div>
          <el-input v-model="patientForm.medical_advice" type="textarea" :rows="3" placeholder="请输入医嘱信息" />
        </el-form-item>
        <el-form-item label="治疗方案" prop="solution">
          <div class="template-input-group">
            <div class="template-row">
              <el-select v-model="selectedTreatmentTemplate" placeholder="从模板中填充" clearable size="small"
                style="width: 200px;" popper-class="template-select-item" @change="fillTreatmentTemplate">
                <el-option v-for="template in templateStore.treatmentTemplates" :key="template.id"
                  :label="template.name" :value="template.id" />
              </el-select>
              <el-select v-model="selectedMedicine" placeholder="开药" clearable filterable remote
                :remote-method="searchMedicines" :loading="medicineLoading" size="small"
                style="width: 200px; margin-left: 10px;" popper-class="medicine-select-item"
                @change="addMedicineToSolution">
                <el-option v-for="medicine in filteredMedicines" :key="medicine.id" :label="medicine.name"
                  :value="medicine.id" />
              </el-select>
            </div>
          </div>
          <el-input v-model="patientForm.solution" type="textarea" :rows="3" placeholder="请输入治疗方案" />
        </el-form-item>
        <el-form-item label="收费">
          <el-input-number v-model="patientForm.fee" :precision="2" style="width: 100%" placeholder="请输入收费金额">
            <template #suffix>
              <span style="color: #909399;">人民币(元)</span>
            </template>
          </el-input-number>
        </el-form-item>
        <el-form-item label="费用备注">
          <el-input v-model="patientForm.money" placeholder="请输入费用备注信息，如：已收费、欠费50元等" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button plain @click="showAddDialog = false">取消</el-button>
          <el-button plain type="primary" @click="handleSubmit" :loading="submitting">
            {{ editingPatient ? '更新' : '添加' }}
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 患者详情查看对话框 -->
    <el-dialog v-model="showDetailDialog" :title="showTimeline ? '患者时间线' : '患者详细信息'" width="680px"
      :close-on-click-modal="false">
      <!-- 患者详细信息视图 -->
      <div v-if="viewingPatient && !showTimeline" class="patient-detail">
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
          <el-descriptions-item label="家庭住址">{{ viewingPatient.address || '无' }}</el-descriptions-item>
          <el-descriptions-item label="监护人">{{ viewingPatient.parent || '无' }}</el-descriptions-item>
          <el-descriptions-item label="工作">{{ viewingPatient.work || '无' }}</el-descriptions-item>
          <el-descriptions-item label="登记日期">{{ viewingPatient.date }}</el-descriptions-item>
          <el-descriptions-item label="患病时间">{{ viewingPatient.ill_time || '无' }}</el-descriptions-item>
          <el-descriptions-item label="诊断医师">{{ viewingPatient.doc }}</el-descriptions-item>
          <el-descriptions-item label="主诉及病史" :span="2">
            <div class="detail-text">{{ viewingPatient.AllergyHistory || '无' }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="诊断" :span="2">
            <div class="detail-text">{{ viewingPatient.detail }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="检查及医嘱" :span="2">
            <div class="detail-text">{{ viewingPatient.medical_advice || '无' }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="治疗方案" :span="2">
            <div class="detail-text">{{ viewingPatient.solution || '无' }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="收费">
            <div class="detail-text">
              {{ viewingPatient.fee ? viewingPatient.fee + ' 元' : '无' }}
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="费用备注">
            <div class="detail-text">
              {{ viewingPatient.money || '无' }}
            </div>
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- 患者时间线视图 -->
      <div v-if="viewingPatient && showTimeline" class="patient-timeline">
        <div class="timeline-header">
          <el-button plain type="primary" link @click="showTimeline = false" :icon="ArrowLeft">
            返回患者详情
          </el-button>
        </div>
        <div class="timeline-content">
          <el-timeline>
            <el-timeline-item v-for="record in patientTimelineData" :key="record.id"
              :timestamp="record.date + ' ' + (record.time || '')" placement="top">
              <el-card class="timeline-card">
                <div class="timeline-item-header">
                  <h4>{{ record.date }} 就诊记录</h4>
                  <el-tag size="small" type="info">{{ record.doc }}</el-tag>
                </div>
                <div class="timeline-item-content">
                  <p><strong>患病时间：</strong>{{ record.ill_time || '无' }}</p>
                  <p><strong>主诉及病史：</strong>{{ record.AllergyHistory || '无' }}</p>
                  <p><strong>诊断：</strong>{{ record.detail }}</p>
                  <p><strong>检查及医嘱：</strong>{{ record.medical_advice || '无' }}</p>
                  <p><strong>治疗方案：</strong>{{ record.solution || '无' }}</p>
                </div>
              </el-card>
            </el-timeline-item>
          </el-timeline>
          <div v-if="patientTimelineData.length === 0" class="empty-timeline">
            <el-empty description="暂无就诊记录" />
          </div>
        </div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button plain @click="showDetailDialog = false">关闭</el-button>
          <el-button v-if="!showTimeline" type="success" plain @click="handleShowTimeline">时间线</el-button>
          <el-button v-if="!showTimeline" type="primary" plain @click="handleEditFromDetail">编辑</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 列配置对话框 -->
    <el-dialog v-model="showColumnConfig" title="表格列配置" class="column-config" :close-on-click-modal="false">
      <div class="column-config-content">
        <p class="config-tip">拖拽调整列的显示顺序，勾选控制列的显示/隐藏</p>
        <el-table :data="columnConfigs" row-key="columnKey" class="config-table">
          <el-table-column label="显示" width="60">
            <template #default="{ row }">
              <el-checkbox v-model="row.visible" />
            </template>
          </el-table-column>
          <el-table-column label="列名" prop="label" />
          <el-table-column label="字段" prop="columnKey" />
        </el-table>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button plain @click="handleCancelColumnConfig">取消</el-button>
          <el-button type="primary" plain @click="handleSaveColumnConfig">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useDataStore } from '@/stores/data'
import { useDoctorStore } from '@/stores/doctor'
import { useTemplateStore } from '@/stores/template'
import { ElMessage, ElMessageBox } from 'element-plus'
import { GetAllMedicines } from '../../wailsjs/go/main/App'
import {
  Plus,
  Refresh,
  Search,
  RefreshLeft,
  Edit,
  Delete,
  CopyDocument,
  ArrowLeft,
  Setting
} from '@element-plus/icons-vue'

const dataStore = useDataStore()
const doctorStore = useDoctorStore()
const templateStore = useTemplateStore()

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const showAddDialog = ref(false)
const showDetailDialog = ref(false)
const editingPatient = ref(null)
const viewingPatient = ref(null)
const showTimeline = ref(false)
const patientTimelineData = ref([])
const searchExpanded = ref(false)
const selectedRows = ref([])
const currentPage = ref(1)
const pageSize = ref(20)

// 列配置相关
const showColumnConfig = ref(false)
const columnConfigs = ref([])
const visibleColumns = ref([])

// 模板选择相关
const selectedDiagnosisTemplate = ref(null)
const selectedAdviceTemplate = ref(null)
const selectedTreatmentTemplate = ref(null)

// 药品搜索相关
const selectedMedicine = ref(null)
const filteredMedicines = ref([])
const medicineLoading = ref(false)

// 工作选项相关
const selectedWorkOption = ref(null)

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
  address: '',
  parent: '',
  work: '',
  date: new Date().toISOString().split('T')[0],
  doc: '',
  ill_time: '',
  AllergyHistory: '',
  detail: '',
  medical_advice: '',
  solution: '',
  fee: null,
  money: ''
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

  // 关键词搜索 - 根据名称进行精准匹配
  if (searchForm.value.keyword) {
    const keyword = searchForm.value.keyword.trim()
    data = data.filter(item =>
      item.name === keyword ||
      item.phone.includes(keyword) ||
      item.detail.toLowerCase().includes(keyword.toLowerCase())
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

// 弹框宽度响应式计算
const dialogWidth = computed(() => {
  if (typeof window !== 'undefined') {
    const windowWidth = window.innerWidth
    if (windowWidth >= 1920) {
      return '900px'  // 大屏幕时更宽
    } else if (windowWidth >= 1440) {
      return '800px'  // 中等屏幕
    } else if (windowWidth >= 1024) {
      return '720px'  // 普通屏幕
    } else {
      return '90%'    // 小屏幕时使用百分比
    }
  }
  return '680px' // 默认宽度
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

const handleCopy = (row) => {
  // 复制患者基本信息，但清空诊断相关信息用于新的就诊记录
  editingPatient.value = null // 确保是添加模式
  patientForm.value = {
    name: row.name,
    sex: row.sex,
    age: Number(row.age) || null, // 确保年龄是数字类型
    phone: row.phone,
    contact: row.contact,
    address: row.address,
    parent: row.parent || '', // 复制监护人信息
    work: row.work || '', // 复制工作信息
    date: new Date().toISOString().split('T')[0], // 设置为当前日期
    doc: row.doc,
    ill_time: row.ill_time || '', // 复制患病时间
    AllergyHistory: row.AllergyHistory,
    detail: '', // 清空诊断
    medical_advice: '', // 清空医嘱
    solution: '', // 清空治疗方案
    fee: null, // 清空收费
    money: '' // 清空费用备注
  }
  showAddDialog.value = true
  ElMessage.success('已复制患者基本信息，请填写新的就诊记录')
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
    address: '',
    parent: '',
    work: '',
    date: new Date().toISOString().split('T')[0],
    doc: '',
    ill_time: '',
    AllergyHistory: '',
    detail: '',
    medical_advice: '',
    solution: '',
    fee: null,
    money: ''
  }
  // 重置模板选择
  selectedDiagnosisTemplate.value = null
  selectedAdviceTemplate.value = null
  selectedTreatmentTemplate.value = null
  selectedWorkOption.value = null

  if (formRef.value) {
    formRef.value.resetFields()
  }
}

/**
 * 填充诊断模板内容
 */
const fillDiagnosisTemplate = (templateId) => {
  if (templateId) {
    const content = templateStore.getTemplateContent('diagnosis', templateId)
    patientForm.value.detail = content
  }
}

/**
 * 填充医嘱模板内容
 */
const fillAdviceTemplate = (templateId) => {
  if (templateId) {
    const content = templateStore.getTemplateContent('advice', templateId)
    patientForm.value.medical_advice = content
  }
}

/**
 * 填充治疗方案模板内容
 */
const fillTreatmentTemplate = (templateId) => {
  if (templateId) {
    const content = templateStore.getTemplateContent('treatment', templateId)
    patientForm.value.solution = content
  }
}

/**
 * 填充工作选项
 */
const fillWorkOption = (workOption) => {
  if (workOption) {
    patientForm.value.work = workOption
    selectedWorkOption.value = null // 清空选择
  }
}

/**
 * 搜索药品
 * @param {string} query - 搜索关键字
 */
const searchMedicines = async (query) => {
  if (!query) {
    filteredMedicines.value = []
    return
  }

  try {
    medicineLoading.value = true
    const medicines = await GetAllMedicines()

    // 根据关键字过滤药品
    filteredMedicines.value = (medicines || []).filter(medicine =>
      medicine.name.toLowerCase().includes(query.toLowerCase())
    )
  } catch (error) {
    console.error('搜索药品失败:', error)
    ElMessage.error('搜索药品失败：' + error.message)
    filteredMedicines.value = []
  } finally {
    medicineLoading.value = false
  }
}

/**
 * 添加药品到治疗方案
 * @param {number} medicineId - 药品ID
 */
const addMedicineToSolution = async (medicineId) => {
  if (!medicineId) return

  try {
    const medicines = await GetAllMedicines()
    const medicine = medicines.find(m => m.id === medicineId)

    if (medicine) {
      // 在治疗方案中添加药品信息
      const medicineInfo = `${medicine.name} 规格:(${medicine.specification}) - 价格:${medicine.price}`

      if (patientForm.value.solution) {
        patientForm.value.solution += `\n${medicineInfo}`
      } else {
        patientForm.value.solution = medicineInfo
      }

      // 清空选择
      selectedMedicine.value = null
      filteredMedicines.value = []
    }
  } catch (error) {
    console.error('添加药品失败:', error)
    ElMessage.error('添加药品失败：' + error.message)
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
      address: String(patientForm.value.address || ''),
      parent: String(patientForm.value.parent || ''),
      work: String(patientForm.value.work || ''),
      doc: String(patientForm.value.doc || ''),
      ill_time: String(patientForm.value.ill_time || ''),
      detail: String(patientForm.value.detail || ''),
      solution: String(patientForm.value.solution || ''),
      medical_advice: String(patientForm.value.medical_advice || ''),
      AllergyHistory: String(patientForm.value.AllergyHistory || '')
    }

    if (editingPatient.value) {
      console.log(submitData);

      await dataStore.updatePatient(submitData)
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

// 窗口大小变化监听器
const handleResize = () => {
  // 触发dialogWidth重新计算
}

// 初始化
onMounted(async () => {
  loading.value = true
  try {
    await dataStore.fetchAllPatients()
    await doctorStore.fetchAllDoctors()
    await templateStore.initTemplates()
    await loadColumnConfigs() // 加载列配置
    
    // 添加窗口大小变化监听器
    window.addEventListener('resize', handleResize)
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
  showTimeline.value = false
  patientTimelineData.value = []
  showDetailDialog.value = true
}

// 从详情对话框编辑患者
const handleEditFromDetail = () => {
  if (viewingPatient.value) {
    showDetailDialog.value = false
    handleEdit(viewingPatient.value)
  }
}

// 显示患者时间线
const handleShowTimeline = async () => {
  if (!viewingPatient.value) return

  try {
    // 获取该患者的所有就诊记录
    await loadPatientTimeline(viewingPatient.value.name)
    showTimeline.value = true
  } catch (error) {
    ElMessage.error('获取时间线数据失败：' + error.message)
  }
}

// 加载患者时间线数据
const loadPatientTimeline = async (patientName) => {
  try {
    // 使用SearchPatients API根据姓名精准搜索该患者的所有记录
    const { SearchPatients } = await import('../../wailsjs/go/main/App')
    const records = await SearchPatients('name', patientName)

    // 按日期和时间排序，从最新到最旧
    patientTimelineData.value = records.sort((a, b) => {
      const dateA = new Date(a.date + ' ' + (a.time || '00:00:00'))
      const dateB = new Date(b.date + ' ' + (b.time || '00:00:00'))
      return dateB - dateA
    })
  } catch (error) {
    console.error('加载患者时间线失败:', error)
    patientTimelineData.value = []
    throw error
  }
}

// 列配置相关方法
const initColumnConfigs = () => {
  // 定义默认列配置
  const defaultColumns = [
    { columnKey: 'id', label: 'ID', visible: true, sortOrder: 0 },
    { columnKey: 'name', label: '姓名', visible: true, sortOrder: 1 },
    { columnKey: 'sex', label: '性别', visible: true, sortOrder: 2 },
    { columnKey: 'age', label: '年龄', visible: true, sortOrder: 3 },
    { columnKey: 'phone', label: '电话', visible: true, sortOrder: 4 },
    { columnKey: 'address', label: '家庭住址', visible: true, sortOrder: 5 },
    { columnKey: 'parent', label: '监护人', visible: true, sortOrder: 6 },
    { columnKey: 'work', label: '工作', visible: true, sortOrder: 7 },
    { columnKey: 'detail', label: '诊断', visible: true, sortOrder: 8 },
    { columnKey: 'contact', label: '联系方式', visible: false, sortOrder: 9 },
    { columnKey: 'date', label: '登记日期', visible: false, sortOrder: 10 },
    { columnKey: 'ill_time', label: '患病时间', visible: false, sortOrder: 11 },
    { columnKey: 'doc', label: '诊断医师', visible: false, sortOrder: 12 },
    { columnKey: 'solution', label: '治疗方案', visible: false, sortOrder: 13 },
    { columnKey: 'medical_advice', label: '医嘱', visible: false, sortOrder: 14 },
    { columnKey: 'fee', label: '费用', visible: false, sortOrder: 15 },
    { columnKey: 'money', label: '费用备注', visible: false, sortOrder: 16 }
  ]

  columnConfigs.value = [...defaultColumns]
  updateVisibleColumns()
}

const loadColumnConfigs = async () => {
  try {
    const { GetColumnConfigs } = await import('../../wailsjs/go/main/App')
    const configs = await GetColumnConfigs()

    if (configs && configs.length > 0) {
      // 从后端获取配置成功
      columnConfigs.value = configs.map(config => ({
        columnKey: config.column_key,
        label: getColumnLabel(config.column_key),
        visible: config.visible,
        sortOrder: config.sort_order
      })).sort((a, b) => a.sortOrder - b.sortOrder)
    } else {
      // 后端无配置，使用默认配置并初始化到后端
      initColumnConfigs()
      await initDefaultColumnConfigs()
    }
  } catch (error) {
    console.error('获取列配置失败:', error)
    // 获取失败，使用前端默认配置
    initColumnConfigs()
  }
  updateVisibleColumns()
}

const initDefaultColumnConfigs = async () => {
  try {
    const { InitDefaultColumnConfigs } = await import('../../wailsjs/go/main/App')
    await InitDefaultColumnConfigs()
  } catch (error) {
    console.error('初始化默认列配置失败:', error)
  }
}

const getColumnLabel = (columnKey) => {
  const labelMap = {
    id: 'ID',
    name: '姓名',
    sex: '性别',
    age: '年龄',
    phone: '电话',
    address: '家庭住址',
    parent: '监护人',
    work: '工作',
    detail: '诊断',
    contact: '联系方式',
    date: '登记日期',
    ill_time: '患病时间',
    doc: '诊断医师',
    solution: '治疗方案',
    medical_advice: '医嘱',
    fee: '费用',
    money: '费用备注'
  }
  return labelMap[columnKey] || columnKey
}

const updateVisibleColumns = () => {
  visibleColumns.value = columnConfigs.value
    .filter(config => config.visible)
    .sort((a, b) => a.sortOrder - b.sortOrder)
    .map(config => config.columnKey)
}

const handleSaveColumnConfig = async () => {
  try {
    const { SaveColumnConfigs } = await import('../../wailsjs/go/main/App')
    const configsToSave = columnConfigs.value.map((config, index) => ({
      column_key: config.columnKey,
      visible: config.visible,
      sort_order: index
    }))

    await SaveColumnConfigs(configsToSave)
    updateVisibleColumns()
    showColumnConfig.value = false
    ElMessage.success('列配置保存成功')
  } catch (error) {
    console.error('保存列配置失败:', error)
    ElMessage.error('保存列配置失败：' + error.message)
  }
}

const handleCancelColumnConfig = () => {
  showColumnConfig.value = false
  // 重新加载配置，恢复原始状态
  loadColumnConfigs()
}

// 组件卸载时清理
onUnmounted(() => {
  isUnmounted.value = true
  // 移除窗口大小变化监听器
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.panel-container {
  padding: 16px;
  height: calc(100vh - 115px);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.toolbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
}

.toolbar-left {
  display: flex;
  gap: 12px;
}

.search-section {
  margin-bottom: 16px;
    user-select: none;
}

.search-card {
  border: 1px solid var(--el-border-color-lighter);
}

.search-content {
  padding: 16px 0;
}

.table-section {
  margin-bottom: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.table-card {
  border: 1px solid var(--el-border-color-lighter);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.table-card :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 16px;
}

.table-card .el-table {
  flex: 1;
  overflow-x: auto;
}

.table-card :deep(.el-table__body-wrapper) {
  max-height: calc(100vh - 350px);
  overflow-y: auto;
  overflow-x: auto;
}

/* 确保分页器不会被遮挡 */
.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  flex-shrink: 0;
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

/* 时间线样式 */
.patient-timeline {
  padding: 16px 0;
}

.timeline-header {
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.timeline-content {
  max-height: 460px;
  overflow-y: auto;
}

.timeline-card {
  margin-bottom: 10px;
}

.timeline-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.timeline-item-header h4 {
  margin: 0;
  color: var(--el-text-color-primary);
  font-size: 16px;
}

.timeline-item-content p {
  margin: 8px 0;
  line-height: 1.5;
  color: var(--el-text-color-regular);
}

.timeline-item-content strong {
  color: var(--el-text-color-primary);
  font-weight: 600;
}

.empty-timeline {
  text-align: center;
  padding: 40px 0;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
  width: 100%;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.toolbar-buttons {
  display: flex;
  align-items: center;
  gap: 8px;
}

.search-input {
  display: flex;
  align-items: center;
}

.search-section .search-card :deep(.el-card__body) {
  padding: 0;
}

.search-card .search-content {
  padding-left: 10px;
  padding-right: 10px;
  padding-bottom: 0;
}



.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.template-input-group {
  margin-bottom: 8px;
}

.template-input-group .el-select {
  width: 200px;
  font-size: 12px;
}

.template-row {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.template-select-item.el-select-dropdown .el-select-dropdown__item {
  font-size: 12px;
}

/* 列配置对话框样式 */
:deep(.column-config) {
  width: 500px;
}

:deep(.column-config .el-dialog__body) {
  padding: 0 16px;
  max-height: 480px;
  overflow-y: auto;
}

.column-config-content {
  padding: 8px 0;
}

.config-tip {
  margin: 0 0 16px 0;
  color: var(--el-text-color-regular);
  font-size: 14px;
  line-height: 1.5;
}

.config-table {
  border: 1px solid var(--el-border-color-light);
  border-radius: 6px;
}

.config-table :deep(.el-table__header) {
  background-color: var(--el-fill-color-light);
}

.config-table :deep(.el-table__row:hover) {
  background-color: var(--el-fill-color-lighter);
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