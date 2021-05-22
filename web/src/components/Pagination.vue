<template>
  <div
      class="mx-auto"
      max-width="1000"
  >
    <div v-if="videos !== null">
      <VideoList :videos="videos" :current-user-id="currentUserId" :show-status="showStatus" :user-page="userPage"  :key="pageNumber"/>
      <v-btn class="v-btn"
             :disabled="pageNumber === 0"
             @click="previousPage">
        Previous
      </v-btn>
      <v-btn>{{ pageNumber + 1 }} of {{countPage + 1}}</v-btn>
      <v-btn class="v-btn"
             :disabled="pageNumber >= countPage"
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
  props: [
    'showStatus',
    'userPage',
  ],
  created() {
    const getter = this.$store.getters["auth/getCurrentUser"];
    this.currentUserId = getter.id;
  },
  data() {
    return {
      error: false,
      pageNumber: 1,
      videos: null,
      countVideoOnPage: 16,
      countPage: null,
      currentUserId: null,
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
      getVideos: "video/getVideos",
      getPageCount: "video/getPageCount"
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
    async fetchVideosByPage(page , countVideoOnPage) {
      await this.getVideoOnPage(page, countVideoOnPage)
          .then(() => {
            this.videos = this.getVideos()
            this.countPage = this.getPageCount()
          })
    }
  },
  mounted() {
    this.fetchVideosByPage(this.pageNumber, this.countVideoOnPage)
  }
}
</script>

<style scoped>

</style>