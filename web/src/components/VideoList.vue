<template>
  <v-container fluid>
    <v-row dense>
      <v-col
          v-for="video in videoList" :key="video.id"
          :cols="3"
      >
        <VideoItem :video="video" :key="video.id"/>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import VideoItem from "./VideoItem";

export default {
  name: "VideoList",
  components: {VideoItem},
  props: [
    'pageNumber',
    'videos'
  ],
  data() {
    return {
      videoList: this.videos,
      currentPage: this.pageNumber,
    }
  },

  mounted() {
    this.videoList = this.videos
  }
}
</script>

<style scoped>

</style>