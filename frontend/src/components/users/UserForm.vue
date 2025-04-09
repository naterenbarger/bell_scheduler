<template>
  <v-card>
    <v-card-title>
      {{ isEdit ? 'Edit User' : 'Create User' }}
    </v-card-title>

    <v-card-text>
      <v-form ref="form" v-model="valid">
        <v-text-field
          v-model="form.username"
          label="Username"
          :rules="[v => !!v || 'Username is required']"
          required
        ></v-text-field>

        <v-text-field
          v-model="form.email"
          label="Email"
          type="email"
          :rules="[
            v => !!v || 'Email is required',
            v => /.+@.+\..+/.test(v) || 'Email must be valid'
          ]"
          required
        ></v-text-field>

        <v-text-field
          v-model="form.password"
          label="Password"
          type="password"
          :rules="[
            v => isEdit || !!v || 'Password is required for new users',
            v => !v || v.length >= 8 || 'Password must be at least 8 characters'
          ]"
          :required="!isEdit"
        ></v-text-field>

        <v-select
          v-model="form.role"
          :items="roles"
          label="Role"
          :rules="[v => !!v || 'Role is required']"
          required
        ></v-select>

        <v-switch
          v-model="form.isActive"
          label="Active"
          color="primary"
        ></v-switch>
      </v-form>
    </v-card-text>

    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn text @click="$emit('cancel')">Cancel</v-btn>
      <v-btn
        color="primary"
        :loading="loading"
        :disabled="!valid"
        @click="handleSubmit"
      >
        {{ isEdit ? 'Update' : 'Create' }}
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
export default {
  name: 'UserForm',
  props: {
    user: {
      type: Object,
      default: null
    },
    isEdit: {
      type: Boolean,
      default: false
    },
    loading: {
      type: Boolean,
      default: false
    },
    error: {
      type: String,
      default: null
    }
  },
  data() {
    return {
      valid: false,
      form: {
        username: '',
        email: '',
        password: '',
        role: '',
        isActive: true
      },
      roles: ['admin', 'user']
    }
  },
  watch: {
    user: {
      handler(user) {
        if (user) {
          this.form = {
            username: user.username,
            email: user.email,
            password: '',
            role: user.role,
            isActive: user.isActive
          }
        } else {
          this.form = {
            username: '',
            email: '',
            password: '',
            role: '',
            isActive: true
          }
        }
      },
      immediate: true
    }
  },
  methods: {
    async handleSubmit() {
      if (!this.$refs.form.validate()) return

      try {
        console.log('Submitting form data:', this.form)
        if (this.isEdit) {
          await this.$store.dispatch('users/updateUser', {
            id: this.user.id,
            userData: this.form
          })
        } else {
          await this.$store.dispatch('users/createUser', this.form)
        }
        this.$emit('success')
      } catch (error) {
        console.error('Failed to save user:', error)
        console.error('Error response:', error.response?.data)
      }
    }
  }
}
</script> 