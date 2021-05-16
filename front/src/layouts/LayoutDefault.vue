<template>
  <div id="app">
    <Navigation/>
    <p>{{ message }}</p>
    <input v-model="message">
    <h1>Первое приложение!</h1>
    <router-view></router-view>
  </div>
</template>

<script>
import Navigation from '../components/Navigation.vue'

export default {
  name: "LayoutDefault",
  components: {
    Navigation,
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