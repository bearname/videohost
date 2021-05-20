<template>
  <layout>
    <div>
      <Navigation/>
      <Search/>
      <h1>Pretty player</h1>
      <Pagination/>
    </div>
  </layout>
</template>

<script>
import Pagination from './components/Pagination.vue'
import Navigation from './components/Navigation.vue'
import Layout from './layouts/LayoutDefault.vue'
import Search from "./components/Search";

export default {
  name: 'App',
  components: {
    Search,
    Layout,
    Navigation,
    Pagination
  }
}
</script>

<style>
a {
  text-decoration: none;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
