import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: {
      title: '首页'
    }
  },
  {
    path: '/panel',
    name: 'Panel',
    component: () => import('@/views/Panel.vue'),
    meta: {
      title: '数据管理'
    }
  },
  {
    path: '/doctor',
    name: 'Doctor',
    component: () => import('@/views/Doctor.vue'),
    meta: {
      title: '医师管理'
    }
  },
  {
    path: '/template',
    name: 'Template',
    component: () => import('@/views/Template.vue'),
    meta: {
      title: '模板管理'
    }
  },
  {
    path: '/medicine',
    name: 'Medicine',
    component: () => import('@/views/Medicine.vue'),
    meta: {
      title: '药品管理'
    }
  },
  {
    path: '/setting',
    name: 'Setting',
    component: () => import('@/views/Setting.vue'),
    meta: {
      title: '设置'
    }
  },
  {
    path: '/about',
    name: 'About',
    component: () => import('@/views/About.vue'),
    meta: {
      title: '关于'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - 医疗表单管理系统`
  }
  next()
})

export default router