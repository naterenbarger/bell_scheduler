<template>
  <v-container fluid fill-height>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="elevation-12">
          <v-toolbar color="primary" dark flat>
            <v-toolbar-title>Reset Password</v-toolbar-title>
          </v-toolbar>
          <v-card-text>
            <ValidationObserver ref="observer" v-slot="{ handleSubmit }">
              <form @submit.prevent="handleSubmit(onSubmit)">
                <ValidationProvider v-slot="{ errors }" name="password" rules="required|min:8|password">
                  <v-text-field
                    v-model="form.password"
                    :error-messages="errors"
                    label="New Password"
                    prepend-icon="mdi-lock"
                    :type="showPassword ? 'text' : 'password'"
                    :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                    @click:append="showPassword = !showPassword"
                    required
                  ></v-text-field>
                </ValidationProvider>

                <ValidationProvider v-slot="{ errors }" name="confirm password" rules="required|confirmed:password">
                  <v-text-field
                    v-model="form.confirmPassword"
                    :error-messages="errors"
                    label="Confirm Password"
                    prepend-icon="mdi-lock-check"
                    :type="showPassword ? 'text' : 'password'"
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
                  Reset Password
                </v-btn>

                <v-btn
                  text
                  block
                  class="mt-2"
                  @click="$router.push('/login')"
                >
                  Back to Login
                </v-btn>
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
import { required, min, confirmed } from 'vee-validate/dist/rules'

// Register validation rules
extend('required', {
  ...required,
  message: '{_field_} is required'
})

extend('min', {
  ...min,
  message: '{_field_} must be at least {length} characters'
})

extend('confirmed', {
  ...confirmed,
  message: 'Passwords do not match'
})

// Custom password validation rule - simplified to match backend requirements
extend('password', {
  validate: value => {
    return value.length >= 8
  },
  message: 'Password must be at least 8 characters long'
})

export default {
  name: 'ResetPasswordView',
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
        password: '',
        confirmPassword: ''
      }
    }
  },
  created() {
    // Get token from URL query parameters
    const token = this.$route.query.token
    if (!token) {
      this.$router.push('/login')
    }
  },
  methods: {
    async onSubmit() {
      this.loading = true
      this.error = null

      try {
        await this.$store.dispatch('auth/resetPassword', {
          token: this.$route.query.token,
          password: this.form.password
        })
        this.$router.push('/login')
      } catch (error) {
        this.error = error.response?.data?.error || 'An error occurred'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>