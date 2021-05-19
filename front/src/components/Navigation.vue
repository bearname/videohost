<template>
  <v-card
      color="grey lighten-4"
      flat
      tile
      class="fixed"
  >
    <v-toolbar dense>
<!--      <v-app-bar-nav-icon></v-app-bar-nav-icon>-->
      <v-toolbar-title class="mr-16">
        <router-link :to="{ name: 'home'}">Каталог</router-link>
      </v-toolbar-title>

      <v-toolbar-title>
        <router-link :to="{ name: 'uploadVideo'}">Загрузить видео</router-link>
      </v-toolbar-title>

      <v-toolbar-title >
        <div v-if="$store.getters['auth/isLoggedIn'] === false">
          <router-link :to="{ name: 'login'}">login</router-link>
        </div>
        <div v-else>
          <v-btn v-on:click="logoutUser">logout</v-btn>
        </div>
      </v-toolbar-title>

      <v-spacer></v-spacer>

<!--      <v-btn icon>-->
<!--        <v-icon>mdi-magnify</v-icon>-->
<!--      </v-btn>-->

<!--      <v-btn icon>-->
<!--        <v-icon>mdi-heart</v-icon>-->
<!--      </v-btn>-->

<!--      <v-btn icon>-->
<!--        <v-icon>mdi-dots-vertical</v-icon>-->
<!--      </v-btn>-->
    </v-toolbar>
  </v-card>
</template>

<script>
import {mapActions} from "vuex";

export default {
  name: "Navigation",
  methods: {
    ...mapActions({
      uploadVideo: "auth/logout"
    }),

    logoutUser() {
      this.logout()
    },
    // isLoggedUser() {
    //   const state1 = this.$store.state;
    //   // if (!state1.user) {
    //   //   return false;
    //   // }
    //   //
    //   // return state1.user.loggedIn;
    //   return this.isLoggedIn(state1)
    // }
  }
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