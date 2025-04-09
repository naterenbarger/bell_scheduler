<template>
  <div>
    <v-row>
      <v-col>
        <h1 class="text-h4 mb-4">Dashboard</h1>
      </v-col>
    </v-row>

    <dashboard-stats class="mb-6" />

    <v-row>
      <v-col cols="12" md="8">
        <current-schedule class="mb-6" @create-schedule="openScheduleModal" />
      </v-col>
      <v-col cols="12" md="4">
        <quick-actions class="mb-6" @create-schedule="openScheduleModal" />
      </v-col>
    </v-row>

    <!-- Schedule Create Modal -->
    <schedule-create-modal 
      v-model="showScheduleModal"
      @success="handleScheduleCreated"
    />

    <!-- Error Snackbar -->
    <v-snackbar
      v-model="showError"
      color="error"
      timeout="5000"
    >
      {{ error }}
    </v-snackbar>
  </div>
</template>

<script>
import DashboardStats from '@/components/dashboard/DashboardStats.vue'
import QuickActions from '@/components/dashboard/QuickActions.vue'
import CurrentSchedule from '@/components/dashboard/CurrentSchedule.vue'
import ScheduleCreateModal from '@/components/schedules/ScheduleCreateModal.vue'

export default {
  name: 'Dashboard',
  components: {
    DashboardStats,
    QuickActions,
    CurrentSchedule,
    ScheduleCreateModal
  },
  data: () => ({
    showError: false,
    showScheduleModal: false
  }),
  computed: {
    error() {
      return this.$store.getters['schedules/error'] || this.$store.getters['settings/error']
    }
  },
  watch: {
    error(newError) {
      this.showError = !!newError
    }
  },
  created() {
    this.initializeData()
  },
  methods: {
    async initializeData() {
      try {
        await Promise.all([
          this.$store.dispatch('schedules/fetchSchedules'),
          this.$store.dispatch('settings/fetchSettings')
        ])
      } catch (error) {
        console.error('Failed to initialize dashboard data:', error)
      }
    },
    openScheduleModal() {
      this.showScheduleModal = true
    },
    handleScheduleCreated() {
      // Refresh schedules data after a new schedule is created
      this.$store.dispatch('schedules/fetchSchedules')
    }
  }
}
</script>