<template>
  <v-form ref="form" v-model="valid">
    <v-card>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="formData.name"
              :rules="[v => !!v || 'Name is required']"
              label="Schedule Name"
              required
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="formData.description"
              label="Description"
              :rules="[v => !!v || 'Description is required']"
              required
            />
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12">
            <v-select
              v-model="formData.timezone"
              :items="timezones"
              label="Timezone"
              :rules="[v => !!v || 'Timezone is required']"
              required
            />
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" md="4">
            <v-switch
              v-model="formData.isActive"
              label="Set as Active Schedule"
              hint="Currently in use schedule"
              persistent-hint
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-switch
              v-model="formData.isDefault"
              label="Set as Default Schedule"
              hint="Becomes active at midnight"
              persistent-hint
              :disabled="isDefault && !isEdit"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-switch
              v-model="formData.isTemporary"
              label="Temporary Schedule"
              hint="Resets to default at midnight"
              persistent-hint
            />
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12">
            <v-textarea
              v-model="formData.notes"
              label="Notes"
              rows="3"
              auto-grow
            />
          </v-col>
        </v-row>
        
        <v-row>
          <v-col cols="12">
            <v-expansion-panels>
              <v-expansion-panel>
                <v-expansion-panel-header>
                  Time Slots
                </v-expansion-panel-header>
                <v-expansion-panel-content>
                  <v-btn color="primary" text @click="addTimeSlot">
                    <v-icon left>mdi-plus</v-icon>
                    Add Time Slot
                  </v-btn>

                  <v-list>
                    <v-list-item v-for="(slot, i) in formData.timeSlots" :key="i">
                      <v-list-item-content>
                        <v-row align="center">
                          <v-col cols="4">
                            <v-text-field
                              v-model="slot.triggerTime"
                              label="Time (HH:MM)"
                              type="time"
                              required
                            ></v-text-field>
                          </v-col>
                          <v-col cols="6">
                            <v-select
                              v-model="slot.days"
                              :items="weekDays"
                              label="Days"
                              multiple
                              chips
                              required
                            ></v-select>
                          </v-col>
                          <v-col cols="2">
                            <v-btn icon @click="removeTimeSlot(i)">
                              <v-icon>mdi-delete</v-icon>
                            </v-btn>
                          </v-col>
                        </v-row>
                      </v-list-item-content>
                    </v-list-item>
                  </v-list>
                </v-expansion-panel-content>
              </v-expansion-panel>
            </v-expansion-panels>
          </v-col>
        </v-row>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn
          text
          @click="$emit('cancel')"
        >
          Cancel
        </v-btn>
        <v-btn
          color="primary"
          :loading="loading"
          :disabled="!valid"
          @click="handleSubmit"
        >
          {{ submitButtonText }}
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
  </v-form>
</template>

<script>
import { timezones } from '@/utils/timezones'

export default {
  name: 'ScheduleForm',
  props: {
    schedule: {
      type: Object,
      default: () => ({})
    },
    isEdit: {
      type: Boolean,
      default: false
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
      name: '',
      description: '',
      timezone: '',
      isDefault: false,
      isTemporary: false,
      isActive: false,
      notes: '',
      timeSlots: []
    },
    timezones,
    weekDays: [
      'Monday',
      'Tuesday',
      'Wednesday',
      'Thursday',
      'Friday',
      'Saturday',
      'Sunday'
    ],
    showError: false
  }),
  computed: {
    submitButtonText() {
      return this.isEdit ? 'Update Schedule' : 'Create Schedule'
    },
    isDefault() {
      return this.schedule.isDefault
    }
  },
  watch: {
    schedule: {
      handler(newSchedule) {
        if (newSchedule && Object.keys(newSchedule).length > 0) {
          // Create a deep copy of the schedule including timeSlots
          this.formData = {
            ...newSchedule,
            timeSlots: newSchedule.timeSlots ? newSchedule.timeSlots.map(slot => ({
              ...slot,
              // Ensure days is properly copied as an array
              days: Array.isArray(slot.days) ? [...slot.days] : JSON.parse(slot.days || '[]'),
              // Ensure triggerTime is preserved
              triggerTime: slot.triggerTime || ''
            })) : []
          }
        } else {
          // Initialize with at least one empty timeslot for new schedules
          if (!this.isEdit && (!this.formData.timeSlots || this.formData.timeSlots.length === 0)) {
            this.formData.timeSlots = [{
              id: 0,
              triggerTime: '',
              days: []
            }]
          }
        }
      },
      immediate: true
    },
    error(newError) {
      this.showError = !!newError
    }
  },
  methods: {
    async handleSubmit() {
      if (!this.$refs.form.validate()) return

      try {
        // Prepare timeslots data - convert days arrays to JSON strings
        const preparedData = {
          ...this.formData,
          timeSlots: this.formData.timeSlots.map(slot => ({
            ...slot,
            id: slot.id || 0, // Ensure ID is 0 for new slots
            days: JSON.stringify(slot.days || []),
            triggerTime: slot.triggerTime || ''
          }))
        }

        if (this.isEdit) {
          await this.$store.dispatch('schedules/updateSchedule', {
            id: this.schedule.id,
            data: preparedData
          })
        } else {
          await this.$store.dispatch('schedules/createSchedule', preparedData)
        }
        this.$emit('success')
      } catch (error) {
        console.error('Failed to save schedule:', error)
      }
    },
    
    addTimeSlot() {
      this.formData.timeSlots.push({
        id: 0, // Set ID to 0 for new slots
        triggerTime: '',
        days: []
      })
    },
    
    removeTimeSlot(index) {
      this.formData.timeSlots.splice(index, 1)
    }
  }
}
</script>