<template>
  <div>
    <Navigation/>

    <Search style="margin-top: 60px"/>
    <router-view></router-view>
    <v-row justify="space-around">
      <v-icon
          large
          color="green darken-2"
      >
        mdi-domain
      </v-icon>

      <v-icon
          large
          color="blue darken-2"
      >
        mdi-message-text
      </v-icon>

      <v-icon
          large
          color="purple darken-2"
      >
        mdi-dialpad
      </v-icon>

      <v-icon
          large
          color="teal darken-2"
      >
        mdi-email
      </v-icon>

      <v-icon
          large
          color="blue-grey darken-2"
      >
        mdi-call-split
      </v-icon>

      <v-icon
          large
          color="orange darken-2"
      >
        mdi-arrow-up-bold-box-outline
      </v-icon>
      <v-icon>arrow_downward</v-icon>
    </v-row>
  </div>
</template>

<script>
import Navigation from '../components/Navigation.vue'
import Search from '../components/Search.vue'

export default {
  name: "LayoutDefault",
  components: {
    Navigation,
    Search,
  },
  data() {
    return {
      drawer: null,
      message: 'Привет, Vue!'
    }
  }
}
</script>

<style scoped>
</style>