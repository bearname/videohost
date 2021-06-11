<template>
  <div
      class="mx-auto"
      max-width="1000"
  >
    <div v-if="videos !== null">
      <VideoList :videos="videos" :current-user-id="currentUserId" :show-status="showStatus" :user-page="userPage"
                 :key="videos[0].id"/>
      <v-btn class="v-btn"
             :disabled="pageNumber === 0"
             @click="previousPage">
        Previous
      </v-btn>
      <v-btn>{{ pageNumber + 1 }} of {{ countPage + 1 }}</v-btn>
      <v-btn class="v-btn"
             :disabled="pageNumber >= countPage - 1"
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
import {VBtn} from 'vuetify/lib'

export default {
  name: "Pagination",
  components: {
    VideoList,
    VBtn
  },
  props: [
    'showStatus',
    'userPage',
  ],
  created() {
    const getter = this.$store.getters["authMod/getCurrentUser"];
    this.currentUserId = getter.id;
  },
  data() {
    return {
      error: false,
      pageNumber: 0,
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
      getVideoOnPage: "videoMod/getVideoOnPage"
    }),
    ...mapGetters({
      getVideos: "videoMod/getVideos",
      getPageCount: "videoMod/getPageCount"
    }),
    async previousPage() {
      if (this.pageNumber > 0) {
        --this.pageNumber;
        await this.fetchVideosByPage(this.pageNumber, this.countVideoOnPage)
      }
    },
    async nextPage() {
      if (this.pageNumber < this.countPage) {
        ++this.pageNumber;
        await this.fetchVideosByPage(this.pageNumber, this.countVideoOnPage)
      }
    },
    async fetchVideosByPage(page, countVideoOnPage) {
      await this.getVideoOnPage(page, countVideoOnPage)
      this.videos = this.getVideos()
      this.countPage = this.getPageCount()
    }
  },
  async mounted() {
    await this.fetchVideosByPage(this.pageNumber, this.countVideoOnPage)
  }
}
</script>

<style scoped>

</style>