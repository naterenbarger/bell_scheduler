<template>
  <v-row>
    <v-col cols="12" sm="6" md="6">
      <v-card>
        <v-card-text>
          <div class="text-h6 mb-2">Default Schedule</div>
          <div class="text-h4">{{ defaultSchedule ? defaultSchedule.name : 'None' }}</div>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col cols="12" sm="6" md="6">
      <v-card>
        <v-card-text>
          <div class="text-h6 mb-2">Next Bell</div>
          <div class="text-h4">{{ nextBellTime || 'No upcoming bells' }}</div>
          <div class="text-caption" v-if="nextBellTime">
            {{ nextBellDate }}
          </div>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
export default {
  name: 'DashboardStats',
  computed: {
    schedules() {
      return this.$store.getters['schedules/schedules']
    },
    defaultSchedule() {
      return this.$store.getters['schedules/defaultSchedule']
    },
    activeSchedule() {
      return this.$store.getters['schedules/activeSchedule']
    },
    settings() {
      return this.$store.getters['settings/settings']
    },
    // Get the schedule to use for next bell calculation
    // First check for active schedule, then fall back to default schedule
    scheduleForNextBell() {
      return this.activeSchedule || this.defaultSchedule
    },
    nextBellTime() {
      if (!this.scheduleForNextBell || !this.scheduleForNextBell.timeSlots) return null
      
      const now = new Date()
      const currentDay = now.getDay() // 0 = Sunday, 1 = Monday, etc.
      const dayNames = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']
      const currentDayName = dayNames[currentDay]
      
      // Current time in minutes since midnight
      const currentHour = now.getHours()
      const currentMinute = now.getMinutes()
      const currentTimeInMinutes = currentHour * 60 + currentMinute
      
      // Filter slots for today
      const todaySlots = this.scheduleForNextBell.timeSlots.filter(slot => {
        const slotDays = JSON.parse(slot.days || '[]')
        return slotDays.includes(currentDayName)
      })
      
      // Find the next slot for today
      let nextSlot = todaySlots.find(slot => {
        const [hours, minutes] = slot.triggerTime.split(':').map(Number)
        const slotTimeInMinutes = hours * 60 + minutes
        return slotTimeInMinutes > currentTimeInMinutes
      })
      
      // If no next slot today, look for the first slot on the next day with bells
      if (!nextSlot) {
        // Check each subsequent day
        for (let i = 1; i <= 7; i++) {
          const nextDay = (currentDay + i) % 7
          const nextDayName = dayNames[nextDay]
          
          // Find slots for the next day
          const nextDaySlots = this.scheduleForNextBell.timeSlots.filter(slot => {
            const slotDays = JSON.parse(slot.days || '[]')
            return slotDays.includes(nextDayName)
          })
          
          // Sort by time
          nextDaySlots.sort((a, b) => {
            return a.triggerTime.localeCompare(b.triggerTime)
          })
          
          // If we found slots for the next day, use the first one
          if (nextDaySlots.length > 0) {
            nextSlot = nextDaySlots[0]
            break
          }
        }
      }
      
      // Format the time if we found a next slot
      if (nextSlot) {
        const [hours, minutes] = nextSlot.triggerTime.split(':').map(Number)
        const hour12 = hours % 12 || 12
        const ampm = hours >= 12 ? 'PM' : 'AM'
        return `${hour12}:${minutes.toString().padStart(2, '0')} ${ampm}`
      }
      
      return null
    },
    nextBellDate() {
      if (!this.scheduleForNextBell || !this.scheduleForNextBell.timeSlots) return null
      
      const now = new Date()
      const currentDay = now.getDay()
      const dayNames = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']
      const currentDayName = dayNames[currentDay]
      
      // Current time in minutes since midnight
      const currentHour = now.getHours()
      const currentMinute = now.getMinutes()
      const currentTimeInMinutes = currentHour * 60 + currentMinute
      
      // Filter slots for today
      const todaySlots = this.scheduleForNextBell.timeSlots.filter(slot => {
        const slotDays = JSON.parse(slot.days || '[]')
        return slotDays.includes(currentDayName)
      })
      
      // Find the next slot for today
      let nextSlot = todaySlots.find(slot => {
        const [hours, minutes] = slot.triggerTime.split(':').map(Number)
        const slotTimeInMinutes = hours * 60 + minutes
        return slotTimeInMinutes > currentTimeInMinutes
      })
      
      // If we found a slot for today
      if (nextSlot) {
        return 'Today'
      }
      
      // If no next slot today, look for the first slot on the next day with bells
      for (let i = 1; i <= 7; i++) {
        const nextDay = (currentDay + i) % 7
        const nextDayName = dayNames[nextDay]
        
        // Find slots for the next day
        const nextDaySlots = this.scheduleForNextBell.timeSlots.filter(slot => {
          const slotDays = JSON.parse(slot.days || '[]')
          return slotDays.includes(nextDayName)
        })
        
        // If we found slots for the next day
        if (nextDaySlots.length > 0) {
          // If it's tomorrow
          if (i === 1) {
            return 'Tomorrow'
          } else {
            // Return the day name
            return nextDayName
          }
        }
      }
      
      return null
    }
  }
}
</script>