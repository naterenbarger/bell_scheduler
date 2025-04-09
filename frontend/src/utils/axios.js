import axios from 'axios'
import store from '@/store'
import router from '@/router'

const instance = axios.create({
  baseURL: process.env.VUE_APP_API_URL || 'http://localhost:8080/api',
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
instance.interceptors.request.use(
  config => {
    const token = store.getters['auth/token']
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// Response interceptor
instance.interceptors.response.use(
  response => response,
  error => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          store.dispatch('auth/logout')
          // Only redirect to login if we're not already on the login page
          if (router.currentRoute.name !== 'Login') {
            router.push({ name: 'Login' })
          }
          break
        case 403:
          // Handle forbidden access
          break
        case 404:
          // Handle not found
          break
        case 500:
          // Handle server error
          break
        default:
          // Handle other errors
      }
    }
    return Promise.reject(error)
  }
)

export default instance 