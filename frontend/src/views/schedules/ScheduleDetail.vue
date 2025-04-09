<template>
  <div>
    <v-row>
      <v-col>
        <h1 class="text-h4 mb-4">Schedule Details</h1>
      </v-col>
      <v-col class="text-right">
        <v-btn
          color="primary"
          @click="$router.push(`/schedules/${schedule.id}/edit`)"
        >
          <v-icon left>mdi-pencil</v-icon>
          Edit Schedule
        </v-btn>
      </v-col>
    </v-row>

    <v-card v-if="schedule" class="mb-4">
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <div class="subtitle-1 font-weight-bold">Schedule Information</div>
            <v-list dense>
              <v-list-item>
                <v-list-item-content>
                  <v-list-item-title>Name</v-list-item-title>
                  <v-list-item-subtitle>{{ schedule.name }}</v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>

              <v-list-item>
                <v-list-item-content>
                  <v-list-item-title>Description</v-list-item-title>
                  <v-list-item-subtitle>{{ schedule.description }}</v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>

              <v-list-item>
                <v-list-item-content>
                  <v-list-item-title>Timezone</v-list-item-title>
                  <v-list-item-subtitle>{{ schedule.timezone }}</v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>

              <v-list-item>
                <v-list-item-content>
                  <v-list-item-title>Status</v-list-item-title>
                  <v-list-item-subtitle>
                    <v-chip
                      v-if="schedule.isDefault"
                      color="primary"
                      small
                    >
                      Default Schedule
                    </v-chip>
                    <span v-else>Regular Schedule</span>
                  </v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>

              <v-list-item>
                <v-list-item-content>
                  <v-list-item-title>Created</v-list-item-title>
                  <v-list-item-subtitle>{{ formatDate(schedule.createdAt) }}</v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>

              <v-list-item>
                <v-list-item-content>
                  <v-list-item-title>Last Updated</v-list-item-title>
                  <v-list-item-subtitle>{{ formatDate(schedule.updatedAt) }}</v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>
            </v-list>
          </v-col>
        </v-row>

        <v-row v-if="schedule.notes">
          <v-col cols="12">
            <div class="subtitle-1 font-weight-bold">Notes</div>
            <div class="body-2">{{ schedule.notes }}</div>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

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
export default {
  name: 'ScheduleDetail',
  data: () => ({
    showError: false
  }),
  computed: {
    schedule() {
      return this.$store.getters['schedules/currentSchedule']
    },
    error() {
      return this.$store.getters['schedules/error']
    }
  },
  watch: {
    error(newError) {
      this.showError = !!newError
    }
  },
  created() {
    this.fetchSchedule()
  },
  methods: {
    async fetchSchedule() {
      const scheduleId = this.$route.params.id
      await this.$store.dispatch('schedules/fetchSchedule', scheduleId)
    },
    formatDate(date) {
      if (!date) return ''
      return new Date(date).toLocaleString()
    }
  }
}
</script> 