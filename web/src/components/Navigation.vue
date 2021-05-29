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
<!--      <SearchRow/>-->
      <v-toolbar-title class="float-right">
        <router-link :to="{ name: 'user'}">User {{ currentUsername }}</router-link>
      </v-toolbar-title>
      <v-toolbar-title class="float-right">
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
// import SearchRow from "./SearchRow";

export default {
  name: "Navigation",
  // components: {SearchRow},
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
      isLogged: "auth/isLoggedIn",
      getUser: "user/getCurrentUser",
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
    currentUsername: function () {
      const currentUser = this.getCurrentUser();
      return currentUser.username;
    },
    isLoggedUser() {
      return this.isLogged()
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