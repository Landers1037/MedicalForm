import { defineStore } from 'pinia'
import { ref } from 'vue'
import { 
  GetAllDoctors, 
  GetDoctorByID, 
  AddDoctor, 
  UpdateDoctor, 
  DeleteDoctor, 
  SearchDoctors
} from '../../wailsjs/go/main/App'

export const useDoctorStore = defineStore('doctor', () => {
  // 医师数据列表
  const doctors = ref([])
  
  // 当前选中的医师
  const currentDoctor = ref(null)
  
  // 加载状态
  const loading = ref(false)
  
  // 获取所有医师数据
  const fetchAllDoctors = async () => {
    loading.value = true
    try {
      const result = await GetAllDoctors()
      doctors.value = result || []
      return result
    } catch (error) {
      console.error('获取医师数据失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }
  
  // 根据ID获取医师数据
  const fetchDoctorById = async (id) => {
    try {
      const result = await GetDoctorByID(id)
      currentDoctor.value = result
      return result
    } catch (error) {
      console.error('获取医师详情失败:', error)
      throw error
    }
  }
  
  // 添加医师
  const addDoctor = async (doctorData) => {
    try {
      await AddDoctor(doctorData)
      // 重新获取列表
      await fetchAllDoctors()
      return true
    } catch (error) {
      console.error('添加医师失败:', error)
      throw error
    }
  }
  
  // 更新医师信息
  const updateDoctor = async (doctorData) => {
    try {
      await UpdateDoctor(doctorData)
      // 重新获取列表
      await fetchAllDoctors()
      return true
    } catch (error) {
      console.error('更新医师失败:', error)
      throw error
    }
  }
  
  // 删除医师
  const deleteDoctor = async (id) => {
    try {
      await DeleteDoctor(id)
      // 从本地列表中移除
      doctors.value = doctors.value.filter(doctor => doctor.id !== id)
      // 如果删除的是当前选中的医师，清空选中状态
      if (currentDoctor.value && currentDoctor.value.id === id) {
        currentDoctor.value = null
      }
      return true
    } catch (error) {
      console.error('删除医师失败:', error)
      throw error
    }
  }
  
  // 搜索医师
  const searchDoctors = async (searchType, keyword) => {
    loading.value = true
    try {
      const result = await SearchDoctors(searchType, keyword)
      return result || []
    } catch (error) {
      console.error('搜索医师失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }
  
  // 根据级别获取医师列表
  const getDoctorsByLevel = (level) => {
    return doctors.value.filter(doctor => doctor.level === level)
  }
  
  // 获取所有医师级别
  const getAllLevels = () => {
    const levels = [...new Set(doctors.value.map(doctor => doctor.level))]
    return levels.filter(level => level && level.trim())
  }
  
  // 获取医师统计信息
  const getDoctorStats = () => {
    const total = doctors.value.length
    const maleCount = doctors.value.filter(doctor => doctor.sex === '男').length
    const femaleCount = doctors.value.filter(doctor => doctor.sex === '女').length
    
    const levelStats = {}
    doctors.value.forEach(doctor => {
      if (doctor.level) {
        levelStats[doctor.level] = (levelStats[doctor.level] || 0) + 1
      }
    })
    
    return {
      total,
      maleCount,
      femaleCount,
      levelStats
    }
  }
  
  // 清空当前选中的医师
  const clearCurrentDoctor = () => {
    currentDoctor.value = null
  }
  
  // 设置当前选中的医师
  const setCurrentDoctor = (doctor) => {
    currentDoctor.value = doctor
  }
  
  return {
    doctors,
    currentDoctor,
    loading,
    fetchAllDoctors,
    fetchDoctorById,
    addDoctor,
    updateDoctor,
    deleteDoctor,
    searchDoctors,
    getDoctorsByLevel,
    getAllLevels,
    getDoctorStats,
    clearCurrentDoctor,
    setCurrentDoctor
  }
})