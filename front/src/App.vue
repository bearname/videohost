<template>
  <layout>
    <div>
      <Navigation/>
      <h1>Pretty player</h1>
      <VideoList/>
    </div>
  </layout>
</template>

<script>
import VideoList from './components/VideoList.vue'
import Navigation from './components/Navigation.vue'
import Layout from './layouts/LayoutDefault.vue'

export default {
  name: 'App',
  components: {
    Layout,
    Navigation: Navigation,
    VideoList: VideoList
  }
}
</script>

<style>
a { text-decoration: none; }

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
