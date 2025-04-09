import { shallowMount } from '@vue/test-utils'
import { createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'
import ScheduleList from './ScheduleList.vue'

const localVue = createLocalVue()
localVue.use(Vuex)

describe('ScheduleList.vue', () => {
  let store
  let wrapper

  beforeEach(() => {
    store = new Vuex.Store({
      modules: {
        schedules: {
          namespaced: true,
          state: {
            schedules: [],
            loading: false,
            error: null
          },
          actions: {
            fetchSchedules: jest.fn(),
            deleteSchedule: jest.fn()
          }
        }
      }
    })

    wrapper = shallowMount(ScheduleList, {
      store,
      localVue,
      mocks: {
        $router: {
          push: jest.fn()
        }
      }
    })
  })

  it('renders schedule list', () => {
    expect(wrapper.find('.schedule-list').exists()).toBe(true)
  })

  it('displays loading state', async () => {
    await store.commit('schedules/SET_LOADING', true)
    expect(wrapper.find('.loading').exists()).toBe(true)
  })

  it('displays error message', async () => {
    const error = 'Failed to load schedules'
    await store.commit('schedules/SET_ERROR', error)
    expect(wrapper.find('.error-message').text()).toBe(error)
  })

  it('renders schedule items', async () => {
    const schedules = [
      { id: 1, name: 'Schedule 1', description: 'Test schedule 1' },
      { id: 2, name: 'Schedule 2', description: 'Test schedule 2' }
    ]
    await store.commit('schedules/SET_SCHEDULES', schedules)

    const items = wrapper.findAll('.schedule-item')
    expect(items).toHaveLength(2)
    expect(items.at(0).text()).toContain('Schedule 1')
    expect(items.at(1).text()).toContain('Schedule 2')
  })

  it('navigates to schedule detail on click', async () => {
    const schedule = { id: 1, name: 'Schedule 1' }
    await store.commit('schedules/SET_SCHEDULES', [schedule])

    const item = wrapper.find('.schedule-item')
    await item.trigger('click')

    expect(wrapper.vm.$router.push).toHaveBeenCalledWith(`/schedules/${schedule.id}`)
  })

  it('deletes schedule on delete button click', async () => {
    const schedule = { id: 1, name: 'Schedule 1' }
    await store.commit('schedules/SET_SCHEDULES', [schedule])

    const deleteButton = wrapper.find('.delete-button')
    await deleteButton.trigger('click')

    expect(store.dispatch).toHaveBeenCalledWith('schedules/deleteSchedule', schedule.id)
  })

  it('shows confirmation dialog before deleting', async () => {
    const schedule = { id: 1, name: 'Schedule 1' }
    await store.commit('schedules/SET_SCHEDULES', [schedule])

    const deleteButton = wrapper.find('.delete-button')
    await deleteButton.trigger('click')

    expect(wrapper.find('.confirm-dialog').exists()).toBe(true)
  })

  it('cancels deletion when confirmed', async () => {
    const schedule = { id: 1, name: 'Schedule 1' }
    await store.commit('schedules/SET_SCHEDULES', [schedule])

    const deleteButton = wrapper.find('.delete-button')
    await deleteButton.trigger('click')

    const confirmButton = wrapper.find('.confirm-button')
    await confirmButton.trigger('click')

    expect(store.dispatch).toHaveBeenCalledWith('schedules/deleteSchedule', schedule.id)
  })

  it('cancels deletion when cancelled', async () => {
    const schedule = { id: 1, name: 'Schedule 1' }
    await store.commit('schedules/SET_SCHEDULES', [schedule])

    const deleteButton = wrapper.find('.delete-button')
    await deleteButton.trigger('click')

    const cancelButton = wrapper.find('.cancel-button')
    await cancelButton.trigger('click')

    expect(store.dispatch).not.toHaveBeenCalled()
  })

  it('fetches schedules on mount', () => {
    expect(store.dispatch).toHaveBeenCalledWith('schedules/fetchSchedules')
  })
}) 