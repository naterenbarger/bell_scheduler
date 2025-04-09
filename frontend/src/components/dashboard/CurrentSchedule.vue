<template>
  <v-card>
    <v-card-title class="d-flex align-center">
      <span>Current Schedule</span>
      <v-spacer />
      <v-chip
        v-if="activeSchedule && activeSchedule.isTemporary"
        color="amber"
        small
        class="mr-2"
      >
        Temporary
      </v-chip>
      <v-chip
        v-if="activeSchedule"
        color="success"
        small
        class="mr-2"
      >
        Active
      </v-chip>
      <v-chip
        v-if="activeSchedule && activeSchedule.isDefault"
        color="primary"
        small
      >
        Default
      </v-chip>
    </v-card-title>

    <v-card-text v-if="activeSchedule">
      <v-row>
        <v-col cols="12" md="6">
          <div class="subtitle-1 font-weight-bold mb-2">{{ activeSchedule.name }}</div>
          <div class="body-2 mb-4">{{ activeSchedule.description }}</div>
          <div v-if="activeSchedule.isTemporary" class="body-2 mt-2">
            <v-icon small class="mr-1">mdi-clock-alert-outline</v-icon>
            <i>Temporary: Will reset to default schedule at midnight</i>
          </div>
          <div v-else-if="activeSchedule.isDefault" class="body-2 mt-2">
            <v-icon small class="mr-1">mdi-check-circle-outline</v-icon>
            <i>This is also the default schedule</i>
          </div>
        </v-col>
      </v-row>

      <v-divider class="my-4" />

      <div class="subtitle-1 font-weight-bold mb-2">Today's Schedule</div>
      <v-list dense>
        <v-list-item v-for="slot in todaySlots" :key="slot.id">
          <v-list-item-content>
            <v-list-item-title>{{ formatTime(slot.triggerTime) }}</v-list-item-title>
            <v-list-item-subtitle>{{ slot.description || 'No description' }}</v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>

        <v-list-item v-if="!todaySlots.length">
          <v-list-item-content>
            <v-list-item-title class="text-center">No bells scheduled for today</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-card-text>

    <v-card-text v-else>
      <div class="text-center pa-4">
        <v-icon size="48" color="grey" class="mb-2">mdi-calendar-off</v-icon>
        <div class="text-h6">No Active Schedule</div>
        <div class="body-2 mt-2">Create a schedule and set it as default to start managing bell times.</div>
        <v-btn
          color="primary"
          class="mt-4"
          @click="$emit('create-schedule')"
        >
          Create Schedule
        </v-btn>
      </div>
    </v-card-text>
  </v-card>
</template>

<script>
export default {
  name: 'CurrentSchedule',
  computed: {
    activeSchedule() {
      return this.$store.getters['schedules/activeSchedule']
    },
    todaySlots() {
      if (!this.activeSchedule || !this.activeSchedule.timeSlots) return []
      
      // Get current day of week (0 = Sunday, 1 = Monday, etc.)
      const today = new Date().getDay()
      // Convert to day name to match our data format
      const dayNames = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']
      const todayName = dayNames[today]
      
      // Filter timeslots for today
      return this.activeSchedule.timeSlots.filter(slot => {
        // Parse the days JSON string to an array
        const slotDays = JSON.parse(slot.days || '[]')
        // Check if today is included in the slot's days
        return slotDays.includes(todayName)
      }).sort((a, b) => {
        // Sort by time
        return a.triggerTime.localeCompare(b.triggerTime)
      })
    }
  },
  methods: {
    formatTime(time) {
      if (!time) return ''
      
      // Assuming time is in HH:MM 24-hour format
      const [hours, minutes] = time.split(':')
      
      // Convert to 12-hour format with AM/PM
      const hour = parseInt(hours, 10)
      const ampm = hour >= 12 ? 'PM' : 'AM'
      const hour12 = hour % 12 || 12
      
      return `${hour12}:${minutes} ${ampm}`
    }
  }
}
</script>