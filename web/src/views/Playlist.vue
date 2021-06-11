<template>
  <div>
    {{listType}}
    <VideoList v-if="videos!== null" :show-status="true" :user-page="true" :videos="videos" :key="page"/>
  </div>
</template>

<script>
import {mapActions, mapGetters} from "vuex";
import VideoList from "@/components/VideoList";

export default {
  name: "Playlist",
  components: {
    VideoList,
  },
  async created() {
    this.listType = this.$route.query.PL;
    await this.fetchUserVideos(this.page, this.countVideo)
  },
  data() {
    return {
      videos: null,
      page: 0,
      countVideo: 30,
      listType: null
    }
  },
  methods: {
    ...mapActions({
      findUserLikedVideos: "user/getUserLikedVideos",
    }),
    ...mapGetters({
      getVideoResult: "user/getUserVideos",
    }),
    async fetchUserVideos(page, countVideo) {
      await this.findUserLikedVideos({page: page, countVideoOnPage: countVideo});
      const result = this.getVideoResult();
      this.videos = result.videos;
    }
  },
}
</script>

<style scoped>

</style>