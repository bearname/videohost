<template>
  <div>
    <Alert/>
    <Navigation/>
    <Search style="margin-top: 60px"/>
    <router-view></router-view>
  </div>
</template>

<script>
import Navigation from '../components/Navigation.vue'
import Search from '../components/Search.vue'
import {VIcon} from "vuetify/lib";
import Alert from "@/components/Alert";

export default {
  name: "LayoutDefault",
  components: {
    Alert,
    Navigation,
    Search,
    VIcon,
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