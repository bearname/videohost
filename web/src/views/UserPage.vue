<template>
  <div>
    <h2>{{ currentUsername }}</h2>
    <h4>my video</h4>
    <h5>Count you video {{ countAllVideos }}</h5>
    <Pagination :show-status="false" :user-page="true"/>
<!--    <VideoList v-if="videos!== null" :show-status="true" :user-page="true" :videos="videos" :key="page"/>-->
  </div>
</template>

<script>
import {mapActions, mapGetters} from "vuex";
// import VideoList from "../components/VideoList";
import Pagination from "@/components/Pagination";

export default {
  name: "User",
  components: {
    Pagination,
    // VideoList
  },
  data() {
    return {
      videos: null,
      countVideoOnPage: 10,
      page: 0,
      countAllVideos: 0
    }
  },
  created() {
    this.getAsyncVideos()
  },
  computed: {
    user() {
      return this.getUser();
    },
    currentUsername() {
      return this.getCurrentUser().username;
    },
    isLoggedUser() {
      return this.isLoggedIn()
    }
  },
  methods: {
    ...mapActions({
      fetchUserVideos: "video/fetchUserVideos"
    }),
    ...mapGetters({
      getUserVideos: "video/getUserVideos",
      getUser: "user/getUser",
      isLoggedIn: "auth/isLoggedIn",
      getCurrentUser: "auth/getCurrentUser"
    }),
    getAsyncVideos() {
      this.fetchUserVideos({page: this.page, countVideoOnPage: this.countVideoOnPage})
          .then(() => {
            const result = this.getUserVideos()
            this.videos = result.videos
            this.countAllVideos = result.countAllVideos
          })
    }
  },
}
</script>

<style scoped>

</style>