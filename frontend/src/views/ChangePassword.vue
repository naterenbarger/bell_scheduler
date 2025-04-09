<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="12" sm="8" md="6">
        <v-card>
          <v-card-title class="headline">Change Password</v-card-title>
          <v-card-text>
            <ValidationObserver ref="observer" v-slot="{ handleSubmit }">
              <form @submit.prevent="handleSubmit(onSubmit)">
                <ValidationProvider v-slot="{ errors }" name="current password" rules="required">
                  <v-text-field
                    v-model="form.currentPassword"
                    label="Current Password"
                    type="password"
                    :error-messages="errors"
                    required
                  ></v-text-field>
                </ValidationProvider>

                <ValidationProvider v-slot="{ errors }" name="new password" rules="required|min:6">
                  <v-text-field
                    v-model="form.newPassword"
                    label="New Password"
                    type="password"
                    :error-messages="errors"
                    required
                  ></v-text-field>
                </ValidationProvider>

                <ValidationProvider v-slot="{ errors }" name="confirm password" rules="required|confirmed:new password">
                  <v-text-field
                    v-model="form.confirmPassword"
                    label="Confirm New Password"
                    type="password"
                    :error-messages="errors"
                    required
                  ></v-text-field>
                </ValidationProvider>

                <v-btn type="submit" color="primary" block :loading="loading">
                  Change Password
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
import { confirmed, min, required } from 'vee-validate/dist/rules'
import axios from '../plugins/axios'

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

export default {
  name: 'ChangePasswordView',
  components: {
    ValidationObserver,
    ValidationProvider
  },
  data() {
    return {
      form: {
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      },
      loading: false
    }
  },
  methods: {
    async onSubmit() {
      this.loading = true
      try {
        await axios.post('/auth/change-password', {
          currentPassword: this.form.currentPassword,
          newPassword: this.form.newPassword
        })
        this.$store.dispatch('snackbar/show', {
          message: 'Password changed successfully',
          color: 'success'
        })
        this.$router.push('/dashboard')
      } catch (error) {
        this.$store.dispatch('snackbar/show', {
          message: error.response?.data?.error || 'Failed to change password',
          color: 'error'
        })
      } finally {
        this.loading = false
      }
    }
  }
}
</script> 