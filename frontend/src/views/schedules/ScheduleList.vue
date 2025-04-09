<template>
  <div>
    <v-row>
      <v-col>
        <h1 class="text-h4 mb-4">Schedules</h1>
      </v-col>
      <v-col class="text-right">
        <v-btn
          color="primary"
          @click="$router.push('/schedules/create')"
        >
          <v-icon left>mdi-plus</v-icon>
          New Schedule
        </v-btn>
      </v-col>
    </v-row>
    
    <v-alert type="info" class="mb-4" outlined>
      <strong>Schedule Status:</strong>
      <ul class="mb-0">
        <li><strong>Active:</strong> Currently in use</li>
        <li><strong>Default:</strong> Becomes active at midnight</li>
        <li><strong>Temporary:</strong> Resets to default at midnight</li>
      </ul>
    </v-alert>

    <v-card>
      <v-card-text>
        <v-data-table
          :headers="headers"
          :items="schedules"
          :loading="loading"
          :items-per-page="10"
          class="elevation-1"
        >
          <template v-slot:item.isActive="{ item }">
            <v-chip
              v-if="item.isActive"
              color="success"
              small
            >
              Active
            </v-chip>
            <span v-else>Inactive</span>
          </template>
          
          <template v-slot:item.isDefault="{ item }">
            <v-chip
              v-if="item.isDefault"
              color="primary"
              small
            >
              Default
            </v-chip>
            <span v-else>-</span>
          </template>
          
          <template v-slot:item.isTemporary="{ item }">
            <v-chip
              v-if="item.isTemporary"
              color="amber"
              small
            >
              Temporary
            </v-chip>
            <span v-else>-</span>
          </template>

          <template v-slot:item.actions="{ item }">
            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-btn
                  icon
                  small
                  class="mr-2"
                  v-bind="attrs"
                  v-on="on"
                  @click="viewSchedule(item)"
                >
                  <v-icon>mdi-eye</v-icon>
                </v-btn>
              </template>
              <span>View Details</span>
            </v-tooltip>

            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-btn
                  icon
                  small
                  class="mr-2"
                  v-bind="attrs"
                  v-on="on"
                  @click="editSchedule(item)"
                >
                  <v-icon>mdi-pencil</v-icon>
                </v-btn>
              </template>
              <span>Edit Schedule</span>
            </v-tooltip>

            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-btn
                  icon
                  small
                  class="mr-2"
                  v-bind="attrs"
                  v-on="on"
                  @click="activateSchedule(item)"
                  :color="item.isDefault ? 'amber' : ''"
                >
                  <v-icon>mdi-star</v-icon>
                </v-btn>
              </template>
              <span>Set as Active Schedule</span>
            </v-tooltip>

            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-btn
                  icon
                  small
                  v-bind="attrs"
                  v-on="on"
                  @click="confirmDelete(item)"
                >
                  <v-icon>mdi-delete</v-icon>
                </v-btn>
              </template>
              <span>Delete Schedule</span>
            </v-tooltip>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h5">
          Delete Schedule
        </v-card-title>
        <v-card-text>
          Are you sure you want to delete this schedule? This action cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            text
            @click="deleteDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn
            color="error"
            :loading="loading"
            @click="deleteSchedule"
          >
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    
    <!-- Activate Schedule Dialog -->
    <v-dialog v-model="activateDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h5">
          Activate Schedule
        </v-card-title>
        <v-card-text>
          <p>Do you want to make <strong>{{ scheduleToActivate?.name }}</strong> the active schedule?</p>
          
          <v-radio-group v-model="activateMode" class="mt-4">
            <v-radio value="active" label="Set as active only" />
            <v-radio value="default" label="Set as default (becomes active at midnight)" />
            <v-radio value="temporary" label="Set as temporary (resets to default at midnight)" />
          </v-radio-group>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            text
            @click="activateDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn
            color="primary"
            :loading="loading"
            @click="confirmActivateSchedule"
          >
            Activate
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Error Snackbar -->
    <v-snackbar
      v-model="showError"
      color="error"
      timeout="5000"
    >
      {{ error }}
    </v-snackbar>

    <!-- Success Snackbar -->
    <v-snackbar
      v-model="showSuccess"
      color="success"
      timeout="3000"
    >
      {{ successMessage }}
    </v-snackbar>
  </div>
</template>

<script>
export default {
  name: 'ScheduleList',
  data: () => ({
    headers: [
      { text: 'Name', value: 'name' },
      { text: 'Description', value: 'description' },
      { text: 'Active', value: 'isActive' },
      { text: 'Default', value: 'isDefault' },
      { text: 'Temporary', value: 'isTemporary' },
      { text: 'Created', value: 'createdAt' },
      { text: 'Actions', value: 'actions', sortable: false }
    ],
    deleteDialog: false,
    scheduleToDelete: null,
    activateDialog: false,
    scheduleToActivate: null,
    activateMode: 'active',
    showError: false,
    showSuccess: false,
    successMessage: ''
  }),
  computed: {
    schedules() {
      return this.$store.getters['schedules/schedules']
    },
    loading() {
      return this.$store.getters['schedules/loading']
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
    this.fetchSchedules()
  },
  methods: {
    async fetchSchedules() {
      await this.$store.dispatch('schedules/fetchSchedules')
    },
    viewSchedule(schedule) {
      this.$router.push(`/schedules/${schedule.id}`)
    },
    editSchedule(schedule) {
      this.$router.push(`/schedules/${schedule.id}/edit`)
    },
    activateSchedule(schedule) {
      this.scheduleToActivate = schedule
      this.activateMode = 'active'
      this.activateDialog = true
    },
    
    async confirmActivateSchedule() {
      if (!this.scheduleToActivate) return
      
      try {
        switch (this.activateMode) {
          case 'active':
            // Set as active only
            await this.$store.dispatch('schedules/setActiveSchedule', this.scheduleToActivate.id)
            this.successMessage = `${this.scheduleToActivate.name} has been set as the active schedule`
            break
          case 'default':
            // Set as default
            await this.$store.dispatch('schedules/setDefaultSchedule', this.scheduleToActivate.id)
            this.successMessage = `${this.scheduleToActivate.name} has been set as the default schedule`
            break
          case 'temporary':
            // Set as temporary
            await this.$store.dispatch('schedules/setTemporarySchedule', {
              id: this.scheduleToActivate.id,
              isTemporary: true
            })
            this.successMessage = `${this.scheduleToActivate.name} has been set as the active schedule and will automatically reset at midnight`
            break
        }
        
        this.showSuccess = true
        this.activateDialog = false
        this.scheduleToActivate = null
      } catch (error) {
        console.error('Failed to set schedule as active:', error)
    },
    confirmDelete(schedule) {
      this.scheduleToDelete = schedule
      this.deleteDialog = true
    },
    async deleteSchedule() {
      if (!this.scheduleToDelete) return

      try {
        await this.$store.dispatch('schedules/deleteSchedule', this.scheduleToDelete.id)
        this.deleteDialog = false
        this.scheduleToDelete = null
      } catch (error) {
        console.error('Failed to delete schedule:', error)
      }
    }
  }
}
</script>