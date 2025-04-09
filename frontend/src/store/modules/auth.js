import axios from '@/utils/axios'

const state = {
  token: localStorage.getItem('token') || null,
  user: JSON.parse(localStorage.getItem('user')) || null,
  loading: false,
  error: null
}

const mutations = {
  SET_TOKEN(state, token) {
    state.token = token
    if (token) {
      localStorage.setItem('token', token)
    } else {
      localStorage.removeItem('token')
    }
  },
  SET_USER(state, user) {
    state.user = user
    if (user) {
      localStorage.setItem('user', JSON.stringify(user))
    } else {
      localStorage.removeItem('user')
    }
  },
  SET_LOADING(state, loading) {
    state.loading = loading
  },
  SET_ERROR(state, error) {
    state.error = error
  }
}

const actions = {
  async login({ commit }, credentials) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      const response = await axios.post('/auth/login', credentials)
      commit('SET_TOKEN', response.data.token)
      commit('SET_USER', response.data.user)
      return response.data
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || 'Login failed')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async register({ commit }, userData) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      const response = await axios.post('/auth/register', userData)
      return response.data
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || 'Registration failed')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async forgotPassword({ commit }, { email }) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      const response = await axios.post('/auth/forgot-password', { email })
      return response.data
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || 'Failed to send reset email')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async resetPassword({ commit }, { token, password }) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      const response = await axios.post('/auth/reset-password', { token, password })
      return response.data
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || 'Failed to reset password')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  logout({ commit }) {
    commit('SET_TOKEN', null)
    commit('SET_USER', null)
  },

  async checkAuth({ commit, state }) {
    if (!state.token) return false

    try {
      const response = await axios.get('/auth/me')
      commit('SET_USER', response.data)
      return true
    } catch (error) {
      commit('SET_TOKEN', null)
      commit('SET_USER', null)
      return false
    }
  },

  async changePassword({ commit, state }, { currentPassword, newPassword }) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      const response = await axios.post('/auth/change-password', {
        currentPassword,
        newPassword
      })
      commit('SET_USER', response.data.user)
      return response.data
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || 'Failed to change password')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  }
}

const getters = {
  isAuthenticated: state => !!state.token,
  user: state => state.user,
  loading: state => state.loading,
  error: state => state.error,
  isAdmin: state => state.user?.role === 'admin',
  token: state => state.token,
  forcePasswordChange: state => state.user?.ForcePasswordChange || false
}

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
} 