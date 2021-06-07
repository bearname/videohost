<template>
  <div>
    <v-btn v-on:click="closeShow">Toggle</v-btn>
    <div v-if="show">
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
  </div>
</template>

<script>
import VideoList from "./VideoList";
import {mapActions, mapGetters} from "vuex";

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
      show: false,
      isNeedDisplay: false
    }
  },
  async created() {
    if (this.search !== undefined && this.search !== null && this.search.length !== 0) {
      await this.searchVideos(this.search);
      console.log(this.search);
      document.addEventListener('keydown', (e) => {
        if (e.key === 'Slash') {
          let searchBoxElement = document.getElementById('searchBox');
          if (searchBoxElement !== null) {
            searchBoxElement.focus();
          }
        }
      });
    }
  },
  methods: {
    ...mapActions({
      searchVideosByQuery: "video/searchVideos"
    }),
    ...mapGetters({
      getVideos: "video/getVideos",
      getPageCount: "video/getPageCount"
    }),
    async searchVideo() {
      await this.searchVideos(this.search);
      this.show = true;
    },
    async searchVideos(searchString) {
      await this.searchVideosByQuery({searchString: searchString});
      this.videos = this.getVideos();
      console.log(this.videos);
    },
    closeShow() {
      this.show = !this.show;
    }
  }
}
</script>

<style scoped>

</style>