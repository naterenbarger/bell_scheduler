<template>
  <v-container fluid fill-height>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="elevation-12">
          <v-toolbar color="primary" dark flat>
            <v-toolbar-title>Login</v-toolbar-title>
          </v-toolbar>
          <v-card-text>
            <ValidationObserver ref="observer" v-slot="{ handleSubmit }">
              <form @submit.prevent="handleSubmit(onSubmit)">
                <ValidationProvider v-slot="{ errors }" name="username" rules="required">
                  <v-text-field
                    v-model="form.username"
                    :error-messages="errors"
                    label="Username"
                    prepend-icon="mdi-account"
                    type="text"
                    required
                  ></v-text-field>
                </ValidationProvider>

                <ValidationProvider v-slot="{ errors }" name="password" rules="required">
                  <v-text-field
                    v-model="form.password"
                    :error-messages="errors"
                    label="Password"
                    prepend-icon="mdi-lock"
                    :type="showPassword ? 'text' : 'password'"
                    :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                    @click:append="showPassword = !showPassword"
                    required
                  ></v-text-field>
                </ValidationProvider>

                <v-alert
                  v-if="error"
                  type="error"
                  text
                  dense
                >
                  {{ error }}
                </v-alert>

                <v-btn
                  color="primary"
                  block
                  type="submit"
                  :loading="loading"
                >
                  Login
                </v-btn>

                <!-- Commented out for now, but keeping for future use
                <v-btn
                  text
                  block
                  class="mt-2"
                  @click="$router.push('/forgot-password')"
                >
                  Forgot Password?
                </v-btn>

                <v-btn
                  text
                  block
                  class="mt-2"
                  @click="$router.push('/register')"
                >
                  Don't have an account? Register
                </v-btn>
                -->
              </form>
            </ValidationObserver>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { ValidationObserver, ValidationProvider, extend } from 'vee-validate'
import { required } from 'vee-validate/dist/rules'

// Register validation rules
extend('required', {
  ...required,
  message: '{_field_} is required'
})

export default {
  name: 'LoginView',
  components: {
    ValidationObserver,
    ValidationProvider
  },
  data() {
    return {
      loading: false,
      showPassword: false,
      error: null,
      form: {
        username: '',
        password: ''
      }
    }
  },
  methods: {
    async onSubmit() {
      this.loading = true
      this.error = null

      try {
        await this.$store.dispatch('auth/login', this.form)
        
        // Check if force password change is required
        const forcePasswordChange = this.$store.getters['auth/forcePasswordChange']
        if (forcePasswordChange) {
          this.$router.push('/force-password-change')
        } else {
          this.$router.push('/')
        }
      } catch (error) {
        this.error = error.message || 'Login failed'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>