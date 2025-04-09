import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '../store'
import ChangePassword from '@/views/ChangePassword.vue'
import Logs from '@/views/Logs.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { guest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue'),
    meta: { guest: true }
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('../views/ForgotPassword.vue'),
    meta: { guest: true }
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: () => import('../views/ResetPassword.vue'),
    meta: { guest: true }
  },
  {
    path: '/force-password-change',
    name: 'ForcePasswordChange',
    component: () => import('../components/auth/ForcePasswordChange.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/schedules',
    name: 'Schedules',
    component: () => import('../views/Schedules.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/users',
    name: 'Users',
    component: () => import('../views/Users.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('../views/Settings.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/change-password',
    name: 'ChangePassword',
    component: ChangePassword,
    meta: { requiresAuth: true }
  },
  {
    path: '/logs',
    name: 'Logs',
    component: Logs,
    meta: { requiresAuth: true }
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

router.beforeEach((to, from, next) => {
  const isAuthenticated = store.getters['auth/isAuthenticated']
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const forcePasswordChange = store.getters['auth/forcePasswordChange']
  const isGuestRoute = to.matched.some(record => record.meta.guest)

  // If not authenticated and trying to access protected route
  if (requiresAuth && !isAuthenticated) {
    next('/login')
    return
  }

  // If authenticated and trying to access guest route
  if (isAuthenticated && isGuestRoute) {
    next('/')
    return
  }

  // If authenticated and force password change is required
  if (isAuthenticated && forcePasswordChange && to.name !== 'ForcePasswordChange') {
    next('/force-password-change')
    return
  }

  next()
})

export default router 