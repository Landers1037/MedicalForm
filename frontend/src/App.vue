<template>
  <el-container class="app-container">
    <!-- 顶部导航栏 -->
    <el-header class="app-header">
      <div class="header-content">
        <div class="logo-section">
          <el-icon class="logo-icon" :size="24">
            <Document />
          </el-icon>
          <span class="app-title">{{ appTitle }}</span>
        </div>
        
        <el-menu
          :default-active="$route.path"
          mode="horizontal"
          class="nav-menu"
          @select="handleMenuSelect"
        >
          <el-menu-item index="/">
            <el-icon><House /></el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="/panel">
            <el-icon><Grid /></el-icon>
            <span>数据管理</span>
          </el-menu-item>
          <el-menu-item index="/doctor">
            <el-icon><User /></el-icon>
            <span>医师管理</span>
          </el-menu-item>
          <el-menu-item index="/setting">
            <el-icon><Setting /></el-icon>
            <span>设置</span>
          </el-menu-item>
          <el-menu-item index="/about">
            <el-icon><InfoFilled /></el-icon>
            <span>关于</span>
          </el-menu-item>
        </el-menu>
        
        <div class="header-actions">
          <!-- 主题切换按钮 -->
          <el-tooltip :content="themeStore.isDark ? '切换到明亮模式' : '切换到暗黑模式'">
            <el-button
              circle
              @click="themeStore.toggleMode()"
              :icon="themeStore.isDark ? Sunny : Moon"
            />
          </el-tooltip>
        </div>
      </div>
    </el-header>
    
    <!-- 主内容区域 -->
    <el-main class="app-main">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </el-main>
  </el-container>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useThemeStore } from './stores/theme'
import {
  Document,
  House,
  Grid,
  User,
  Setting,
  InfoFilled,
  Sunny,
  Moon
} from '@element-plus/icons-vue'
import { GetMetaInfo } from '../wailsjs/go/main/App'

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

const router = useRouter()
const themeStore = useThemeStore()

// 处理菜单选择
const handleMenuSelect = async (index) => {
  // 防止重复导航
  if (router.currentRoute.value.path === index) {
    return
  }
  
  try {
    // 使用nextTick确保DOM更新完成
    await nextTick()
    await router.push(index)
  } catch (error) {
    // 忽略导航相关错误
    if (error.name !== 'NavigationDuplicated' && error.name !== 'NavigationCancelled') {
      console.error('路由跳转失败:', error)
    }
  }
}

// 初始化主题和程序信息
onMounted(async () => {
  await themeStore.initTheme()
  await loadAppInfo()
})
</script>

<style scoped>
.app-container {
  height: 100vh;
  width: 100vw;
}

.app-header {
  background-color: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);
  padding: 0;
  height: 60px;
  box-shadow: var(--el-box-shadow-light);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  padding: 0 20px;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  color: var(--el-color-primary);
}

.app-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.nav-menu {
  flex: 1;
  margin: 0 40px;
  border-bottom: none;
  background-color: transparent;
}

.nav-menu .el-menu-item {
  border-radius: var(--el-border-radius-small);
  margin: 0 4px;
  transition: all 0.3s;
}

.nav-menu .el-menu-item:hover {
  background-color: var(--el-fill-color-light);
}

.nav-menu .el-menu-item.is-active {
  background-color: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.app-main {
  padding: 0;
  background-color: var(--el-bg-color-page);
  overflow: auto;
}

/* 路由过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
