import axios from '@/utils/axios'

const state = {
  users: [],
  loading: false,
  error: null,
  pagination: {
    page: 1,
    itemsPerPage: 10,
    totalItems: 0
  },
  sort: {
    by: ['username'],
    desc: [false]
  },
  filter: {
    search: ''
  }
}

const mutations = {
  SET_USERS(state, { users, total }) {
    state.users = users
    state.pagination.totalItems = total
  },
  SET_LOADING(state, loading) {
    state.loading = loading
  },
  SET_ERROR(state, error) {
    state.error = error
  },
  SET_PAGINATION(state, { page, itemsPerPage }) {
    state.pagination.page = page
    state.pagination.itemsPerPage = itemsPerPage
  },
  SET_SORT(state, { by, desc }) {
    state.sort.by = by
    state.sort.desc = desc
  },
  SET_FILTER(state, { search }) {
    state.filter.search = search
  }
}

const actions = {
  async fetchUsers({ commit, state }) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      const { page, itemsPerPage } = state.pagination
      const { by, desc } = state.sort
      const { search } = state.filter

      const params = {
        page,
        limit: itemsPerPage,
        sort_by: by.join(','),
        sort_desc: desc.join(','),
        search
      }

      const response = await axios.get('/users', { params })
      commit('SET_USERS', {
        users: response.data.users || [],
        total: response.data.total || 0
      })
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || error.response?.data?.message || 'Failed to fetch users')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async createUser({ commit, dispatch }, userData) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      await axios.post('/users', userData)
      await dispatch('fetchUsers')
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || error.response?.data?.message || 'Failed to create user')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async updateUser({ commit, dispatch }, { id, userData }) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      await axios.put(`/users/${id}`, userData)
      await dispatch('fetchUsers')
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || error.response?.data?.message || 'Failed to update user')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async deleteUser({ commit, dispatch }, id) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    try {
      await axios.delete(`/users/${id}`)
      await dispatch('fetchUsers')
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || error.response?.data?.message || 'Failed to delete user')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  updatePagination({ commit }, { page, itemsPerPage }) {
    commit('SET_PAGINATION', { page, itemsPerPage })
    // Removed automatic fetchUsers call to prevent duplicate API calls
  },

  updateSort({ commit }, { by, desc }) {
    commit('SET_SORT', { by, desc })
    // Removed automatic fetchUsers call to prevent duplicate API calls
  },

  updateFilter({ commit }, { search }) {
    commit('SET_FILTER', { search })
    // Removed automatic fetchUsers call to prevent duplicate API calls
  }
}

const getters = {
  users: state => state.users,
  loading: state => state.loading,
  error: state => state.error,
  pagination: state => state.pagination,
  sort: state => state.sort,
  filter: state => state.filter
}

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
}