<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            Bell Schedules
            <v-spacer></v-spacer>
            <v-text-field
              v-model="search"
              append-icon="mdi-magnify"
              label="Search"
              single-line
              hide-details
            ></v-text-field>
            <v-btn color="primary" class="ml-4" @click="openDialog()">
              <v-icon left>mdi-plus</v-icon>
              New Schedule
            </v-btn>
          </v-card-title>

          <v-data-table
            :headers="headers"
            :items="schedules"
            :search="search"
            :loading="loading"
          >
            <template v-slot:item.isDefault="{ item }">
              <v-chip :color="item.isDefault ? 'primary' : ''">
                {{ item.isDefault ? 'Yes' : 'No' }}
              </v-chip>
            </template>

            <template v-slot:item.isActive="{ item }">
              <v-chip :color="item.isActive ? 'success' : ''">
                {{ item.isActive ? 'Yes' : 'No' }}
              </v-chip>
            </template>

            <template v-slot:item.isTemporary="{ item }">
              <v-chip :color="item.isTemporary ? 'amber' : ''">
                {{ item.isTemporary ? 'Yes' : 'No' }}
              </v-chip>
            </template>

            <template v-slot:item.actions="{ item }">
              <v-btn icon small class="mr-2" @click="openDialog(item)">
                <v-icon>mdi-pencil</v-icon>
              </v-btn>
              <v-btn icon small class="mr-2" @click="triggerNow(item)">
                <v-icon>mdi-bell</v-icon>
              </v-btn>
              <v-btn icon small class="mr-2" @click="setActive(item)">
                <v-icon>mdi-star</v-icon>
              </v-btn>
              <v-btn icon small color="error" @click="confirmDelete(item)">
                <v-icon>mdi-delete</v-icon>
              </v-btn>
            </template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>

    <!-- Schedule Dialog -->
    <v-dialog v-model="dialog" max-width="600px">
      <v-card>
        <v-card-title>
          {{ editedItem.id ? 'Edit Schedule' : 'New Schedule' }}
        </v-card-title>

        <v-card-text>
          <v-form ref="form" v-model="valid">
            <v-text-field
              v-model="editedItem.name"
              label="Schedule Name"
              :rules="[v => !!v || 'Name is required']"
              required
            ></v-text-field>

            <v-textarea
              v-model="editedItem.description"
              label="Description"
              rows="3"
            ></v-textarea>

            <v-switch
              v-model="editedItem.isDefault"
              label="Set as Default Schedule"
            ></v-switch>

            <v-switch
              v-model="editedItem.isActive"
              label="Set as Active Schedule"
              hint="This will make this schedule the currently active one"
              persistent-hint
            ></v-switch>

            <v-switch
              v-model="editedItem.isTemporary"
              label="Temporary Schedule (resets at midnight)"
              hint="Schedule will automatically revert to default at the end of the day"
              persistent-hint
            ></v-switch>

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
                    <v-list-item v-for="(slot, i) in editedItem.timeSlots" :key="i">
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
          </v-form>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="closeDialog">Cancel</v-btn>
          <v-btn color="primary" :disabled="!valid" @click="save">Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="deleteDialog" max-width="400px">
      <v-card>
        <v-card-title>Delete Schedule</v-card-title>
        <v-card-text>
          Are you sure you want to delete this schedule?
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="deleteDialog = false">Cancel</v-btn>
          <v-btn color="error" text @click="deleteSchedule">Delete</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Set Active Dialog -->
    <v-dialog v-model="activeDialog" max-width="400px">
      <v-card>
        <v-card-title>Set Active Schedule</v-card-title>
        <v-card-text>
          <p>Do you want to make <strong>{{ editedItem.name }}</strong> the active schedule?</p>
          <p>This will replace the current active schedule.</p>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="activeDialog = false">Cancel</v-btn>
          <v-btn color="primary" text @click="confirmSetActive">Confirm</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
import axios from '../plugins/axios'

export default {
  name: 'Schedules',
  data() {
    return {
      search: '',
      loading: false,
      dialog: false,
      deleteDialog: false,
      activeDialog: false,
      valid: true,
      headers: [
        { text: 'Name', value: 'name' },
        { text: 'Description', value: 'description' },
        { text: 'Default', value: 'isDefault' },
        { text: 'Active', value: 'isActive' },
        { text: 'Temporary', value: 'isTemporary' },
        { text: 'Actions', value: 'actions', sortable: false }
      ],
      schedules: [],
      editedItem: {
        id: null,
        name: '',
        description: '',
        isDefault: false,
        timeSlots: []
      },
      defaultItem: {
        id: null,
        name: '',
        description: '',
        isDefault: false,
        isActive: false,
        isTemporary: false,
        timeSlots: [{
          triggerTime: '',
          days: []
        }]
      },
      weekDays: [
        'Monday',
        'Tuesday',
        'Wednesday',
        'Thursday',
        'Friday',
        'Saturday',
        'Sunday'
      ]
    }
  },
  created() {
    this.fetchSchedules()
  },
  methods: {
    async fetchSchedules() {
      this.loading = true
      try {
        const response = await axios.get('/schedules')
        console.log('Raw response:', response.data) // Log raw response
        // Convert days JSON strings back to arrays for each time slot
        this.schedules = response.data.map(schedule => {
          console.log('Processing schedule:', schedule) // Log each schedule
          return {
            ...schedule,
            timeSlots: schedule.timeSlots.map(slot => {
              console.log('Processing time slot:', slot) // Log each time slot
              return {
                ...slot,
                triggerTime: slot.triggerTime || '', // Ensure triggerTime is preserved
                days: JSON.parse(slot.days || '[]')
              }
            })
          }
        })
        console.log('Final processed schedules:', this.schedules) // Log final result
      } catch (error) {
        console.error('Error fetching schedules:', error) // Log any errors
        this.$store.dispatch('snackbar/show', {
          message: 'Failed to fetch schedules',
          color: 'error'
        })
      } finally {
        this.loading = false
      }
    },
    openDialog(item) {
      if (item) {
        // Create a deep copy of the item including timeSlots
        this.editedItem = {
          ...item,
          timeSlots: item.timeSlots.map(slot => ({
            ...slot,
            // Ensure days is properly copied as an array
            days: Array.isArray(slot.days) ? [...slot.days] : JSON.parse(slot.days || '[]'),
            // Ensure triggerTime is preserved
            triggerTime: slot.triggerTime || ''
          }))
        }
        console.log('Opening edit dialog with timeslots:', this.editedItem.timeSlots)
      } else {
        this.editedItem = { ...this.defaultItem }
      }
      this.dialog = true
    },
    closeDialog() {
      this.dialog = false
      this.$nextTick(() => {
        this.editedItem = { ...this.defaultItem }
        this.$refs.form.reset()
      })
    },
    addTimeSlot() {
      this.editedItem.timeSlots.push({
        id: 0, // Set ID to 0 for new slots
        triggerTime: '',
        days: []
      })
    },
    removeTimeSlot(index) {
      this.editedItem.timeSlots.splice(index, 1)
    },
    async save() {
      if (!this.$refs.form.validate()) return

      try {
        // Create a copy of the edited item to modify
        const scheduleToSave = { ...this.editedItem }
        
        // If isActive is true, we need to handle it specially
        const isActivating = scheduleToSave.isActive
        
        // Remove isActive from the schedule object as it's handled separately
        delete scheduleToSave.isActive
        
        // Convert days arrays to JSON strings for each time slot
        scheduleToSave.timeSlots = scheduleToSave.timeSlots.map(slot => ({
          ...slot,
          id: slot.id || 0, // Ensure ID is 0 for new slots
          days: JSON.stringify(slot.days || []),
          triggerTime: slot.triggerTime || '', // Ensure triggerTime is preserved
          scheduleId: scheduleToSave.id // Ensure scheduleId is set
        }))

        console.log('Saving schedule:', scheduleToSave) // Add logging

        let savedScheduleId = this.editedItem.id
        
        if (this.editedItem.id) {
          await axios.put(`/schedules/${this.editedItem.id}`, scheduleToSave)
        } else {
          const response = await axios.post('/schedules', scheduleToSave)
          savedScheduleId = response.data.id
        }
        
        // If isActive was set to true, set this schedule as active
        if (isActivating && savedScheduleId) {
          await axios.put(`/schedules/${savedScheduleId}/active`)
        }
        this.closeDialog()
        this.fetchSchedules()
        this.$store.dispatch('snackbar/show', {
          message: 'Schedule saved successfully',
          color: 'success'
        })
      } catch (error) {
        console.error('Error saving schedule:', error) // Add error logging
        this.$store.dispatch('snackbar/show', {
          message: error.response?.data?.error || 'Failed to save schedule',
          color: 'error'
        })
      }
    },
    confirmDelete(item) {
      this.editedItem = { ...item }
      this.deleteDialog = true
    },
    async deleteSchedule() {
      try {
        await axios.delete(`/schedules/${this.editedItem.id}`)
        this.deleteDialog = false
        this.fetchSchedules()
        this.$store.dispatch('snackbar/show', {
          message: 'Schedule deleted successfully',
          color: 'success'
        })
      } catch (error) {
        this.$store.dispatch('snackbar/show', {
          message: error.response?.data?.error || 'Failed to delete schedule',
          color: 'error'
        })
      }
    },
    async triggerNow(item) {
      try {
        await axios.post(`/schedules/${item.id}/trigger`)
        this.$store.dispatch('snackbar/show', {
          message: 'Bell triggered successfully',
          color: 'success'
        })
      } catch (error) {
        this.$store.dispatch('snackbar/show', {
          message: error.response?.data?.error || 'Failed to trigger bell',
          color: 'error'
        })
      }
    },
    
    setActive(item) {
      this.editedItem = { ...item }
      this.activeDialog = true
    },
    
    async confirmSetActive() {
      try {
        // Set schedule as active
        await axios.put(`/schedules/${this.editedItem.id}/active`)
        
        this.activeDialog = false
        this.fetchSchedules()
        this.$store.dispatch('snackbar/show', {
          message: `${this.editedItem.name} set as active schedule successfully`,
          color: 'success'
        })
      } catch (error) {
        this.$store.dispatch('snackbar/show', {
          message: error.response?.data?.error || 'Failed to set active schedule',
          color: 'error'
        })
      }
    }
  }
}
</script>