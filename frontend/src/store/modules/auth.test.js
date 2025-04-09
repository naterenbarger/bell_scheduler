import axios from 'axios'
import { auth } from './auth'

jest.mock('axios')

describe('auth store', () => {
  let state
  let commit
  let dispatch

  beforeEach(() => {
    state = {
      token: null,
      user: null,
      loading: false,
      error: null
    }
    commit = jest.fn()
    dispatch = jest.fn()
  })

  describe('mutations', () => {
    it('SET_TOKEN sets token', () => {
      const token = 'test_token'
      auth.mutations.SET_TOKEN(state, token)
      expect(state.token).toBe(token)
    })

    it('SET_USER sets user', () => {
      const user = { id: 1, username: 'testuser' }
      auth.mutations.SET_USER(state, user)
      expect(state.user).toBe(user)
    })

    it('SET_LOADING sets loading state', () => {
      auth.mutations.SET_LOADING(state, true)
      expect(state.loading).toBe(true)
    })

    it('SET_ERROR sets error message', () => {
      const error = 'Test error'
      auth.mutations.SET_ERROR(state, error)
      expect(state.error).toBe(error)
    })

    it('CLEAR_AUTH clears auth state', () => {
      state.token = 'test_token'
      state.user = { id: 1 }
      auth.mutations.CLEAR_AUTH(state)
      expect(state.token).toBeNull()
      expect(state.user).toBeNull()
    })
  })

  describe('actions', () => {
    it('login success', async () => {
      const credentials = {
        username: 'testuser',
        password: 'password123'
      }
      const response = {
        data: {
          token: 'test_token',
          user: { id: 1, username: 'testuser' }
        }
      }
      axios.post.mockResolvedValueOnce(response)

      await auth.actions.login({ commit, dispatch }, credentials)

      expect(commit).toHaveBeenCalledWith('SET_LOADING', true)
      expect(commit).toHaveBeenCalledWith('SET_TOKEN', 'test_token')
      expect(commit).toHaveBeenCalledWith('SET_USER', response.data.user)
      expect(commit).toHaveBeenCalledWith('SET_LOADING', false)
      expect(commit).toHaveBeenCalledWith('SET_ERROR', null)
    })

    it('login failure', async () => {
      const credentials = {
        username: 'testuser',
        password: 'wrong_password'
      }
      const error = new Error('Invalid credentials')
      axios.post.mockRejectedValueOnce(error)

      await auth.actions.login({ commit, dispatch }, credentials)

      expect(commit).toHaveBeenCalledWith('SET_LOADING', true)
      expect(commit).toHaveBeenCalledWith('SET_LOADING', false)
      expect(commit).toHaveBeenCalledWith('SET_ERROR', error.message)
    })

    it('logout', async () => {
      await auth.actions.logout({ commit, dispatch })

      expect(commit).toHaveBeenCalledWith('CLEAR_AUTH')
      expect(dispatch).toHaveBeenCalledWith('clearAuthData')
    })

    it('register success', async () => {
      const userData = {
        username: 'newuser',
        email: 'new@example.com',
        password: 'password123'
      }
      const response = {
        data: {
          token: 'test_token',
          user: { id: 1, username: 'newuser' }
        }
      }
      axios.post.mockResolvedValueOnce(response)

      await auth.actions.register({ commit, dispatch }, userData)

      expect(commit).toHaveBeenCalledWith('SET_LOADING', true)
      expect(commit).toHaveBeenCalledWith('SET_TOKEN', 'test_token')
      expect(commit).toHaveBeenCalledWith('SET_USER', response.data.user)
      expect(commit).toHaveBeenCalledWith('SET_LOADING', false)
      expect(commit).toHaveBeenCalledWith('SET_ERROR', null)
    })

    it('register failure', async () => {
      const userData = {
        username: 'existing',
        email: 'existing@example.com',
        password: 'password123'
      }
      const error = new Error('Username already exists')
      axios.post.mockRejectedValueOnce(error)

      await auth.actions.register({ commit, dispatch }, userData)

      expect(commit).toHaveBeenCalledWith('SET_LOADING', true)
      expect(commit).toHaveBeenCalledWith('SET_LOADING', false)
      expect(commit).toHaveBeenCalledWith('SET_ERROR', error.message)
    })
  })

  describe('getters', () => {
    it('isAuthenticated returns true when token exists', () => {
      state.token = 'test_token'
      expect(auth.getters.isAuthenticated(state)).toBe(true)
    })

    it('isAuthenticated returns false when token is null', () => {
      state.token = null
      expect(auth.getters.isAuthenticated(state)).toBe(false)
    })

    it('currentUser returns user object', () => {
      const user = { id: 1, username: 'testuser' }
      state.user = user
      expect(auth.getters.currentUser(state)).toBe(user)
    })

    it('isAdmin returns true for admin user', () => {
      state.user = { role: 'admin' }
      expect(auth.getters.isAdmin(state)).toBe(true)
    })

    it('isAdmin returns false for non-admin user', () => {
      state.user = { role: 'user' }
      expect(auth.getters.isAdmin(state)).toBe(false)
    })
  })
}) 