<template>
  <div>
    <h2>{{ currentUsername }}</h2>
    <h4>my video</h4>
    <h5>Count you video {{ countAllVideos }}</h5>
    <VideoList v-if="videos!== null" :show-status="true" :user-page="true" :videos="videos" :key="page"/>
  </div>
</template>

<script>
import {mapActions, mapGetters} from "vuex";
import VideoList from "../components/VideoList";

export default {
  name: "User",
  components: {
    VideoList,
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
      fetchUserVideos: "video/getUserVideos"
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