<template>
  <div id="app">

    <Navigation/>
    <Search/>
    <router-view></router-view>
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
      message: 'Привет, Vue!'
    }
  }
}
</script>

<style scoped>

</style>