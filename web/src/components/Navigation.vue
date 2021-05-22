<template>
  <v-card
      color="grey lighten-4"
      flat
      tile
      class="fixed"
  >
    <v-toolbar dense>
      <v-toolbar-title class="mr-16">
        <router-link :to="{ name: 'home'}">Каталог</router-link>
      </v-toolbar-title>

      <v-toolbar-title>
        <router-link :to="{ name: 'uploadVideo'}">Загрузить видео</router-link>
      </v-toolbar-title>
      <v-toolbar-title>
        <router-link :to="{ name: 'user'}">User</router-link>
      </v-toolbar-title>
      <v-toolbar-title>
        <div v-if="isLoggedUser === false">
          <router-link :to="{ name: 'login'}">login</router-link>
        </div>
        <div v-else>
          <v-btn v-on:click="logoutUser()">logout</v-btn>
        </div>
      </v-toolbar-title>
    </v-toolbar>
  </v-card>
</template>

<script>
import {mapActions, mapGetters} from "vuex";

export default {
  name: "Navigation",
  data() {
    return {
      userid: {type: String, default: ""}
    }
  },
  methods: {
    ...mapActions({
      logout: "auth/logout",
    }),
    ...mapGetters({
      isLoggedIn: "auth/isLoggedIn",
      getUser: "user/getUser",
      getCurrentUser: "auth/getCurrentUser"
    }),
    logoutUser() {
      this.logout()
      this.isLoggedUser()
    },
  },
  computed: {
    user() {
      return this.getUser(this.userid);
    },
    currentUsername() {
      return this.getCurrentUser().username;
    },
    isLoggedUser() {
      return this.isLoggedIn()
    }
  },
}
</script>

<style scoped>
.fixed {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 1;
}
</style>