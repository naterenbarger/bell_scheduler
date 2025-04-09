<template>
  <div>
    <v-row>
      <v-col>
        <h1 class="text-h4 mb-4">Global Settings</h1>
      </v-col>
    </v-row>

    <settings-form
      :settings="settings"
      :loading="loading"
      :error="error"
    />
  </div>
</template>

<script>
import SettingsForm from '@/components/settings/SettingsForm.vue'

export default {
  name: 'Settings',
  components: {
    SettingsForm
  },
  computed: {
    settings() {
      return this.$store.getters['settings/settings']
    },
    loading() {
      return this.$store.getters['settings/loading']
    },
    error() {
      return this.$store.getters['settings/error']
    }
  },
  created() {
    this.fetchSettings()
  },
  methods: {
    async fetchSettings() {
      await this.$store.dispatch('settings/fetchSettings')
    }
  }
}
</script> 