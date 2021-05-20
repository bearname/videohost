<template>
  <div
      class="mx-auto"
      max-width="1000"
  >
    <div v-if="videos !== null">
      <VideoList :videos="videos" :key="pageNumber"/>
      <v-btn class="v-btn"
             :disabled="pageNumber === 0"
             @click="previousPage">
        Previous
      </v-btn>
      <v-btn>{{ pageNumber + 1 }}</v-btn>
      <v-btn class="v-btn"
             :disabled="(pageNumber >= countPage && countPage < 2) || (pageNumber >= countPage - 1 && countPage >= 2) "
             @click="nextPage">
        Next
      </v-btn>
    </div>
    <p v-else>Not found any video</p>
  </div>

</template>

<script>
import VideoList from "./VideoList";
import {mapActions, mapGetters} from "vuex";

export default {
  name: "Pagination",
  components: {
    VideoList
  },

  data() {
    return {
      error: false,
      pageNumber: 1,
      videos: null,
      countVideoOnPage: 12,
      countPage: 0,
      url: {
        type: String,
        required: false,
        default: "list"
      }
    }
  },
  methods: {
    ...mapActions({
      getVideoOnPage: "video/getVideoOnPage"
    }),
    ...mapGetters({
      getVideos: "video/getVideos"
    }),
    async previousPage() {
      if (this.pageNumber > 0) {
        this.pageNumber--;
        await this.fetchVideosByPage(this.pageNumber, this.countVideoOnPage)
      }
    },

    async nextPage() {
      if (this.pageNumber < this.countPage) {
        this.pageNumber++;
        await this.fetchVideosByPage(this.pageNumber, this.countVideoOnPage)
      }
    },

    async fetchVideosByPage(page = '1', countVideoOnPage = '10') {
      await this.getVideoOnPage(page, countVideoOnPage)
          .then(() => {
            this.videos = this.getVideos()
          })
    }
  },
  mounted() {
    console.log("axios")
    this.fetchVideosByPage()
  }
}
</script>

<style scoped>

</style>