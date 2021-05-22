<template>
  <div>
    <v-btn
        class="hidden-xs-only"
        v-on:click="toggleDisplay"
    >
      <span v-if="isNeedDisplay">close </span>
      <span v-else>open </span>
      <span> search </span>
    </v-btn>
    <div v-if="isNeedDisplay">
      <div class="search-wrapper">
        <input type="text" v-model="search" placeholder="Search title.."/>
      </div>
      <v-btn v-on:click="onSearchSubmit">Search</v-btn>
      <div v-if="videos !== null && videos !== []">
        <VideoList :videos="videos" :show-status="false" :user-page="false" :key="videos"/>
      </div>
      <div v-else>
        Not Found
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import VideoList from "./VideoList";

export default {
  name: "Search",
  components: {VideoList},
  data() {
    return {
      search: null,
      videos: null,
      error: false,
      page: 1,
      isNeedDisplay: false
    }
  },
  methods: {
    toggleDisplay() {
      this.isNeedDisplay = !this.isNeedDisplay
    },
    onSearchSubmit() {
      this.searchVideos(this.search)
    },
    searchVideos(searchString, page = '1', countVideoOnPage = '10') {
      const videoServerAddress = process.env.VUE_APP_VIDEO_SERVER_ADDRESS;
      let url = videoServerAddress + '/api/v1/videos/search?page=' + page + '&limit=' + countVideoOnPage + '&search=' + searchString;
      console.log(url)
      axios.get(url)
          .then(response => {
            console.log(response.data)
            if (Object.keys(response.data).includes("pageCount")) {
              this.countPage = response.data.pageCount
            }
            if (Object.keys(response.data).includes("videos")) {
              this.videos = response.data.videos
              this.videos.forEach(function (part, index) {
                this[index].thumbnail = videoServerAddress + this[index].thumbnail
              }, this.videos)
              console.log("videos", this.videos)
            }
          })
          .catch(function (error) {
            if (error.response) {
              console.log(error.response.data);
              console.log(error.response.status);
              console.log(error.response.headers);
            } else if (error.request) {
              console.log(error.request);
            } else {
              console.log('Error', error.message);
            }
          });
    }
  }
}
</script>

<style scoped>

</style>