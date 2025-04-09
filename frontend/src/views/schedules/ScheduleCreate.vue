<template>
  <div>
    <v-row>
      <v-col>
        <h1 class="text-h4 mb-4">Create New Schedule</h1>
      </v-col>
    </v-row>

    <schedule-form
      :loading="loading"
      :error="error"
      @success="handleSuccess"
      @cancel="handleCancel"
    />
  </div>
</template>

<script>
import ScheduleForm from '@/components/schedules/ScheduleForm.vue'

export default {
  name: 'ScheduleCreate',
  components: {
    ScheduleForm
  },
  computed: {
    loading() {
      return this.$store.getters['schedules/loading']
    },
    error() {
      return this.$store.getters['schedules/error']
    }
  },
  methods: {
    handleSuccess() {
      this.$router.push('/schedules')
    },
    handleCancel() {
      this.$router.push('/schedules')
    }
  }
}
</script> 