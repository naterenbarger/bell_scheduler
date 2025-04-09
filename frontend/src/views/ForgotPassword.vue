<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="12" sm="8" md="6">
        <v-card>
          <v-card-title class="headline">Forgot Password</v-card-title>
          <v-card-text>
            <p class="mb-4">Enter your email address and we'll send you a link to reset your password.</p>
            <ValidationObserver ref="observer" v-slot="{ handleSubmit }">
              <form @submit.prevent="handleSubmit(onSubmit)">
                <ValidationProvider v-slot="{ errors }" name="email" rules="required|email">
                  <v-text-field
                    v-model="email"
                    label="Email"
                    type="email"
                    :error-messages="errors"
                    required
                  ></v-text-field>
                </ValidationProvider>

                <v-btn type="submit" color="primary" block :loading="loading">
                  Send Reset Link
                </v-btn>
              </form>
            </ValidationObserver>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn text to="/login">Back to Login</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { ValidationObserver, ValidationProvider, extend } from 'vee-validate'
import { email, required } from 'vee-validate/dist/rules'
import axios from '../plugins/axios'

// Register validation rules
extend('required', {
  ...required,
  message: '{_field_} is required'
})

extend('email', {
  ...email,
  message: 'Please enter a valid email address'
})

export default {
  name: 'ForgotPasswordView',
  components: {
    ValidationObserver,
    ValidationProvider
  },
  data() {
    return {
      email: '',
      loading: false
    }
  },
  methods: {
    async onSubmit() {
      this.loading = true
      try {
        await axios.post('/auth/forgot-password', { email: this.email })
        this.$store.dispatch('snackbar/show', {
          message: 'If the email exists, a password reset link has been sent',
          color: 'success'
        })
        this.$router.push('/login')
      } catch (error) {
        this.$store.dispatch('snackbar/show', {
          message: error.response?.data?.error || 'Failed to send reset link',
          color: 'error'
        })
      } finally {
        this.loading = false
      }
    }
  }
}
</script> 