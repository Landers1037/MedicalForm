<template>
  <div class="setting-container">
    <!-- 主题设置 -->
    <div class="setting-section">
      <el-card class="setting-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <el-icon class="header-icon">
              <Brush />
            </el-icon>
            <span class="card-title">主题设置</span>
          </div>
        </template>
        <div class="setting-content">
          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">主题模式</span>
              <span class="label-desc">选择明亮或暗黑主题模式</span>
            </div>
            <div class="setting-control">
              <el-switch
                v-model="themeStore.isDark"
                @change="handleThemeChange"
                :active-icon="Moon"
                :inactive-icon="Sunny"
                active-text="暗黑模式"
                inactive-text="明亮模式"
                inline-prompt
                style="--el-switch-on-color: #2c2c2c; --el-switch-off-color: #409EFF"
              />
            </div>
          </div>
          
          <el-divider />
          
          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">主题颜色</span>
              <span class="label-desc">选择系统的主色调</span>
            </div>
            <div class="setting-control">
              <div class="color-options">
                <div
                  v-for="color in colorOptions"
                  :key="color.name"
                  class="color-option"
                  :class="{ active: themeStore.primaryColor === color.value }"
                  @click="handleColorChange(color.value)"
                >
                  <div class="color-circle" :style="{ backgroundColor: color.value }">
                    <el-icon v-if="themeStore.primaryColor === color.value" class="check-icon">
                      <Check />
                    </el-icon>
                  </div>
                  <span class="color-name">{{ color.name }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 数据库设置 -->
    <div class="setting-section">
      <el-card class="setting-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <el-icon class="header-icon">
              <Coin />
            </el-icon>
            <span class="card-title">数据库设置</span>
          </div>
        </template>
        <div class="setting-content">
          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">数据库信息</span>
              <span class="label-desc">当前数据库的基本信息</span>
            </div>
            <div class="setting-control">
              <div class="db-info">
                <div class="info-item">
                  <span class="info-label">数据库大小：</span>
                  <span class="info-value">{{ formatFileSize(dataStore.dbInfo.size) }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">记录数量：</span>
                  <span class="info-value">{{ dataStore.dbInfo.count }} 条</span>
                </div>
                <div class="info-item">
                  <span class="info-label">最后备份：</span>
                  <span class="info-value">{{ lastBackupTime }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">数据库路径：</span>
                  <div class="path-control">
                    <el-input
                      v-model="dbPath"
                      placeholder="请选择数据库文件路径"
                      readonly
                      class="path-input"
                    />
                    <el-button
                      type="primary"
                      :icon="FolderOpened"
                      @click="handleBrowseDatabase"
                      class="browse-btn"
                    >
                      浏览
                    </el-button>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          <el-divider />
          
          <div class="setting-item">
            <div class="setting-label" style="align-items: center;">
              <span class="label-text">数据库操作</span>
              <span class="label-desc">备份、重置或清空数据库</span>
            </div>
            <div class="setting-control">
              <div class="db-actions">
                <el-button type="success" @click="handleMigrateOldDatabase" :icon="Upload">
                  合并旧版本
                </el-button>
                <el-button type="primary" @click="handleBackup" :icon="Download">
                  备份数据库
                </el-button>
                <el-button type="warning" @click="handleBackupAndClear" :icon="RefreshLeft">
                  备份并清空
                </el-button>
                <el-button type="danger" @click="handleClearDatabase" :icon="Delete">
                  清空数据库
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 系统信息 -->
    <div class="setting-section">
      <el-card class="setting-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <el-icon class="header-icon">
              <Monitor />
            </el-icon>
            <span class="card-title">系统信息</span>
          </div>
        </template>
        <div class="setting-content">
          <div class="system-info">
            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">系统名称：</span>
                <span class="info-value">{{ appInfo.title }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">版本号：</span>
                <span class="info-value">v{{ appInfo.version }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">技术栈：</span>
                <span class="info-value">Go + Wails + Vue3</span>
              </div>
              <div class="info-item">
                <span class="info-label">UI框架：</span>
                <span class="info-value">Element Plus</span>
              </div>
              <div class="info-item">
                <span class="info-label">数据库：</span>
                <span class="info-value">SQLite3</span>
              </div>
              <div class="info-item">
                <span class="info-label">开发者：</span>
                <span class="info-value">Landers</span>
              </div>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 导入导出 -->
    <div class="setting-section">
      <el-card class="setting-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <el-icon class="header-icon">
              <FolderOpened />
            </el-icon>
            <span class="card-title">数据导入导出</span>
          </div>
        </template>
        <div class="setting-content">
          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">数据导出</span>
              <span class="label-desc">将患者数据导出为Excel文件</span>
            </div>
            <div class="setting-control">
              <el-button type="success" @click="handleExport" :icon="Upload">
                导出数据
              </el-button>
            </div>
          </div>
          
          <el-divider />
          
          <div class="setting-item">
            <div class="setting-label">
              <span class="label-text">数据导入</span>
              <span class="label-desc">从Excel文件导入患者数据</span>
            </div>
            <div class="setting-control">
              <el-upload
                ref="uploadRef"
                :auto-upload="false"
                :show-file-list="false"
                accept=".xlsx,.xls"
                @change="handleImport"
              >
                <el-button type="primary" :icon="Download">
                  选择文件导入
                </el-button>
              </el-upload>
            </div>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useThemeStore } from '@/stores/theme'
import { useDataStore } from '@/stores/data'
import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
import {
  Brush,
  Coin,
  Monitor,
  FolderOpened,
  Moon,
  Sunny,
  Check,
  Download,
  RefreshLeft,
  Delete,
  Upload
} from '@element-plus/icons-vue'
import { GetMetaInfo, GetDatabasePath, SelectDatabaseFile, MigrateOldDatabase } from '../../wailsjs/go/main/App'

const themeStore = useThemeStore()
const dataStore = useDataStore()
const uploadRef = ref()

// 数据库路径
const dbPath = ref('')

// 程序信息
const appInfo = ref({
  title: '医疗表单管理系统',
  version: '1.0.0'
})

// 主题颜色选项
const colorOptions = [
  { name: '蓝色', value: '#409EFF' },
  { name: '绿色', value: '#67C23A' },
  { name: '紫色', value: '#9C27B0' },
  { name: '橙色', value: '#FF9800' }
]

// 计算属性
const lastBackupTime = computed(() => {
  if (!dataStore.dbInfo.lastBackup) return '无备份记录'
  const date = new Date(dataStore.dbInfo.lastBackup)
  return date.toLocaleString()
})

// 方法
const formatFileSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const handleThemeChange = async () => {
  try {
    await themeStore.saveTheme()
    ElMessage.success('主题设置已保存')
  } catch (error) {
    ElMessage.error('保存主题设置失败：' + error.message)
  }
}

const handleColorChange = async (color) => {
  try {
    themeStore.setPrimaryColor(color)
    await themeStore.saveTheme()
    ElMessage.success('主题颜色已更新')
  } catch (error) {
    ElMessage.error('保存主题颜色失败：' + error.message)
  }
}

const handleBackup = async () => {
  try {
    await dataStore.backupDatabase(false)
    await dataStore.fetchDatabaseInfo()
    ElMessage.success('数据库备份成功')
  } catch (error) {
    ElMessage.error('数据库备份失败：' + error.message)
  }
}

const handleBackupAndClear = async () => {
  try {
    await ElMessageBox.confirm(
      '此操作将备份当前数据库并清空所有数据，是否继续？',
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await dataStore.backupDatabase(true)
    await dataStore.fetchAllPatients()
    await dataStore.fetchDatabaseInfo()
    ElMessage.success('数据库已备份并清空')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败：' + error.message)
    }
  }
}

const handleClearDatabase = async () => {
  try {
    await ElMessageBox.confirm(
      '此操作将永久删除所有患者数据，且无法恢复，是否继续？',
      '危险操作',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'error',
        confirmButtonClass: 'el-button--danger'
      }
    )
    
    await dataStore.dropDatabase()
    await dataStore.createDatabase()
    await dataStore.fetchAllPatients()
    await dataStore.fetchDatabaseInfo()
    ElMessage.success('数据库已清空')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('清空数据库失败：' + error.message)
    }
  }
}

const handleExport = () => {
  // 这里可以实现数据导出功能
  ElMessage.info('导出功能开发中...')
}

const handleImport = (file) => {
  // 这里可以实现数据导入功能
  ElMessage.info('导入功能开发中...')
}

// 浏览数据库文件
const handleBrowseDatabase = async () => {
  try {
    // 调用后端API选择数据库文件
    const result = await SelectDatabaseFile()
    if (result) {
      dbPath.value = result
      ElMessage.success('数据库路径已更新')
      // 重新加载数据库信息
      await dataStore.fetchDatabaseInfo()
    }
  } catch (error) {
    ElMessage.error('选择数据库文件失败：' + error.message)
  }
}

// 合并旧版本数据库
const handleMigrateOldDatabase = async () => {
  let loadingInstance = null
  try {
    await ElMessageBox.confirm(
      '此操作将选择旧版本数据库文件并合并到当前数据库中，是否继续？',
      '确认合并',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 调用后端API选择旧数据库文件
    const oldDbPath = await SelectDatabaseFile()
    if (oldDbPath) {
      // 显示全局loading
      loadingInstance = ElLoading.service({
        lock: true,
        text: '正在迁移数据库，请稍候...',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      
      // 调用后端迁移方法
      await MigrateOldDatabase(oldDbPath)
      
      // 重新加载数据
      await dataStore.fetchAllPatients()
      await dataStore.fetchDatabaseInfo()
      
      ElMessage.success('数据库迁移成功')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('数据库迁移失败：' + error)
    }
  } finally {
    // 关闭loading
    if (loadingInstance) {
      loadingInstance.close()
    }
  }
}

// 加载程序信息
const loadAppInfo = async () => {
  try {
    const metaInfo = await GetMetaInfo()
    if (metaInfo && metaInfo.title && metaInfo.version) {
      appInfo.value = metaInfo
    }
  } catch (error) {
    console.error('获取程序信息失败：', error)
  }
}

// 初始化
onMounted(async () => {
  try {
    await dataStore.fetchDatabaseInfo()
    await loadAppInfo()
    // 获取当前数据库路径
    const currentPath = await GetDatabasePath()
    if (currentPath) {
      dbPath.value = currentPath
    }
  } catch (error) {
    console.error('加载数据失败：', error)
    ElMessage({
      message: '加载数据库信息失败：' + error.message,
      type: 'error',
      duration: 1500,
    })
  }
})
</script>

<style scoped>
.setting-container {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
  height: calc(100vh - 60px);
  overflow: auto;
}

.setting-section {
  margin-bottom: 24px;
}

.setting-card {
  transition: transform 0.3s;
}

.setting-card:hover {
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-icon {
  color: var(--el-color-primary);
  font-size: 18px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.setting-content {
  padding: 20px 0;
}

.setting-item {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 24px;
}

.setting-item:last-child {
  margin-bottom: 0;
}

.setting-label {
  width: 100%;
  flex: 1;
  margin-right: 20px;
}

.label-text {
  display: block;
  font-size: 16px;
  font-weight: 500;
  color: var(--el-text-color-primary);
  margin-bottom: 4px;
}

.label-desc {
  font-size: 14px;
  color: var(--el-text-color-regular);
  line-height: 1.4;
}

.setting-control {
  flex-shrink: 0;
}

.color-options {
  display: flex;
  gap: 16px;
}

.color-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.color-option:hover {
  transform: translateY(-2px);
}

.color-circle {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 3px solid transparent;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.color-option.active .color-circle {
  border-color: var(--el-color-primary);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.check-icon {
  color: white;
  font-size: 18px;
  font-weight: bold;
}

.color-name {
  font-size: 12px;
  color: var(--el-text-color-regular);
  font-weight: 500;
}

.db-info,
.system-info {
  width: 100%;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 14px;
  color: var(--el-text-color-regular);
  font-weight: 500;
}

.info-value {
  font-size: 14px;
  color: var(--el-text-color-primary);
  font-weight: 600;
}

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.db-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.path-control {
  display: flex;
  gap: 8px;
  align-items: center;
  flex: 1;
}

.path-input {
  flex: 1;
}

.browse-btn {
  flex-shrink: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .setting-container {
    padding: 12px;
  }
  
  .setting-item {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }
  
  .setting-label {
    margin-right: 0;
  }
  
  .color-options {
    justify-content: center;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
  }
  
  .db-actions {
    justify-content: center;
  }
}
</style>