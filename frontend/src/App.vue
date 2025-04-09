<template>
  <v-app>
    <v-navigation-drawer
      v-model="drawer"
      :mini-variant="miniVariant"
      :clipped="clipped"
      fixed
      app
    >
      <v-list>
        <v-list-item
          v-for="(item, i) in menuItems"
          :key="i"
          :to="item.action ? undefined : item.to"
          router
          exact
          @click="handleMenuItemClick(item)"
        >
          <v-list-item-action>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title v-text="item.title" />
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar
      :clipped-left="clipped"
      fixed
      app
    >
      <v-app-bar-nav-icon @click.stop="drawer = !drawer" />
      <v-toolbar-title v-text="title" />
      <v-spacer />
      <v-btn
        v-if="isAuthenticated"
        icon
        @click="logout"
      >
        <v-icon>mdi-logout</v-icon>
      </v-btn>
    </v-app-bar>

    <v-main>
      <v-container fluid>
        <router-view />
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
export default {
  name: 'App',
  data: () => ({
    clipped: false,
    drawer: false,
    fixed: false,
    miniVariant: false,
    right: true,
    rightDrawer: false,
    title: 'Bell Scheduler',
    menuItems: [
      {
        icon: 'mdi-view-dashboard',
        title: 'Dashboard',
        to: '/'
      },
      {
        icon: 'mdi-bell',
        title: 'Schedules',
        to: '/schedules'
      },
      {
        icon: 'mdi-account-group',
        title: 'Users',
        to: '/users'
      },
      {
        icon: 'mdi-clipboard-text-clock',
        title: 'Logs',
        to: '/logs'
      },
      {
        icon: 'mdi-cog',
        title: 'Settings',
        to: '/settings'
      },
      {
        icon: 'mdi-logout',
        title: 'Logout',
        to: '#',
        action: 'logout'
      }
    ]
  }),
  computed: {
    isAuthenticated() {
      return this.$store.getters['auth/isAuthenticated']
    }
  },
  methods: {
    async logout() {
      await this.$store.dispatch('auth/logout')
      this.$router.push('/login')
    },
    handleMenuItemClick(item) {
      if (item.action === 'logout') {
        this.logout()
      }
      this.drawer = false
    }
  }
}
</script>