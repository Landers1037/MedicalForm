<template>
  <div class="home-container">
    <!-- 欢迎区域 -->
    <div class="welcome-section">
      <el-card class="welcome-card" shadow="hover">
        <div class="welcome-content">
          <div class="welcome-text">
            <h1 class="welcome-title">欢迎使用 {{ appTitle }}</h1>
            <p class="welcome-subtitle">高效管理患者信息，提升医疗服务质量</p>
          </div>
          <div class="welcome-icon">
            <el-icon :size="80" color="var(--el-color-primary)">
              <Document />
            </el-icon>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 统计信息 -->
    <div class="stats-section">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-card class="stat-card" shadow="hover">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon :size="32" color="#409EFF">
                  <User />
                </el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ dataStore.totalCount }}</div>
                <div class="stat-label">总患者数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card" shadow="hover">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon :size="32" color="#67C23A">
                  <Calendar />
                </el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ todayCount }}</div>
                <div class="stat-label">今日新增</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card" shadow="hover">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon :size="32" color="#E6A23C">
                  <FolderOpened />
                </el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ formatFileSize(dataStore.dbInfo.size) }}</div>
                <div class="stat-label">数据库大小</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card" shadow="hover">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon :size="32" color="#F56C6C">
                  <Clock />
                </el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ lastBackupTime }}</div>
                <div class="stat-label">最后备份</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 收入统计 -->
    <div class="revenue-section">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card class="revenue-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <span class="card-title">收入明细</span>
                <el-button type="primary" link @click="showRevenueDialog = true">
                  <el-icon><TrendCharts /></el-icon>
                  查看详情
                </el-button>
              </div>
            </template>
            <div class="revenue-content">
              <div class="revenue-icon">
                <el-icon :size="48" color="var(--el-color-primary)">
                  <TrendCharts />
                </el-icon>
              </div>
              <div class="revenue-info">
                <div class="revenue-text">点击查看收入统计图表</div>
                <div class="revenue-desc">年度、月度、周度收入分析</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="revenue-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <span class="card-title">今日收入</span>
              </div>
            </template>
            <div class="revenue-content">
              <div class="revenue-icon">
                <el-icon :size="48" color="var(--el-color-success)">
                  <Money />
                </el-icon>
              </div>
              <div class="revenue-info">
                <div class="revenue-number">¥{{ todayRevenue.toFixed(2) }}</div>
                <div class="revenue-desc">今日总收入</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 快速操作 -->
    <div class="actions-section">
      <el-card class="actions-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span class="card-title">快速操作</span>
          </div>
        </template>
        <div class="actions-content">
          <el-row :gutter="20">
            <el-col :span="6">
              <div class="action-item" @click="$router.push('/panel')">
                <el-icon :size="40" color="var(--el-color-primary)">
                  <Plus />
                </el-icon>
                <span class="action-text">添加患者</span>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="action-item" @click="$router.push('/panel')">
                <el-icon :size="40" color="var(--el-color-success)">
                  <Search />
                </el-icon>
                <span class="action-text">查询患者</span>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="action-item" @click="handleBackup">
                <el-icon :size="40" color="var(--el-color-warning)">
                  <Download />
                </el-icon>
                <span class="action-text">备份数据</span>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="action-item" @click="$router.push('/setting')">
                <el-icon :size="40" color="var(--el-color-info)">
                  <Setting />
                </el-icon>
                <span class="action-text">系统设置</span>
              </div>
            </el-col>
          </el-row>
        </div>
      </el-card>
    </div>

    <!-- 最近记录 -->
    <div class="recent-section">
      <el-card class="recent-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span class="card-title">最近记录</span>
            <el-button type="primary" link @click="$router.push('/panel')">
              查看全部
            </el-button>
          </div>
        </template>
        <div class="recent-content">
          <el-table :data="recentData" style="width: 100%" stripe>
            <el-table-column prop="name" label="姓名" width="120" />
            <el-table-column prop="sex" label="性别" width="80">
              <template #default="{ row }">
                <el-tag :type="row.sex === '男' ? 'primary' : 'danger'" size="small">
                  {{ row.sex }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="age" label="年龄" width="80" />
            <el-table-column prop="phone" label="电话" width="140" />
            <el-table-column prop="date" label="登记日期" width="120" />
            <el-table-column prop="detail" label="诊断" show-overflow-tooltip />
          </el-table>
        </div>
      </el-card>
    </div>

    <!-- 收入明细弹框 -->
    <el-dialog
      v-model="showRevenueDialog"
      title="收入明细统计"
      width="80%"
      :before-close="() => showRevenueDialog = false"
    >
      <div class="revenue-dialog-content">
        <div class="revenue-tabs">
          <el-radio-group v-model="revenueChartType" class="chart-type-selector">
            <el-radio-button label="year">年度统计</el-radio-button>
            <el-radio-button label="month">月度统计</el-radio-button>
            <el-radio-button label="week">周度统计</el-radio-button>
          </el-radio-group>
        </div>
        <div class="revenue-chart-container">
          <div id="revenueChart" style="width: 100%; height: 400px;"></div>
        </div>
        <div class="revenue-summary">
          <el-row :gutter="20">
            <el-col :span="8">
              <el-statistic title="总收入" :value="totalRevenue" suffix="元" />
            </el-col>
            <el-col :span="8">
              <el-statistic title="平均收入" :value="averageRevenue" suffix="元" :precision="2" />
            </el-col>
            <el-col :span="8">
              <el-statistic title="统计周期" :value="chartPeriodText" />
            </el-col>
          </el-row>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useDataStore } from '@/stores/data'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import {
  Document,
  User,
  Calendar,
  FolderOpened,
  Clock,
  Plus,
  Search,
  Download,
  Setting,
  TrendCharts,
  Money
} from '@element-plus/icons-vue'
import { GetMetaInfo } from '../../wailsjs/go/main/App'

// 程序标题
const appTitle = ref('医疗表单管理系统')

// 加载程序信息
const loadAppInfo = async () => {
  try {
    const metaInfo = await GetMetaInfo()
    if (metaInfo && metaInfo.title) {
      appTitle.value = metaInfo.title
    }
  } catch (error) {
    console.error('获取程序信息失败:', error)
  }
}

const dataStore = useDataStore()
const recentData = ref([])

// 收入统计相关
const showRevenueDialog = ref(false)
const revenueChartType = ref('year') // year, month, week
let revenueChart = null

// 计算今日新增数量
const todayCount = computed(() => {
  if (!dataStore.patientItems || !Array.isArray(dataStore.patientItems)) return 0
  const today = new Date().toISOString().split('T')[0]
  return dataStore.patientItems.filter(p => p.date === today).length
})

// 计算今日收入
const todayRevenue = computed(() => {
  if (!dataStore.patientItems || !Array.isArray(dataStore.patientItems)) return 0
  const today = new Date().toISOString().split('T')[0]
  return dataStore.patientItems
    .filter(p => p.date === today)
    .reduce((total, patient) => {
      const fee = parseFloat(patient.fee) || 0
      return total + fee
    }, 0)
})

// 计算收入图表数据
const revenueChartData = computed(() => {
  if (!dataStore.patientItems || !Array.isArray(dataStore.patientItems)) return { labels: [], data: [] }
  
  const now = new Date()
  const data = {}
  
  if (revenueChartType.value === 'year') {
    // 年度统计 - 最近12个月
    for (let i = 11; i >= 0; i--) {
      const date = new Date(now.getFullYear(), now.getMonth() - i, 1)
      const key = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
      data[key] = 0
    }
    
    dataStore.patientItems.forEach(patient => {
      const patientDate = new Date(patient.date)
      const key = `${patientDate.getFullYear()}-${String(patientDate.getMonth() + 1).padStart(2, '0')}`
      if (data.hasOwnProperty(key)) {
        data[key] += parseFloat(patient.fee) || 0
      }
    })
  } else if (revenueChartType.value === 'month') {
    // 月度统计 - 最近30天
    for (let i = 29; i >= 0; i--) {
      const date = new Date(now.getTime() - i * 24 * 60 * 60 * 1000)
      const key = date.toISOString().split('T')[0]
      data[key] = 0
    }
    
    dataStore.patientItems.forEach(patient => {
      if (data.hasOwnProperty(patient.date)) {
        data[patient.date] += parseFloat(patient.fee) || 0
      }
    })
  } else if (revenueChartType.value === 'week') {
    // 周度统计 - 最近12周
    for (let i = 11; i >= 0; i--) {
      const startOfWeek = new Date(now.getTime() - i * 7 * 24 * 60 * 60 * 1000)
      startOfWeek.setDate(startOfWeek.getDate() - startOfWeek.getDay())
      const key = `第${12-i}周`
      data[key] = 0
    }
    
    dataStore.patientItems.forEach(patient => {
      const patientDate = new Date(patient.date)
      const weeksDiff = Math.floor((now.getTime() - patientDate.getTime()) / (7 * 24 * 60 * 60 * 1000))
      if (weeksDiff >= 0 && weeksDiff < 12) {
        const key = `第${12-weeksDiff}周`
        data[key] += parseFloat(patient.fee) || 0
      }
    })
  }
  
  return {
    labels: Object.keys(data),
    data: Object.values(data)
  }
})

// 计算总收入
const totalRevenue = computed(() => {
  return revenueChartData.value.data.reduce((sum, value) => sum + value, 0)
})

// 计算平均收入
const averageRevenue = computed(() => {
  const data = revenueChartData.value.data
  return data.length > 0 ? totalRevenue.value / data.length : 0
})

// 图表周期文本
const chartPeriodText = computed(() => {
  switch (revenueChartType.value) {
    case 'year': return '最近12个月'
    case 'month': return '最近30天'
    case 'week': return '最近12周'
    default: return ''
  }
})

// 格式化最后备份时间
const lastBackupTime = computed(() => {
  if (!dataStore.dbInfo.lastBackup) return '无'
  const date = new Date(dataStore.dbInfo.lastBackup)
  return date.toLocaleDateString()
})

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 处理备份操作
const handleBackup = async () => {
  try {
    await dataStore.backupDatabase(false)
    ElMessage.success('数据备份成功')
  } catch (error) {
    ElMessage.error('数据备份失败：' + error.message)
  }
}

// 初始化收入图表
const initRevenueChart = () => {
  const chartDom = document.getElementById('revenueChart')
  if (!chartDom) return
  
  revenueChart = echarts.init(chartDom)
  updateRevenueChart()
}

// 更新收入图表
const updateRevenueChart = () => {
  if (!revenueChart) return
  
  const { labels, data } = revenueChartData.value
  
  const option = {
    title: {
      text: `收入统计 - ${chartPeriodText.value}`,
      left: 'center'
    },
    tooltip: {
      trigger: 'axis',
      formatter: function(params) {
        return `${params[0].name}<br/>收入: ¥${params[0].value.toFixed(2)}`
      }
    },
    xAxis: {
      type: 'category',
      data: labels,
      axisLabel: {
        rotate: revenueChartType.value === 'month' ? 45 : 0
      }
    },
    yAxis: {
      type: 'value',
      name: '收入(元)',
      axisLabel: {
        formatter: '¥{value}'
      }
    },
    series: [{
      name: '收入',
      type: 'bar',
      data: data,
      itemStyle: {
        color: '#409EFF'
      },
      emphasis: {
        itemStyle: {
          color: '#66b1ff'
        }
      }
    }]
  }
  
  revenueChart.setOption(option)
}

// 组件是否已卸载
const isUnmounted = ref(false)

// 初始化数据
onMounted(async () => {
  try {
    await loadAppInfo();
    await dataStore.fetchAllPatients()
    await dataStore.fetchDatabaseInfo()
    
    // 检查组件是否已卸载
    if (isUnmounted.value) return
    
    // 获取最近5条记录
    if (dataStore.patientItems && Array.isArray(dataStore.patientItems)) {
      recentData.value = dataStore.patientItems
        .sort((a, b) => new Date(b.date) - new Date(a.date))
        .slice(0, 5)
    }
  } catch (error) {
    if (!isUnmounted.value) {
      console.error('加载数据失败：', error)
      ElMessage({
        message: '加载数据失败：' + error.message,
        type: 'error',
        duration: 1500,
      })
    }
  }
})

// 监听图表类型变化
watch(revenueChartType, () => {
  updateRevenueChart()
})

// 监听弹框打开状态
watch(showRevenueDialog, async (newVal) => {
  if (newVal) {
    await nextTick()
    initRevenueChart()
  }
})

// 组件卸载时清理
onUnmounted(() => {
  isUnmounted.value = true
  if (revenueChart) {
    revenueChart.dispose()
    revenueChart = null
  }
})
</script>

<style scoped>
.home-container {
  padding: 16px;
  max-width: 1200px;
  margin: 0 auto;
  user-select: none;
}

.welcome-section {
  margin-bottom: 24px;
}

.welcome-card {
  background: linear-gradient(135deg, var(--el-color-primary-light-9) 0%, var(--el-color-primary-light-8) 100%);
  border: none;
}

.welcome-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px;
}

.welcome-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 8px 0;
}

.welcome-subtitle {
  font-size: 16px;
  color: var(--el-text-color-regular);
  margin: 0;
}

.stats-section {
  margin-bottom: 24px;
}

.revenue-section {
  margin-bottom: 24px;
}

.revenue-card {
  height: 180px;
  transition: transform 0.3s;
}

.revenue-card:hover {
  transform: translateY(-4px);
}

.revenue-content {
  display: flex;
  align-items: center;
  height: 100%;
  padding: 20px;
}

.revenue-icon {
  margin-right: 20px;
}

.revenue-info {
  flex: 1;
}

.revenue-number {
  font-size: 28px;
  font-weight: 600;
  color: var(--el-color-success);
  line-height: 1;
  margin-bottom: 8px;
}

.revenue-text {
  font-size: 16px;
  color: var(--el-text-color-primary);
  margin-bottom: 8px;
}

.revenue-desc {
  font-size: 14px;
  color: var(--el-text-color-regular);
}

.stat-card {
  height: 120px;
  transition: transform 0.3s;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-content {
  display: flex;
  align-items: center;
  height: 100%;
  padding: 20px;
}

.stat-icon {
  margin-right: 16px;
}

.stat-number {
  font-size: 24px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  line-height: 1;
}

.stat-label {
  font-size: 14px;
  color: var(--el-text-color-regular);
  margin-top: 4px;
}

.actions-section {
  margin-bottom: 24px;
}

.actions-card {
  min-height: 200px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.actions-content {
  padding: 20px 0;
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  border-radius: var(--el-border-radius-base);
  cursor: pointer;
  transition: all 0.3s;
  background-color: var(--el-fill-color-lighter);
}

.action-item:hover {
  background-color: var(--el-fill-color-light);
  transform: translateY(-2px);
  box-shadow: var(--el-box-shadow-light);
}

.action-text {
  margin-top: 12px;
  font-size: 14px;
  font-weight: 500;
  color: var(--el-text-color-primary);
}

.recent-section {
  margin-bottom: 24px;
}

.recent-card {
  min-height: 300px;
}

.recent-content {
  padding-top: 10px;
}

/* 收入明细弹框样式 */
.revenue-dialog-content {
  padding: 20px 0;
}

.revenue-tabs {
  margin-bottom: 20px;
  text-align: center;
}

.chart-type-selector {
  display: inline-block;
}

.revenue-chart-container {
  margin-bottom: 30px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: var(--el-border-radius-base);
  padding: 20px;
  background-color: var(--el-bg-color-page);
}

.revenue-summary {
  padding: 20px;
  background-color: var(--el-fill-color-lighter);
  border-radius: var(--el-border-radius-base);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .home-container {
    padding: 12px;
  }
  
  .welcome-content {
    flex-direction: column;
    text-align: center;
    gap: 20px;
  }
  
  .welcome-title {
    font-size: 24px;
  }
  
  .stat-content {
    padding: 16px;
  }
  
  .actions-content {
    padding: 16px 0;
  }
  
  .action-item {
    padding: 16px;
  }
}
</style>