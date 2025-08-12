import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { SetConfig, GetConfig } from '../../wailsjs/go/main/App'

export const useThemeStore = defineStore('theme', () => {
  // 主题模式：light | dark
  const mode = ref('light')
  
  // 主题颜色：blue | green | purple | orange
  const color = ref('blue')
  
  // 主题颜色配置
  const themeColors = {
    blue: '#409EFF',
    green: '#67C23A', 
    purple: '#9C27B0',
    orange: '#FF9800'
  }
  
  // 计算当前主题颜色值
  const currentColor = computed(() => themeColors[color.value])
  
  // 是否为暗黑模式
  const isDark = computed(() => mode.value === 'dark')
  
  // 设置主题模式
  const setMode = async (newMode) => {
    mode.value = newMode
    await saveThemeConfig()
    applyTheme()
  }
  
  // 设置主题颜色
  const setColor = async (newColor) => {
    color.value = newColor
    await saveThemeConfig()
    applyTheme()
  }
  
  // 切换主题模式
  const toggleMode = async () => {
    const newMode = mode.value === 'light' ? 'dark' : 'light'
    await setMode(newMode)
  }
  
  // 应用主题到DOM
  const applyTheme = () => {
    const html = document.documentElement
    
    // 设置暗黑模式类名
    if (isDark.value) {
      html.classList.add('dark')
    } else {
      html.classList.remove('dark')
    }
    
    // 设置CSS变量
    html.style.setProperty('--el-color-primary', currentColor.value)
    
    // 设置Element Plus主题
    if (isDark.value) {
      html.setAttribute('data-theme', 'dark')
    } else {
      html.removeAttribute('data-theme')
    }
  }
  
  // 保存主题配置到后端
  const saveThemeConfig = async () => {
    try {
      await SetConfig('theme_mode', mode.value)
      await SetConfig('theme_color', color.value)
    } catch (error) {
      console.error('保存主题配置失败:', error)
    }
  }
  
  // 从后端加载主题配置
  const loadThemeConfig = async () => {
    try {
      const savedMode = await GetConfig('theme_mode')
      const savedColor = await GetConfig('theme_color')
      
      if (savedMode) {
        mode.value = savedMode
      }
      
      if (savedColor && themeColors[savedColor]) {
        color.value = savedColor
      }
      
      applyTheme()
    } catch (error) {
      console.error('加载主题配置失败:', error)
      // 使用默认配置
      applyTheme()
    }
  }
  
  // 初始化主题
  const initTheme = async () => {
    await loadThemeConfig()
  }
  
  return {
    mode,
    color,
    themeColors,
    currentColor,
    isDark,
    setMode,
    setColor,
    toggleMode,
    applyTheme,
    loadThemeConfig,
    initTheme
  }
})