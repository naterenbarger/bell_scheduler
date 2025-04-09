<template>
  <div>
    <v-row>
      <v-col>
        <h1 class="text-h4 mb-4">Edit Schedule</h1>
      </v-col>
    </v-row>

    <v-card v-if="schedule" class="mb-4">
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <div class="subtitle-1 font-weight-bold">Current Schedule Details</div>
            <div class="body-2">{{ schedule.name }}</div>
            <div class="body-2">{{ schedule.description }}</div>
            <div class="body-2">Timezone: {{ schedule.timezone }}</div>
            <div class="body-2" v-if="schedule.isDefault">
              <v-chip color="primary" small>Default Schedule</v-chip>
            </div>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <schedule-form
      :schedule="schedule"
      :loading="loading"
      :error="error"
      :is-edit="true"
      @success="handleSuccess"
      @cancel="handleCancel"
    />
  </div>
</template>

<script>
import ScheduleForm from '@/components/schedules/ScheduleForm.vue'

export default {
  name: 'ScheduleEdit',
  components: {
    ScheduleForm
  },
  computed: {
    schedule() {
      return this.$store.getters['schedules/currentSchedule']
    },
    loading() {
      return this.$store.getters['schedules/loading']
    },
    error() {
      return this.$store.getters['schedules/error']
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
    handleSuccess() {
      this.$router.push('/schedules')
    },
    handleCancel() {
      this.$router.push('/schedules')
    }
  }
}
</script> 