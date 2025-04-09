<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="12" sm="8" md="6">
        <v-card>
          <v-card-title class="headline">Register</v-card-title>
          <v-card-text>
            <ValidationObserver ref="observer" v-slot="{ handleSubmit }">
              <form @submit.prevent="handleSubmit(onSubmit)">
                <ValidationProvider v-slot="{ errors }" name="username" rules="required|min:3">
                  <v-text-field
                    v-model="form.username"
                    label="Username"
                    :error-messages="errors"
                    required
                  ></v-text-field>
                </ValidationProvider>

                <ValidationProvider v-slot="{ errors }" name="email" rules="required|email">
                  <v-text-field
                    v-model="form.email"
                    label="Email"
                    :error-messages="errors"
                    required
                  ></v-text-field>
                </ValidationProvider>

                <ValidationProvider v-slot="{ errors }" name="password" rules="required|min:6">
                  <v-text-field
                    v-model="form.password"
                    label="Password"
                    type="password"
                    :error-messages="errors"
                    required
                  ></v-text-field>
                </ValidationProvider>

                <v-btn type="submit" color="primary" block :loading="loading">
                  Register
                </v-btn>
              </form>
            </ValidationObserver>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn text to="/login">Already have an account? Login</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { ValidationObserver, ValidationProvider } from 'vee-validate'
import axios from '../plugins/axios'

export default {
  name: 'Register',
  components: {
    ValidationObserver,
    ValidationProvider
  },
  data() {
    return {
      form: {
        username: '',
        email: '',
        password: ''
      },
      loading: false
    }
  },
  methods: {
    async onSubmit() {
      this.loading = true
      try {
        await axios.post('/auth/register', this.form)
        this.$router.push('/login')
        this.$store.dispatch('snackbar/show', {
          message: 'Registration successful! Please login.',
          color: 'success'
        })
      } catch (error) {
        this.$store.dispatch('snackbar/show', {
          message: error.response?.data?.error || 'Registration failed',
          color: 'error'
        })
      } finally {
        this.loading = false
      }
    }
  }
}
</script> 