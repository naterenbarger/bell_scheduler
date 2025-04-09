import Vue from 'vue'
import axios from 'axios'
import store from '../store'

// Create axios instance with base URL
const instance = axios.create({
  baseURL: process.env.VUE_APP_API_URL || 'http://localhost:8080/api',
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Add request interceptor
instance.interceptors.request.use(
  config => {
    const token = store.state.auth.token
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// Add response interceptor
instance.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      store.dispatch('auth/logout')
    }
    return Promise.reject(error)
  }
)

// Register axios instance as a Vue prototype property
Vue.prototype.$axios = instance

export default instance