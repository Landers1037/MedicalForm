<template>
  <div class="about-container">
    <!-- 系统介绍 -->
    <div class="intro-section">
      <el-card class="intro-card" shadow="hover">
        <div class="intro-content">
          <div class="logo-section">
            <el-icon class="app-logo" :size="80">
              <Document />
            </el-icon>
            <h1 class="app-name">{{ appInfo.title }}</h1>
            <p class="app-version">Version {{ appInfo.version }}</p>
          </div>
          <div class="description">
            <p class="desc-text">
              {{ appInfo.title }}是一款基于现代化技术栈开发的患者信息管理软件，
              旨在为医疗机构提供高效、安全、易用的患者数据管理解决方案。
            </p>
            <p class="desc-text">
              系统采用Go语言作为后端，Vue3作为前端框架，通过Wails技术实现桌面应用，
              为用户提供流畅的操作体验和可靠的数据存储。
            </p>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 功能特性 -->
    <div class="features-section">
      <el-card class="features-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <el-icon class="header-icon">
              <Star />
            </el-icon>
            <span class="card-title">核心功能</span>
          </div>
        </template>
        <div class="features-content">
          <el-row :gutter="20">
            <el-col :span="12" v-for="feature in features" :key="feature.title">
              <div class="feature-item">
                <div class="feature-icon">
                  <el-icon :size="24" :color="feature.color">
                    <component :is="feature.icon" />
                  </el-icon>
                </div>
                <div class="feature-info">
                  <h3 class="feature-title">{{ feature.title }}</h3>
                  <p class="feature-desc">{{ feature.description }}</p>
                </div>
              </div>
            </el-col>
          </el-row>
        </div>
      </el-card>
    </div>

    <!-- 技术栈 -->
    <div class="tech-section">
      <el-card class="tech-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <el-icon class="header-icon">
              <Cpu />
            </el-icon>
            <span class="card-title">技术栈</span>
          </div>
        </template>
        <div class="tech-content">
          <div class="tech-category">
            <h3 class="category-title">前端技术</h3>
            <div class="tech-items">
              <div class="tech-item" v-for="tech in frontendTech" :key="tech.name">
                <div class="tech-icon" :style="{ backgroundColor: tech.color }">
                  <span class="tech-letter">{{ tech.name.charAt(0) }}</span>
                </div>
                <div class="tech-info">
                  <span class="tech-name">{{ tech.name }}</span>
                  <span class="tech-desc">{{ tech.description }}</span>
                </div>
              </div>
            </div>
          </div>
          
          <el-divider />
          
          <div class="tech-category">
            <h3 class="category-title">后端技术</h3>
            <div class="tech-items">
              <div class="tech-item" v-for="tech in backendTech" :key="tech.name">
                <div class="tech-icon" :style="{ backgroundColor: tech.color }">
                  <span class="tech-letter">{{ tech.name.charAt(0) }}</span>
                </div>
                <div class="tech-info">
                  <span class="tech-name">{{ tech.name }}</span>
                  <span class="tech-desc">{{ tech.description }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 开发团队 -->
    <div class="team-section">
      <el-card class="team-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <el-icon class="header-icon">
              <UserFilled />
            </el-icon>
            <span class="card-title">开发团队</span>
          </div>
        </template>
        <div class="team-content">
          <div class="team-info">
            <p class="team-desc">
              本系统致力于为医疗机构提供优质的信息化解决方案。
            </p>
          </div>
          <div class="contact-info">
            <h3 class="contact-title">联系我们</h3>
            <div class="contact-items">
              <div class="contact-item">
                <el-icon class="contact-icon">
                  <Message />
                </el-icon>
                <span class="contact-text">邮箱：mail@renj.io</span>
              </div>
              <div class="contact-item">
                <el-icon class="contact-icon">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                  </svg>
                </el-icon>
                <span class="contact-text">Github: <a href="https://github.com/landers1037" target="_blank" rel="noopener noreferrer" style="color: var(--el-color-primary); text-decoration: none;">https://github.com/landers1037</a></span>
              </div>
              <div class="contact-item">
                <el-icon class="contact-icon">
                  <Location />
                </el-icon>
                <span class="contact-text">地址：中国</span>
              </div>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 版权信息 -->
    <div class="copyright-section">
      <el-card class="copyright-card" shadow="never">
        <div class="copyright-content">
          <p class="copyright-text">
            © 2025 {{ appInfo.title }}. All rights reserved.
          </p>
          <p class="copyright-desc">
            本软件受著作权法保护，未经授权不得复制、分发或修改。
          </p>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  Document,
  Star,
  Cpu,
  UserFilled,
  Message,
  Phone,
  Location,
  User,
  Search,
  Setting,
  PieChart,
  Lock,
  Lightning
} from '@element-plus/icons-vue'
import { GetMetaInfo } from '../../wailsjs/go/main/App'

// 程序信息
const appInfo = ref({
  title: '医疗表单管理系统',
  version: '1.0.0'
})

// 加载程序信息
const loadAppInfo = async () => {
  try {
    const metaInfo = await GetMetaInfo()
    if (metaInfo && metaInfo.title && metaInfo.version) {
      appInfo.value = metaInfo
    }
  } catch (error) {
    console.error('获取程序信息失败:', error)
  }
}

// 组件挂载时加载信息
onMounted(() => {
  loadAppInfo()
})

// 功能特性
const features = [
  {
    title: '患者管理',
    description: '完整的患者信息录入、编辑、查询和删除功能',
    icon: User,
    color: '#409EFF'
  },
  {
    title: '智能搜索',
    description: '支持多条件搜索，快速定位患者信息',
    icon: Search,
    color: '#67C23A'
  },
  {
    title: '数据统计',
    description: '提供丰富的数据统计和分析功能',
    icon: PieChart,
    color: '#E6A23C'
  },
  {
    title: '系统设置',
    description: '灵活的主题配置和系统参数设置',
    icon: Setting,
    color: '#F56C6C'
  },
  {
    title: '数据安全',
    description: '本地数据存储，确保患者信息安全',
    icon: Lock,
    color: '#9C27B0'
  },
  {
    title: '高性能',
    description: '基于现代化技术栈，响应速度快',
    icon: Lightning,
    color: '#FF9800'
  }
]

// 前端技术栈
const frontendTech = [
  {
    name: 'Vue 3',
    description: '渐进式JavaScript框架',
    color: '#4FC08D'
  },
  {
    name: 'Element Plus',
    description: '基于Vue 3的组件库',
    color: '#409EFF'
  },
  {
    name: 'Vite',
    description: '下一代前端构建工具',
    color: '#646CFF'
  },
  {
    name: 'Pinia',
    description: 'Vue的状态管理库',
    color: '#FFD43B'
  }
]

// 后端技术栈
const backendTech = [
  {
    name: 'Go',
    description: '高性能编程语言',
    color: '#00ADD8'
  },
  {
    name: 'Wails',
    description: '使用Go构建桌面应用',
    color: '#FF6B6B'
  },
  {
    name: 'GORM',
    description: 'Go语言ORM库',
    color: '#00D4AA'
  },
  {
    name: 'SQLite',
    description: '轻量级数据库',
    color: '#003B57'
  }
]
</script>

<style scoped>
.about-container {
  padding: 16px;
  max-width: 1000px;
  margin: 0 auto;
  height: calc(100vh - 60px);
  overflow: auto;
}

.intro-section,
.features-section,
.tech-section,
.team-section,
.copyright-section {
  margin-bottom: 24px;
}

.intro-card,
.features-card,
.tech-card,
.team-card {
  transition: transform 0.3s;
}

.intro-card:hover,
.features-card:hover,
.tech-card:hover,
.team-card:hover {
  transform: translateY(-4px);
}

.intro-content {
  text-align: center;
  padding: 40px 20px;
}

.logo-section {
  margin-bottom: 32px;
}

.app-logo {
  color: var(--el-color-primary);
  margin-bottom: 16px;
}

.app-name {
  font-size: 32px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 8px 0;
}

.app-version {
  font-size: 16px;
  color: var(--el-text-color-regular);
  margin: 0;
}

.description {
  max-width: 600px;
  margin: 0 auto;
}

.desc-text {
  font-size: 16px;
  line-height: 1.6;
  color: var(--el-text-color-regular);
  margin-bottom: 16px;
  text-align: left;
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
  font-size: 18px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.features-content {
  padding: 20px 0;
}

.feature-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 20px;
  border-radius: var(--el-border-radius-base);
  background-color: var(--el-fill-color-lighter);
  margin-bottom: 16px;
  transition: all 0.3s;
}

.feature-item:hover {
  background-color: var(--el-fill-color-light);
  transform: translateY(-2px);
}

.feature-icon {
  flex-shrink: 0;
}

.feature-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 8px 0;
}

.feature-desc {
  font-size: 14px;
  color: var(--el-text-color-regular);
  line-height: 1.5;
  margin: 0;
}

.tech-content {
  padding: 20px 0;
}

.tech-category {
  margin-bottom: 24px;
}

.category-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 16px 0;
}

.tech-items {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 16px;
}

.tech-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border-radius: var(--el-border-radius-base);
  background-color: var(--el-fill-color-lighter);
  transition: all 0.3s;
}

.tech-item:hover {
  background-color: var(--el-fill-color-light);
  transform: translateY(-2px);
}

.tech-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.tech-letter {
  color: white;
  font-size: 16px;
  font-weight: 600;
}

.tech-info {
  display: flex;
  flex-direction: column;
}

.tech-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin-bottom: 4px;
}

.tech-desc {
  font-size: 14px;
  color: var(--el-text-color-regular);
}

.team-content {
  padding: 20px 0;
}

.team-info {
  margin-bottom: 32px;
}

.team-desc {
  font-size: 16px;
  line-height: 1.6;
  color: var(--el-text-color-regular);
  margin: 0;
}

.contact-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 16px 0;
}

.contact-items {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.contact-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.contact-icon {
  color: var(--el-color-primary);
  font-size: 16px;
}

.contact-text {
  font-size: 14px;
  color: var(--el-text-color-regular);
}

.copyright-card {
  background-color: var(--el-fill-color-lighter);
  border: none;
}

.copyright-content {
  text-align: center;
  padding: 20px;
}

.copyright-text {
  font-size: 14px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 8px 0;
}

.copyright-desc {
  font-size: 12px;
  color: var(--el-text-color-regular);
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .about-container {
    padding: 12px;
  }
  
  .intro-content {
    padding: 20px 12px;
  }
  
  .app-name {
    font-size: 24px;
  }
  
  .features-content .el-col {
    margin-bottom: 16px;
  }
  
  .tech-items {
    grid-template-columns: 1fr;
  }
  
  .contact-items {
    align-items: flex-start;
  }
}
</style>