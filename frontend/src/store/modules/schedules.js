import axios from '@/utils/axios'

const state = {
  schedules: [],
  currentSchedule: null,
  loading: false,
  error: null
}

const mutations = {
  SET_SCHEDULES(state, schedules) {
    state.schedules = schedules
  },
  SET_CURRENT_SCHEDULE(state, schedule) {
    state.currentSchedule = schedule
  },
  SET_LOADING(state, loading) {
    state.loading = loading
  },
  SET_ERROR(state, error) {
    state.error = error
  },
  ADD_SCHEDULE(state, schedule) {
    state.schedules.push(schedule)
  },
  UPDATE_SCHEDULE(state, updatedSchedule) {
    const index = state.schedules.findIndex(s => s.id === updatedSchedule.id)
    if (index !== -1) {
      state.schedules.splice(index, 1, updatedSchedule)
    }
    if (state.currentSchedule?.id === updatedSchedule.id) {
      state.currentSchedule = updatedSchedule
    }
  },
  DELETE_SCHEDULE(state, scheduleId) {
    state.schedules = state.schedules.filter(s => s.id !== scheduleId)
    if (state.currentSchedule?.id === scheduleId) {
      state.currentSchedule = null
    }
  }
}

const actions = {
  async fetchSchedules({ commit }) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      const response = await axios.get('/schedules')
      commit('SET_SCHEDULES', response.data)
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || 'Failed to fetch schedules')
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async fetchSchedule({ commit }, id) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      const response = await axios.get(`/schedules/${id}`)
      commit('SET_CURRENT_SCHEDULE', response.data)
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || 'Failed to fetch schedule')
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async createSchedule({ commit }, scheduleData) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      const response = await axios.post('/schedules', scheduleData)
      commit('ADD_SCHEDULE', response.data)
      return response.data
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || 'Failed to create schedule')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async updateSchedule({ commit }, { id, data }) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      const response = await axios.put(`/schedules/${id}`, data)
      commit('UPDATE_SCHEDULE', response.data)
      return response.data
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || 'Failed to update schedule')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async deleteSchedule({ commit }, id) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      await axios.delete(`/schedules/${id}`)
      commit('DELETE_SCHEDULE', id)
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || 'Failed to delete schedule')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async setDefaultSchedule({ commit }, id) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      await axios.put(`/schedules/${id}/default`)
      await this.dispatch('fetchSchedules')
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || 'Failed to set default schedule')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async setTemporarySchedule({ commit }, { id, isTemporary }) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      await axios.put(`/schedules/${id}/temporary`, { isTemporary })
      await this.dispatch('fetchSchedules')
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || 'Failed to set temporary schedule')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async setActiveSchedule({ commit }, id) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      await axios.put(`/schedules/${id}/active`)
      await this.dispatch('fetchSchedules')
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || 'Failed to set active schedule')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async triggerBell({ commit }) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      // Since this is a manual trigger without a specific schedule ID,
      // we'll use the default schedule if available
      const defaultSchedule = this.getters['schedules/defaultSchedule']
      if (defaultSchedule) {
        await axios.post(`/schedules/${defaultSchedule.id}/trigger`)
        return true
      } else if (state.schedules.length > 0) {
        // If no default, use the first available schedule
        await axios.post(`/schedules/${state.schedules[0].id}/trigger`)
        return true
      } else {
        throw new Error('No schedules available to trigger')
      }
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || error.message || 'Failed to trigger bell')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  }
}

const getters = {
  schedules: state => state.schedules,
  currentSchedule: state => state.currentSchedule,
  defaultSchedule: state => state.schedules.find(s => s.isDefault),
  activeSchedule: state => state.schedules.find(s => s.isActive),
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