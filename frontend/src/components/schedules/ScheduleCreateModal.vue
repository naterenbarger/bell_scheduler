<template>
  <v-dialog v-model="showDialog" max-width="600px">
    <v-card>
      <v-card-title>
        Create New Schedule
      </v-card-title>
      <v-card-text>
        <schedule-form
          :loading="loading"
          :error="error"
          @success="handleSuccess"
          @cancel="closeDialog"
        />
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script>
import ScheduleForm from '@/components/schedules/ScheduleForm.vue'

export default {
  name: 'ScheduleCreateModal',
  components: {
    ScheduleForm
  },
  props: {
    value: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    showDialog: {
      get() {
        return this.value
      },
      set(value) {
        this.$emit('input', value)
      }
    },
    loading() {
      return this.$store.getters['schedules/loading']
    },
    error() {
      return this.$store.getters['schedules/error']
    }
  },
  methods: {
    handleSuccess() {
      this.$emit('success')
      this.closeDialog()
    },
    closeDialog() {
      this.showDialog = false
    }
  }
}
</script>