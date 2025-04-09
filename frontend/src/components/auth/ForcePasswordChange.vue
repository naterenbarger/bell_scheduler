<template>
  <v-container fluid fill-height>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="elevation-12">
          <v-toolbar color="primary" dark flat>
            <v-toolbar-title>Change Password Required</v-toolbar-title>
          </v-toolbar>
          <v-card-text>
            <p class="subtitle-1 mb-4">
              You are required to change your password before continuing.
            </p>
            <ValidationObserver ref="observer" v-slot="{ handleSubmit }">
              <form @submit.prevent="handleSubmit(onSubmit)">
                <ValidationProvider v-slot="{ errors }" name="current password" rules="required">
                  <v-text-field
                    v-model="form.currentPassword"
                    :error-messages="errors"
                    label="Current Password"
                    prepend-icon="mdi-lock"
                    :type="showCurrentPassword ? 'text' : 'password'"
                    :append-icon="showCurrentPassword ? 'mdi-eye' : 'mdi-eye-off'"
                    @click:append="showCurrentPassword = !showCurrentPassword"
                    required
                  ></v-text-field>
                </ValidationProvider>

                <ValidationProvider v-slot="{ errors }" name="new password" rules="required|min:8">
                  <v-text-field
                    v-model="form.newPassword"
                    :error-messages="errors"
                    label="New Password"
                    prepend-icon="mdi-lock"
                    :type="showNewPassword ? 'text' : 'password'"
                    :append-icon="showNewPassword ? 'mdi-eye' : 'mdi-eye-off'"
                    @click:append="showNewPassword = !showNewPassword"
                    required
                  ></v-text-field>
                </ValidationProvider>

                <ValidationProvider v-slot="{ errors }" name="confirm password" rules="required|confirmed:new password">
                  <v-text-field
                    v-model="form.confirmPassword"
                    :error-messages="errors"
                    label="Confirm New Password"
                    prepend-icon="mdi-lock"
                    :type="showConfirmPassword ? 'text' : 'password'"
                    :append-icon="showConfirmPassword ? 'mdi-eye' : 'mdi-eye-off'"
                    @click:append="showConfirmPassword = !showConfirmPassword"
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
  message: '{_field_} must match {target}'
})

export default {
  name: 'ForcePasswordChange',
  components: {
    ValidationObserver,
    ValidationProvider
  },
  data() {
    return {
      loading: false,
      error: null,
      showCurrentPassword: false,
      showNewPassword: false,
      showConfirmPassword: false,
      form: {
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      }
    }
  },
  methods: {
    async onSubmit() {
      this.loading = true
      this.error = null

      try {
        await this.$store.dispatch('auth/changePassword', {
          currentPassword: this.form.currentPassword,
          newPassword: this.form.newPassword
        })
        
        // After successful password change, redirect to dashboard
        this.$router.push('/')
      } catch (error) {
        this.error = error.message || 'Failed to change password'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>