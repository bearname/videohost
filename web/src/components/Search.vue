<template>
  <div>
    {{ search }}
    <div class="search-wrapper">
      <input id="searchBox" type="text" v-model="search" placeholder="Search title.."/>
    </div>
    <v-btn v-on:click="searchVideo">Search</v-btn>
    <div v-if="videos !== null && videos !== []">
      <VideoList :videos="videos" :show-status="false" :user-page="false" :key="videos"/>
    </div>
    <div v-else>
      Not Found
    </div>
  </div>
</template>

<script>
import VideoList from "./VideoList";
import axios from "axios";
import logError from "../util/logger";
import {updateThumbnail} from "../store/videoStore/video";

export default {
  name: "Search",
  components: {VideoList},
  props: [
    'searchQuery'
  ],
  data() {
    return {
      search: this.searchQuery,
      videos: null,
      error: false,
      page: 1,
      isNeedDisplay: false
    }
  },
  async created() {
    await this.searchVideos(this.search)
    document.addEventListener('keydown', (e) => {
      if (e.key === 'Slash') {
        let searchBoxElement = document.getElementById('searchBox');
        if (searchBoxElement !== null) {
          searchBoxElement.focus();
        }
      }
    })
  },
  methods: {
    // ...mapActions({
    //   searchVideosByQuery: "video/searchVideos"
    // }),
    // ...mapGetters({
    //   getVideos: "video/getVideos",
    //   getPageCount: "video/getPageCount"
    // }),
    // async searchVideo() {
    //   await this.searchVideosByQuery({searchString: this.search});
    //   this.videos = this.getVideos();
    // },
    async searchVideo() {
      await this.searchVideos(this.search)
    },
    async searchVideos(searchString, page = '1', countVideoOnPage = '10') {
      try {
        const videoServerAddress = process.env.VUE_APP_VIDEO_API;
        const url = videoServerAddress + '/api/v1/videos/search?page=' + page + '&limit=' + countVideoOnPage + '&search=' + searchString;
        const response = await axios.get(url);
        const data = response.data;

        console.log(data)

        if (Object.keys(data).includes("pageCount")) {
          this.countPage = data.pageCount
        }
        if (Object.keys(data).includes("videos")) {
          this.videos = data.videos
          this.videos.forEach((part, index) => {
            updateThumbnail(part, index)
          }, this.videos);
        }
      } catch (error) {
        logError(error)
      }
    }
  }
}
</script>

<style scoped>

</style>