<template>
  <div>
    <v-row>
      <v-col>
        <h1 class="text-h4 mb-4">Bell Ringing Logs</h1>
      </v-col>
      <v-col cols="auto">
        <v-btn
          color="primary"
          @click="refreshLogs"
          :loading="loading"
        >
          <v-icon left>mdi-refresh</v-icon>
          Refresh
        </v-btn>
      </v-col>
    </v-row>

    <v-card>
      <v-card-text>
        <v-row>
          <v-col cols="12" sm="6" md="3">
            <v-menu
              ref="startMenu"
              v-model="startMenu"
              :close-on-content-click="false"
              transition="scale-transition"
              offset-y
            >
              <template v-slot:activator="{ on, attrs }">
                <v-text-field
                  v-model="startDate"
                  label="Start Date"
                  readonly
                  v-bind="attrs"
                  v-on="on"
                ></v-text-field>
              </template>
              <v-date-picker
                v-model="startDate"
                @change="startMenu = false"
                :max="endDate"
              ></v-date-picker>
            </v-menu>
          </v-col>
          <v-col cols="12" sm="6" md="3">
            <v-menu
              ref="endMenu"
              v-model="endMenu"
              :close-on-content-click="false"
              transition="scale-transition"
              offset-y
            >
              <template v-slot:activator="{ on, attrs }">
                <v-text-field
                  v-model="endDate"
                  label="End Date"
                  readonly
                  v-bind="attrs"
                  v-on="on"
                ></v-text-field>
              </template>
              <v-date-picker
                v-model="endDate"
                @change="endMenu = false"
                :min="startDate"
              ></v-date-picker>
            </v-menu>
          </v-col>
          <v-col cols="12" sm="6" md="3">
            <v-select
              v-model="triggerFilter"
              :items="triggerOptions"
              label="Trigger Type"
              clearable
              @change="filterLogs"
            ></v-select>
          </v-col>
          <v-col cols="12" sm="6" md="3">
            <v-btn
              color="primary"
              block
              @click="filterLogs"
            >
              Filter
            </v-btn>
          </v-col>
        </v-row>

        <v-data-table
          :headers="headers"
          :items="filteredLogs"
          :loading="loading"
          :items-per-page="10"
          class="elevation-1"
        >
          <template v-slot:item.timestamp="{ item }">
            {{ formatDate(item.timestamp) }}
          </template>
          <template v-slot:item.trigger="{ item }">
            <v-chip
              :color="item.trigger === 'schedule' ? 'primary' : 'success'"
              small
            >
              {{ item.trigger }}
            </v-chip>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import { format } from 'date-fns'

export default {
  name: 'Logs',
  data: () => ({
    loading: false,
    logs: [],
    startMenu: false,
    endMenu: false,
    startDate: null,
    endDate: null,
    triggerFilter: null,
    triggerOptions: ['schedule', 'manual'],
    headers: [
      { text: 'Timestamp', value: 'timestamp' },
      { text: 'Trigger', value: 'trigger' },
      { text: 'User', value: 'username' },
      { text: 'Schedule', value: 'scheduleName' },
      { text: 'Time', value: 'scheduleTime' }
    ]
  }),
  computed: {
    filteredLogs() {
      return this.logs.filter(log => {
        if (this.triggerFilter && log.trigger !== this.triggerFilter) {
          return false
        }
        return true
      })
    }
  },
  methods: {
    formatDate(date) {
      return format(new Date(date), 'PPpp')
    },
    async fetchLogs() {
      try {
        this.loading = true
        const response = await this.$axios.get('/logs')
        this.logs = response.data
      } catch (error) {
        console.error('Failed to fetch logs:', error)
        this.$store.dispatch('notifications/showError', 'Failed to fetch logs')
      } finally {
        this.loading = false
      }
    },
    async filterLogs() {
      if (!this.startDate || !this.endDate) {
        await this.fetchLogs()
        return
      }

      try {
        this.loading = true
        const start = new Date(this.startDate)
        const end = new Date(this.endDate)
        end.setHours(23, 59, 59, 999)

        const response = await this.$axios.get('/logs/range', {
          params: {
            start: start.toISOString(),
            end: end.toISOString()
          }
        })
        this.logs = response.data
      } catch (error) {
        console.error('Failed to filter logs:', error)
        this.$store.dispatch('notifications/showError', 'Failed to filter logs')
      } finally {
        this.loading = false
      }
    },
    refreshLogs() {
      this.startDate = null
      this.endDate = null
      this.triggerFilter = null
      this.fetchLogs()
    }
  },
  created() {
    this.fetchLogs()
  }
}
</script>