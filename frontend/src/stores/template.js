import { defineStore } from 'pinia'
import { ref } from 'vue'
import { GetTemplatesByType } from '../../wailsjs/go/main/App'

/**
 * 模板管理Store
 * 用于管理诊断、医嘱、治疗方案模板数据
 */
export const useTemplateStore = defineStore('template', () => {
  // 诊断模板列表
  const diagnosisTemplates = ref([])
  // 医嘱模板列表
  const adviceTemplates = ref([])
  // 治疗方案模板列表
  const treatmentTemplates = ref([])
  // 加载状态
  const loading = ref(false)

  /**
   * 获取诊断模板
   */
  const fetchDiagnosisTemplates = async () => {
    try {
      loading.value = true
      const templates = await GetTemplatesByType('diagnosis')
      diagnosisTemplates.value = templates || []
    } catch (error) {
      console.error('获取诊断模板失败:', error)
      diagnosisTemplates.value = []
    } finally {
      loading.value = false
    }
  }

  /**
   * 获取医嘱模板
   */
  const fetchAdviceTemplates = async () => {
    try {
      loading.value = true
      const templates = await GetTemplatesByType('advice')
      adviceTemplates.value = templates || []
    } catch (error) {
      console.error('获取医嘱模板失败:', error)
      adviceTemplates.value = []
    } finally {
      loading.value = false
    }
  }

  /**
   * 获取治疗方案模板
   */
  const fetchTreatmentTemplates = async () => {
    try {
      loading.value = true
      const templates = await GetTemplatesByType('treatment')
      treatmentTemplates.value = templates || []
    } catch (error) {
      console.error('获取治疗方案模板失败:', error)
      treatmentTemplates.value = []
    } finally {
      loading.value = false
    }
  }

  /**
   * 初始化所有模板数据
   */
  const initTemplates = async () => {
    await Promise.all([
      fetchDiagnosisTemplates(),
      fetchAdviceTemplates(),
      fetchTreatmentTemplates()
    ])
  }

  /**
   * 根据ID获取模板内容
   * @param {string} type - 模板类型
   * @param {number} id - 模板ID
   * @returns {string} 模板内容
   */
  const getTemplateContent = (type, id) => {
    let templates = []
    switch (type) {
      case 'diagnosis':
        templates = diagnosisTemplates.value
        break
      case 'advice':
        templates = adviceTemplates.value
        break
      case 'treatment':
        templates = treatmentTemplates.value
        break
    }
    const template = templates.find(t => t.id === id)
    return template ? template.content : ''
  }

  return {
    // 状态
    diagnosisTemplates,
    adviceTemplates,
    treatmentTemplates,
    loading,
    // 方法
    fetchDiagnosisTemplates,
    fetchAdviceTemplates,
    fetchTreatmentTemplates,
    initTemplates,
    getTemplateContent
  }
})