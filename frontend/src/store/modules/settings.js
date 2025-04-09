import axios from '@/utils/axios'

const state = {
  settings: {
    ringDuration: 30, // Default duration in seconds
    timezone: 'UTC'
  },
  loading: false,
  error: null
}

const mutations = {
  SET_SETTINGS(state, settings) {
    state.settings = settings
  },
  SET_LOADING(state, loading) {
    state.loading = loading
  },
  SET_ERROR(state, error) {
    state.error = error
  },
  UPDATE_SETTING(state, { key, value }) {
    state.settings[key] = value
  }
}

const actions = {
  async fetchSettings({ commit }) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      const response = await axios.get('/settings')
      commit('SET_SETTINGS', response.data)
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || 'Failed to fetch settings')
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async updateSettings({ commit }, settings) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      const response = await axios.put('/settings', {
        ringDuration: settings.ringDuration,
        timezone: settings.timezone,
        gpioPin: settings.gpioPin
      })
      commit('SET_SETTINGS', response.data)
      return response.data
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || 'Failed to update settings')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  }
}

const getters = {
  settings: state => state.settings,
  ringDuration: state => state.settings.ringDuration,
  timezone: state => state.settings.timezone,
  loading: state => state.loading,
  error: state => state.error
}

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
} 