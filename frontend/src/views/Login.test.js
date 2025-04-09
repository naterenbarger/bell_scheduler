import { shallowMount } from '@vue/test-utils'
import { createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'
import Login from './Login.vue'

const localVue = createLocalVue()
localVue.use(Vuex)

describe('Login.vue', () => {
  let store
  let wrapper

  beforeEach(() => {
    store = new Vuex.Store({
      modules: {
        auth: {
          namespaced: true,
          state: {
            loading: false,
            error: null
          },
          actions: {
            login: jest.fn()
          }
        }
      }
    })

    wrapper = shallowMount(Login, {
      store,
      localVue,
      mocks: {
        $router: {
          push: jest.fn()
        }
      }
    })
  })

  it('renders login form', () => {
    expect(wrapper.find('form').exists()).toBe(true)
    expect(wrapper.find('input[name="username"]').exists()).toBe(true)
    expect(wrapper.find('input[name="password"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('validates required fields', async () => {
    const form = wrapper.find('form')
    await form.trigger('submit')
    
    expect(wrapper.find('.error-message').text()).toBe('Username is required')
    expect(store.dispatch).not.toHaveBeenCalled()
  })

  it('submits form with valid data', async () => {
    const username = 'testuser'
    const password = 'password123'

    await wrapper.setData({
      username,
      password
    })

    const form = wrapper.find('form')
    await form.trigger('submit')

    expect(store.dispatch).toHaveBeenCalledWith('auth/login', {
      username,
      password
    })
  })

  it('displays error message from store', async () => {
    const error = 'Invalid credentials'
    await store.commit('auth/SET_ERROR', error)

    expect(wrapper.find('.error-message').text()).toBe(error)
  })

  it('shows loading state', async () => {
    await store.commit('auth/SET_LOADING', true)

    expect(wrapper.find('button[type="submit"]').attributes('disabled')).toBeTruthy()
  })

  it('navigates to dashboard on successful login', async () => {
    const username = 'testuser'
    const password = 'password123'

    await wrapper.setData({
      username,
      password
    })

    store.dispatch.mockResolvedValueOnce()
    const form = wrapper.find('form')
    await form.trigger('submit')

    expect(wrapper.vm.$router.push).toHaveBeenCalledWith('/dashboard')
  })
}) 