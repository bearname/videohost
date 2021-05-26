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
    this.searchVideos(this.search)
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
    searchVideo() {
      this.searchVideos(this.search)
    },
    searchVideos(searchString, page = '1', countVideoOnPage = '10') {

      const videoServerAddress = process.env.VUE_APP_VIDEO_SERVER_ADDRESS;
      let url = videoServerAddress + '/api/v1/videos/search?page=' + page + '&limit=' + countVideoOnPage + '&search=' + searchString;
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