<template>
  <v-container fluid>
    <v-row dense>
      <v-col
          v-for="video in videoList" :key="video.id"
          :cols="3"
      >
        <VideoItem :video="video" :current-user-id="currentUserId" :show-status="showStatus" :user-page="isUserPage" :key="video.id"/>
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
    'videos',
    'showStatus',
    'userPage',
    'currentUserId',
  ],
  data() {
    return {
      videoList: this.videos,
      currentPage: this.pageNumber,
      isUserPage: this.userPage,
    }
  },

  mounted() {
    this.videoList = this.videos
  }
}
</script>

<style scoped>

</style>