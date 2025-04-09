<template>
  <v-form ref="form" v-model="valid">
    <v-card>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
              v-model.number="formData.ringDuration"
              type="number"
              label="Bell Ring Duration (seconds)"
              :rules="[
                v => !!v || 'Duration is required',
                v => v > 0 || 'Duration must be greater than 0',
                v => v <= 300 || 'Duration cannot exceed 5 minutes'
              ]"
              required
            >
              <template v-slot:append>
                <div class="text-caption">seconds</div>
              </template>
            </v-text-field>
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field
              v-model.number="formData.gpioPin"
              type="number"
              label="GPIO Pin"
              :rules="[
                v => !!v || 'GPIO Pin is required',
                v => v > 0 || 'GPIO Pin must be greater than 0'
              ]"
              required
            />
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" md="6">
            <v-select
              v-model="formData.timezone"
              :items="timezones"
              label="Default Timezone"
              :rules="[v => !!v || 'Timezone is required']"
              required
            />
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12">
            <v-alert
              type="info"
              dense
              text
            >
              The bell ring duration setting will be applied to all schedules unless overridden in individual schedules.
            </v-alert>
          </v-col>
        </v-row>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn
          color="primary"
          :loading="loading"
          :disabled="!valid"
          @click="handleSubmit"
        >
          Save Settings
        </v-btn>
      </v-card-actions>
    </v-card>

    <!-- Error Snackbar -->
    <v-snackbar
      v-model="showError"
      color="error"
      timeout="5000"
    >
      {{ error }}
    </v-snackbar>

    <!-- Success Snackbar -->
    <v-snackbar
      v-model="showSuccess"
      color="success"
      timeout="3000"
    >
      Settings saved successfully
    </v-snackbar>
  </v-form>
</template>

<script>
import { timezones } from '@/utils/timezones'

export default {
  name: 'SettingsForm',
  props: {
    settings: {
      type: Object,
      required: true
    },
    loading: {
      type: Boolean,
      default: false
    },
    error: {
      type: String,
      default: ''
    }
  },
  data: () => ({
    valid: false,
    formData: {
      ringDuration: 30,
      timezone: 'UTC',
      gpioPin: 17
    },
    timezones,
    showError: false,
    showSuccess: false
  }),
  watch: {
    settings: {
      handler(newSettings) {
        if (newSettings) {
          this.formData = {
            ringDuration: newSettings.ringDuration || 30,
            timezone: newSettings.timezone || 'UTC',
            gpioPin: newSettings.gpioPin || 17
          }
        }
      },
      immediate: true,
      deep: true
    },
    error(newError) {
      this.showError = !!newError
    }
  },
  methods: {
    async handleSubmit() {
      if (!this.$refs.form.validate()) return

      try {
        await this.$store.dispatch('settings/updateSettings', this.formData)
        this.showSuccess = true
      } catch (error) {
        console.error('Failed to save settings:', error)
      }
    }
  }
}
</script> 