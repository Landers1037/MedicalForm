import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { 
  GetAllPatients, 
  GetPatientByID, 
  AddPatient, 
  UpdatePatient, 
  DeletePatient, 
  SearchPatients,
  GetDatabaseInfo,
  CreateDatabase,
  DropDatabase,
  BackupDatabase
} from '../../wailsjs/go/main/App'

export const useDataStore = defineStore('data', () => {
  // 患者数据列表
  const patientItems = ref([])
  
  // 当前选中的患者
  const currentPatient = ref(null)
  
  // 加载状态
  const loading = ref(false)
  
  // 数据库信息
  const dbInfo = ref({
    size: 0,
    count: 0,
    lastBackup: '',
    backupList: []
  })
  
  // 获取所有患者数据
  const fetchAllPatients = async () => {
    loading.value = true
    try {
      const result = await GetAllPatients()
      patientItems.value = result || []
      return result
    } catch (error) {
      console.error('获取患者数据失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }
  
  // 根据ID获取患者数据
  const fetchPatientById = async (id) => {
    try {
      const result = await GetPatientByID(id)
      currentPatient.value = result
      return result
    } catch (error) {
      console.error('获取患者详情失败:', error)
      throw error
    }
  }
  
  // 添加患者数据
  const addPatient = async (patientData) => {
    try {
      await AddPatient(patientData)
      await fetchAllPatients() // 刷新列表
      return true
    } catch (error) {
      console.error('添加患者数据失败:', error)
      throw error
    }
  }
  
  // 更新患者数据
  const updatePatient = async (patientData) => {
    try {
      await UpdatePatient(patientData)
      await fetchAllPatients() // 刷新列表
      return true
    } catch (error) {
      console.error('更新患者数据失败:', error)
      throw error
    }
  }
  
  // 删除患者数据
  const deletePatient = async (id) => {
    try {
      await DeletePatient(id)
      await fetchAllPatients() // 刷新列表
      return true
    } catch (error) {
      console.error('删除患者数据失败:', error)
      throw error
    }
  }
  
  // 搜索患者数据
  const searchPatients = async (searchType, keyword) => {
    loading.value = true
    try {
      const result = await SearchPatients(searchType, keyword)
      patientItems.value = result || []
      return result
    } catch (error) {
      console.error('搜索患者数据失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }
  
  // 获取数据库信息
  const fetchDatabaseInfo = async () => {
    try {
      const result = await GetDatabaseInfo()
      dbInfo.value = result || { size: 0, count: 0, lastBackup: '', backupList: [] }
      return result
    } catch (error) {
      console.error('获取数据库信息失败:', error)
      throw error
    }
  }
  
  // 创建数据库
  const createDatabase = async () => {
    try {
      await CreateDatabase()
      return true
    } catch (error) {
      console.error('创建数据库失败:', error)
      throw error
    }
  }
  
  // 清空数据库
  const dropDatabase = async () => {
    try {
      await DropDatabase()
      await fetchAllPatients() // 刷新列表
      return true
    } catch (error) {
      console.error('清空数据库失败:', error)
      throw error
    }
  }
  
  // 备份数据库
  const backupDatabase = async (clearAfterBackup = false) => {
    try {
      await BackupDatabase(clearAfterBackup)
      if (clearAfterBackup) {
        await fetchAllPatients() // 如果清空了数据库，刷新列表
      }
      return true
    } catch (error) {
      console.error('备份数据库失败:', error)
      throw error
    }
  }
  
  // 计算总患者数量
  const totalCount = computed(() => {
    return patientItems.value ? patientItems.value.length : 0
  })

  return {
    patientItems,
    currentPatient,
    loading,
    dbInfo,
    totalCount,
    fetchAllPatients,
    fetchPatientById,
    addPatient,
    updatePatient,
    deletePatient,
    searchPatients,
    fetchDatabaseInfo,
    createDatabase,
    dropDatabase,
    backupDatabase
  }
})