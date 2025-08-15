<template>
  <div class="template-container">
    <!-- 诊断模板 -->
    <div class="template-section">
      <el-card class="template-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span>诊断模板</span>
            <el-button type="primary" plain @click="showAddDialog('diagnosis')" :icon="Plus">
              新增诊断模板
            </el-button>
          </div>
        </template>
        <el-table :data="diagnosisTemplates" style="width: 100%" stripe v-loading="loading">
          <el-table-column prop="name" label="模板名称" width="200" />
          <el-table-column prop="content" label="模板内容" show-overflow-tooltip />
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
      </el-card>
    </div>

    <!-- 医嘱模板 -->
    <div class="template-section">
      <el-card class="template-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span>医嘱模板</span>
            <el-button type="primary" plain @click="showAddDialog('advice')" :icon="Plus">
              新增医嘱模板
            </el-button>
          </div>
        </template>
        <el-table :data="adviceTemplates" style="width: 100%" stripe v-loading="loading">
          <el-table-column prop="name" label="模板名称" width="200" />
          <el-table-column prop="content" label="模板内容" show-overflow-tooltip />
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
      </el-card>
    </div>

    <!-- 治疗方案模板 -->
    <div class="template-section">
      <el-card class="template-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span>治疗方案模板</span>
            <el-button type="primary" plain @click="showAddDialog('treatment')" :icon="Plus">
              新增治疗方案模板
            </el-button>
          </div>
        </template>
        <el-table :data="treatmentTemplates" style="width: 100%" stripe v-loading="loading">
          <el-table-column prop="name" label="模板名称" width="200" />
          <el-table-column prop="content" label="模板内容" show-overflow-tooltip />
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
      </el-card>
    </div>

    <!-- 新增/编辑模板对话框 -->
    <el-dialog
      v-model="showDialog"
      :title="isEdit ? '编辑模板' : '新增模板'"
      width="600px"
      @close="resetForm"
    >
      <el-form
        ref="templateFormRef"
        :model="templateForm"
        :rules="templateRules"
        label-width="100px"
      >
        <el-form-item label="模板类型" prop="type">
          <el-select v-model="templateForm.type" placeholder="选择模板类型" :disabled="isEdit">
            <el-option label="诊断" value="diagnosis" />
            <el-option label="医嘱" value="advice" />
            <el-option label="治疗方案" value="treatment" />
          </el-select>
        </el-form-item>
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="templateForm.name" placeholder="请输入模板名称" />
        </el-form-item>
        <el-form-item label="模板内容" prop="content">
          <el-input
            v-model="templateForm.content"
            type="textarea"
            :rows="6"
            placeholder="请输入模板内容"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showDialog = false" plain>取消</el-button>
          <el-button type="primary" plain @click="handleSubmit" :loading="submitting">
            {{ isEdit ? '更新' : '创建' }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElNotification, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import { GetAllTemplates, CreateTemplate, UpdateTemplate, DeleteTemplate } from '../../wailsjs/go/main/App'

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const showDialog = ref(false)
const isEdit = ref(false)
const templateFormRef = ref()
const allTemplates = ref([])

// 表单数据
const templateForm = reactive({
  id: null,
  name: '',
  type: '',
  content: ''
})

// 表单验证规则
const templateRules = {
  type: [{ required: true, message: '请选择模板类型', trigger: 'change' }],
  name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
  content: [{ required: true, message: '请输入模板内容', trigger: 'blur' }]
}

// 计算属性：按类型分组的模板
const diagnosisTemplates = computed(() => 
  allTemplates.value.filter(template => template.type === 'diagnosis')
)

const adviceTemplates = computed(() => 
  allTemplates.value.filter(template => template.type === 'advice')
)

const treatmentTemplates = computed(() => 
  allTemplates.value.filter(template => template.type === 'treatment')
)

// 获取所有模板数据
const fetchTemplates = async () => {
  try {
    loading.value = true
    const response = await GetAllTemplates()
    allTemplates.value = response || []
  } catch (error) {
    console.error('获取模板数据失败:', error)
    ElNotification.error({ message: '获取模板数据失败', duration: 1500, position: 'bottom-right' })
  } finally {
    loading.value = false
  }
}

// 显示新增对话框
const showAddDialog = (type) => {
  isEdit.value = false
  templateForm.type = type
  showDialog.value = true
}

// 编辑模板
const handleEdit = (row) => {
  isEdit.value = true
  templateForm.id = row.id
  templateForm.name = row.name
  templateForm.type = row.type
  templateForm.content = row.content
  showDialog.value = true
}

// 删除模板
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除模板 "${row.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await DeleteTemplate(row.id)
    
    // 从本地数组中删除
    const index = allTemplates.value.findIndex(template => template.id === row.id)
    if (index > -1) {
      allTemplates.value.splice(index, 1)
    }
    
    ElNotification.success({ message: '删除成功', duration: 1500, position: 'bottom-right' })
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除模板失败:', error)
      ElNotification.error({ message: '删除模板失败', duration: 1500, position: 'bottom-right' })
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  try {
    await templateFormRef.value.validate()
    submitting.value = true
    
    if (isEdit.value) {
      // 更新模板
      const response = await UpdateTemplate(templateForm)
      
      // 更新本地数据
      const index = allTemplates.value.findIndex(template => template.id === templateForm.id)
      if (index > -1) {
        allTemplates.value[index] = response || { ...templateForm }
      }
      
      ElNotification.success({ message: '更新成功', duration: 1500, position: 'bottom-right' })
    } else {
      // 检查名称是否重复
      const exists = allTemplates.value.some(template => 
        template.name === templateForm.name && template.type === templateForm.type
      )
      if (exists) {
        ElNotification.error({ message: '模板名称已存在，请使用其他名称', duration: 1500, position: 'bottom-right' })
        return
      }
      
      // 创建新模板
      const response = await CreateTemplate(templateForm)
      
      // 添加到本地数据
      if (response) {
        allTemplates.value.push(response)
      }
      
      ElNotification.success({ message: '创建成功', duration: 1500, position: 'bottom-right' })
    }
    
    showDialog.value = false
    resetForm()
  } catch (error) {
    if (error.message && error.message.includes('duplicate')) {
      ElNotification.error({ message: '模板名称已存在，请使用其他名称', duration: 1500, position: 'bottom-right' })
    } else {
      console.error('提交失败:', error)
      ElNotification.error({ message: '操作失败', duration: 1500, position: 'bottom-right' })
    }
  } finally {
    submitting.value = false
  }
}

// 重置表单
const resetForm = () => {
  templateForm.id = null
  templateForm.name = ''
  templateForm.type = ''
  templateForm.content = ''
  templateFormRef.value?.resetFields()
}

// 组件挂载时获取数据
onMounted(() => {
  fetchTemplates()
})
</script>

<style scoped>
.template-container {
  height: calc(100vh - 120px);
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 16px;
  overflow-y: auto;
}

.template-section {
  flex-shrink: 0;
}

.template-card {
  border: 1px solid var(--el-border-color-lighter);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .template-container {
    padding: 12px;
    gap: 16px;
  }
  
  .card-header {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
}
</style>