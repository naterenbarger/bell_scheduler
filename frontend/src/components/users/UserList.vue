<template>
  <v-card>
    <v-data-table
      :headers="headers"
      :items="processedUsers"
      :loading="loading"
      class="elevation-1"
      dense
      :items-per-page="itemsPerPage"
      :page="currentPage"
      :sort-by="sortBy"
      :sort-desc="sortDesc"
      :search="search"
      :server-items-length="totalUsers"
      :options.sync="options"
      @update:options="handleOptionsUpdate"
    >
      <template v-slot:top>
        <v-toolbar flat>
          <v-toolbar-title>Users</v-toolbar-title>
          <v-divider class="mx-4" inset vertical></v-divider>
          <v-spacer></v-spacer>
          <v-text-field
            v-model="search"
            append-icon="mdi-magnify"
            label="Search"
            single-line
            hide-details
            @input="handleSearch"
          ></v-text-field>
        </v-toolbar>
      </template>

      <template v-slot:item.role="{ item }">
        <v-chip
          :color="item.role === 'admin' ? 'primary' : 'grey'"
          small
        >
          {{ item.role }}
        </v-chip>
      </template>

      <template v-slot:item.isActive="{ item }">
        <v-chip
          :color="item.isActive ? 'success' : 'error'"
          small
        >
          {{ item.isActive ? 'Active' : 'Inactive' }}
        </v-chip>
      </template>

      <template v-slot:item.actions="{ item }">
        <v-icon small class="mr-2" @click="$emit('edit', item)">
          mdi-pencil
        </v-icon>
        <v-icon small @click="confirmDelete(item)">
          mdi-delete
        </v-icon>
      </template>

      <template v-slot:no-data>
        <v-text>No users found</v-text>
      </template>
    </v-data-table>

    <!-- Delete Confirmation Dialog -->
    <v-dialog
      v-model="deleteDialog"
      max-width="400"
    >
      <v-card>
        <v-card-title class="headline">
          Delete User
        </v-card-title>

        <v-card-text>
          Are you sure you want to delete this user? This action cannot be undone.
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
            @click="handleDelete"
          >
            Delete
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
  </v-card>
</template>

<script>
export default {
  name: 'UserList',
  props: {
    loading: {
      type: Boolean,
      default: false
    },
    error: {
      type: String,
      default: ''
    },
    users: {
      type: Array,
      default: () => []
    },
    totalUsers: {
      type: Number,
      default: 0
    }
  },
  data: () => ({
    search: '',
    deleteDialog: false,
    showError: false,
    showSuccess: false,
    userToDelete: null,
    successMessage: '',
    itemsPerPage: 10,
    currentPage: 1,
    sortBy: ['username'],
    sortDesc: [false],
    options: {},
    headers: [
      { text: 'Username', value: 'username' },
      { text: 'Email', value: 'email' },
      { text: 'Role', value: 'role' },
      { text: 'Status', value: 'isActive' },
      { text: 'Actions', value: 'actions', sortable: false }
    ]
  }),
  computed: {
    processedUsers() {
      if (!Array.isArray(this.users)) {
        return []
      }
      return this.users.map(user => ({
        id: user.id,
        username: user.username,
        email: user.email,
        role: user.role,
        isActive: user.isActive,
        createdAt: user.createdAt,
        updatedAt: user.updatedAt
      }))
    }
  },
  watch: {
    error(newError) {
      this.showError = !!newError
    }
  },
  methods: {
    handleOptionsUpdate(options) {
      this.options = options
      this.$emit('update:options', options)
    },
    handleSearch() {
      this.$emit('update:search', this.search)
    },
    confirmDelete(user) {
      this.userToDelete = user
      this.deleteDialog = true
    },
    async handleDelete() {
      try {
        await this.$store.dispatch('users/deleteUser', this.userToDelete.id)
        this.showSuccess = true
        this.successMessage = 'User deleted successfully'
        this.deleteDialog = false
        this.$emit('delete-success')
      } catch (error) {
        console.error('Failed to delete user:', error)
      }
    }
  }
}
</script> 