<template>
  <v-container>
    <v-row>
      <v-col>
        <h1>Manage Users</h1>
      </v-col>
      <v-col class="text-right">
        <v-btn
          color="primary"
          @click="handleAdd"
          :loading="loading"
        >
          Add User
        </v-btn>
      </v-col>
    </v-row>

    <v-row>
      <v-col>
        <user-list
          :users="users"
          :loading="loading"
          :error="error"
          :total-users="pagination.totalItems"
          @update:options="handleOptionsUpdate"
          @update:search="handleSearchUpdate"
          @edit="handleEdit"
          @delete="handleDelete"
        />
      </v-col>
    </v-row>

    <v-dialog v-model="showForm" max-width="600px">
      <user-form
        :user="selectedUser"
        :is-edit="!!selectedUser"
        :loading="loading"
        :error="error"
        @success="handleFormSuccess"
        @cancel="showForm = false"
      />
    </v-dialog>
  </v-container>
</template>

<script>
import UserList from '@/components/users/UserList.vue'
import UserForm from '@/components/users/UserForm.vue'

export default {
  name: 'Users',
  components: {
    UserList,
    UserForm
  },
  data() {
    return {
      showForm: false,
      selectedUser: null,
      isInitialLoad: true
    }
  },
  computed: {
    loading() {
      return this.$store.getters['users/loading']
    },
    error() {
      return this.$store.getters['users/error']
    },
    users() {
      return this.$store.getters['users/users'] || []
    },
    pagination() {
      return this.$store.getters['users/pagination']
    },
    sort() {
      return this.$store.getters['users/sort']
    },
    filter() {
      return this.$store.getters['users/filter']
    }
  },
  created() {
    this.fetchUsers()
  },
  methods: {
    async fetchUsers() {
      try {
        await this.$store.dispatch('users/fetchUsers')
      } catch (error) {
        // Handle error silently or show error message to user
      }
    },
    handleAdd() {
      this.selectedUser = null
      this.showForm = true
    },
    handleEdit(user) {
      this.selectedUser = user
      this.showForm = true
    },
    async handleDelete(user) {
      try {
        await this.$store.dispatch('users/deleteUser', user.id)
        // Removed duplicate fetchUsers call as the store action already handles it
      } catch (error) {
        // Handle error silently or show error message to user
      }
    },
    handleFormSuccess() {
      this.showForm = false
      // Removed duplicate fetchUsers call as the store actions already handle it
    },
    handleOptionsUpdate(options) {
      const { page, itemsPerPage, sortBy, sortDesc } = options
      // Update store state without triggering additional API calls
      this.$store.commit('users/SET_PAGINATION', {
        page,
        itemsPerPage
      })
      this.$store.commit('users/SET_SORT', {
        by: sortBy,
        desc: sortDesc
      })
      // Only fetch users if this is not the initial component load
      if (!this.isInitialLoad) {
        this.fetchUsers()
      }
      // Reset the flag after initial load
      this.isInitialLoad = false
    },
    handleSearchUpdate(search) {
      this.$store.commit('users/SET_FILTER', { search })
      this.fetchUsers()
    }
  }
}
</script>