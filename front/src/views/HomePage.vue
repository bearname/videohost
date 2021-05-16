<template>
  <div>
    <h2>Home</h2>
    <VideoList/>
  </div>
</template>

<script>
import VideoList from '../components/VideoList.vue'

export default {
  name: "Home",
  components: {
    VideoList: VideoList
  }
}
</script>

<style scoped>

</style>