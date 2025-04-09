import axios from 'axios'
import { schedules } from './schedules'

jest.mock('axios')

describe('schedules store', () => {
  let state
  let commit
  let dispatch

  beforeEach(() => {
    state = {
      schedules: [],
      loading: false,
      error: null
    }
    commit = jest.fn()
    dispatch = jest.fn()
  })

  describe('mutations', () => {
    it('SET_SCHEDULES sets schedules', () => {
      const schedules = [
        { id: 1, name: 'Schedule 1' },
        { id: 2, name: 'Schedule 2' }
      ]
      schedules.mutations.SET_SCHEDULES(state, schedules)
      expect(state.schedules).toEqual(schedules)
    })

    it('SET_LOADING sets loading state', () => {
      schedules.mutations.SET_LOADING(state, true)
      expect(state.loading).toBe(true)
    })

    it('SET_ERROR sets error message', () => {
      const error = 'Test error'
      schedules.mutations.SET_ERROR(state, error)
      expect(state.error).toBe(error)
    })

    it('ADD_SCHEDULE adds a schedule', () => {
      const schedule = { id: 1, name: 'New Schedule' }
      schedules.mutations.ADD_SCHEDULE(state, schedule)
      expect(state.schedules).toContainEqual(schedule)
    })

    it('UPDATE_SCHEDULE updates a schedule', () => {
      const oldSchedule = { id: 1, name: 'Old Name' }
      const newSchedule = { id: 1, name: 'New Name' }
      state.schedules = [oldSchedule]
      schedules.mutations.UPDATE_SCHEDULE(state, newSchedule)
      expect(state.schedules[0]).toEqual(newSchedule)
    })

    it('DELETE_SCHEDULE removes a schedule', () => {
      const schedule = { id: 1, name: 'Schedule 1' }
      state.schedules = [schedule]
      schedules.mutations.DELETE_SCHEDULE(state, 1)
      expect(state.schedules).not.toContainEqual(schedule)
    })
  })

  describe('actions', () => {
    it('fetchSchedules success', async () => {
      const response = {
        data: [
          { id: 1, name: 'Schedule 1' },
          { id: 2, name: 'Schedule 2' }
        ]
      }
      axios.get.mockResolvedValueOnce(response)

      await schedules.actions.fetchSchedules({ commit, dispatch })

      expect(commit).toHaveBeenCalledWith('SET_LOADING', true)
      expect(commit).toHaveBeenCalledWith('SET_SCHEDULES', response.data)
      expect(commit).toHaveBeenCalledWith('SET_LOADING', false)
      expect(commit).toHaveBeenCalledWith('SET_ERROR', null)
    })

    it('fetchSchedules failure', async () => {
      const error = new Error('Failed to fetch schedules')
      axios.get.mockRejectedValueOnce(error)

      await schedules.actions.fetchSchedules({ commit, dispatch })

      expect(commit).toHaveBeenCalledWith('SET_LOADING', true)
      expect(commit).toHaveBeenCalledWith('SET_LOADING', false)
      expect(commit).toHaveBeenCalledWith('SET_ERROR', error.message)
    })

    it('createSchedule success', async () => {
      const schedule = { name: 'New Schedule', description: 'Test' }
      const response = { data: { id: 1, ...schedule } }
      axios.post.mockResolvedValueOnce(response)

      await schedules.actions.createSchedule({ commit, dispatch }, schedule)

      expect(commit).toHaveBeenCalledWith('SET_LOADING', true)
      expect(commit).toHaveBeenCalledWith('ADD_SCHEDULE', response.data)
      expect(commit).toHaveBeenCalledWith('SET_LOADING', false)
      expect(commit).toHaveBeenCalledWith('SET_ERROR', null)
    })

    it('updateSchedule success', async () => {
      const schedule = { id: 1, name: 'Updated Schedule' }
      const response = { data: schedule }
      axios.put.mockResolvedValueOnce(response)

      await schedules.actions.updateSchedule({ commit, dispatch }, schedule)

      expect(commit).toHaveBeenCalledWith('SET_LOADING', true)
      expect(commit).toHaveBeenCalledWith('UPDATE_SCHEDULE', schedule)
      expect(commit).toHaveBeenCalledWith('SET_LOADING', false)
      expect(commit).toHaveBeenCalledWith('SET_ERROR', null)
    })

    it('deleteSchedule success', async () => {
      const scheduleId = 1
      axios.delete.mockResolvedValueOnce({})

      await schedules.actions.deleteSchedule({ commit, dispatch }, scheduleId)

      expect(commit).toHaveBeenCalledWith('SET_LOADING', true)
      expect(commit).toHaveBeenCalledWith('DELETE_SCHEDULE', scheduleId)
      expect(commit).toHaveBeenCalledWith('SET_LOADING', false)
      expect(commit).toHaveBeenCalledWith('SET_ERROR', null)
    })

    it('setDefaultSchedule success', async () => {
      const scheduleId = 1
      axios.post.mockResolvedValueOnce({})

      await schedules.actions.setDefaultSchedule({ commit, dispatch }, scheduleId)

      expect(commit).toHaveBeenCalledWith('SET_LOADING', true)
      expect(commit).toHaveBeenCalledWith('SET_LOADING', false)
      expect(commit).toHaveBeenCalledWith('SET_ERROR', null)
    })
  })

  describe('getters', () => {
    it('allSchedules returns all schedules', () => {
      const schedules = [
        { id: 1, name: 'Schedule 1' },
        { id: 2, name: 'Schedule 2' }
      ]
      state.schedules = schedules
      expect(schedules.getters.allSchedules(state)).toEqual(schedules)
    })

    it('defaultSchedule returns the default schedule', () => {
      const schedules = [
        { id: 1, name: 'Schedule 1', is_default: false },
        { id: 2, name: 'Schedule 2', is_default: true }
      ]
      state.schedules = schedules
      expect(schedules.getters.defaultSchedule(state)).toEqual(schedules[1])
    })

    it('defaultSchedule returns null when no default schedule', () => {
      const schedules = [
        { id: 1, name: 'Schedule 1', is_default: false },
        { id: 2, name: 'Schedule 2', is_default: false }
      ]
      state.schedules = schedules
      expect(schedules.getters.defaultSchedule(state)).toBeNull()
    })
  })
}) 