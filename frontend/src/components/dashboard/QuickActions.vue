<template>
  <v-card>
    <v-card-title>Quick Actions</v-card-title>
    <v-card-text>
      <v-row justify="start" class="mx-n2">
        <v-col cols="12" sm="6" md="6" lg="4" xl="4" class="px-2 mb-3">
          <v-btn
            block
            color="primary"
            @click="$emit('create-schedule')"
            height="48"
          >
            <v-icon left>mdi-plus</v-icon>
            New Schedule
          </v-btn>
        </v-col>

        <v-col cols="12" sm="6" md="6" lg="4" xl="4" class="px-2 mb-3">
          <v-btn
            block
            color="primary"
            @click="$router.push('/schedules')"
            height="48"
          >
            <v-icon left>mdi-calendar</v-icon>
            Schedules
          </v-btn>
        </v-col>

        <v-col cols="12" sm="6" md="6" lg="4" xl="4" class="px-2 mb-3">
          <v-btn
            block
            color="primary"
            @click="$router.push('/settings')"
            height="48"
          >
            <v-icon left>mdi-cog</v-icon>
            Settings
          </v-btn>
        </v-col>

        <v-col cols="12" sm="6" md="6" lg="4" xl="4" class="px-2 mb-3">
          <v-btn
            block
            color="primary"
            @click="$router.push('/users')"
            v-if="isAdmin"
            height="48"
          >
            <v-icon left>mdi-account-group</v-icon>
            Manage Users
          </v-btn>
        </v-col>

        <v-col cols="12" sm="6" md="6" lg="4" xl="4" class="px-2 mb-3">
          <v-btn
            block
            color="success"
            :loading="isRinging"
            :disabled="isRinging"
            @click="triggerBell"
            height="48"
          >
            <v-icon left>mdi-bell-ring</v-icon>
            Ring Bell Now
          </v-btn>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script>
export default {
  name: 'QuickActions',
  data: () => ({
    isRinging: false
  }),
  computed: {
    isAdmin() {
      return this.$store.getters['auth/isAdmin']
    }
  },
  methods: {
    async triggerBell() {
      try {
        this.isRinging = true
        await this.$store.dispatch('schedules/triggerBell')
        this.$store.dispatch('notifications/showSuccess', 'Bell triggered successfully')
      } catch (error) {
        this.$store.dispatch('notifications/showError', 'Failed to trigger bell')
        console.error('Failed to trigger bell:', error)
      } finally {
        this.isRinging = false
      }
    }
  }
}
</script>