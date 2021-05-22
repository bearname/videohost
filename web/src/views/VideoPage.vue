<template>
  <div>
    Video page
    <div class="text-align-left" v-if="video !== null">
      <div v-if="video.status === '3'">
        <Player :videoId="videoId" :key="key"/>
      </div>
      <div v-else> status {{video.status}}{{ videoStatus }}</div>
      <h3>{{ video.name }}</h3>
      <p class="subtitle-1">Watch video {{ video.description }}</p>
      <p class="subtitle-2">Добавлено {{ video.uploaded }}</p>
      <p class="subtitle-2">{{ video.views }} views</p>
      <div v-if="currentUserId !== null && video.ownerId === currentUserId">
        <v-btn v-on:click="deleteItemPermanent(video.id)" :data-id="video.id">delete</v-btn>
      </div>
    </div>
    <div v-else>
      Video not exist
    </div>

    <v-spacer></v-spacer>
    <div>
      <h4>Also see</h4>
      <Pagination :show-status="false" :user-page="false"/>
    </div>
  </div>
</template>

<script>
import Pagination from '../components/Pagination.vue'
import Player from '../components/Player.vue'
import axios from 'axios'
import {mapActions, mapGetters} from "vuex";
import Cookie from "../util/cookie";
import VideoStatus from "../store/videoStore/videoStatus";

export default {
  name: "StreamPage",
  components: {
    Player,
    Pagination
  },
  data() {
    return {
      videoId: null,
      key: 0,
      video: null,
      currentUserId: null,
      error: false,
      videoStatus: null
    }
  },
  created() {
    this.setVideoId()
    this.key = Date.now()
    // this.currentUserId =
        Cookie.getCookie("userId");
    console.log(this.videoId)
  },
  watch: {
    '$route'() {
      console.log("router")
      this.setVideoId()
      this.key = Date.now()

      console.log(this.videoId)
    }
  },
  methods: {
    ...mapActions({
      deleteVideoPermanent: "video/deleteVideoPermanent"
    }),
    ...mapGetters({
      getStatus: "video/getStatus"
    }),
    async setVideoId() {
      this.videoId = this.$route.params.videoId
      this.fetchVideo(this.videoId)
    },
    fetchVideo(videoId) {
      let url = process.env.VUE_APP_VIDEO_SERVER_ADDRESS + '/api/v1/videos/' + videoId;
      console.log(url)
      axios.get(url)
          .then(response => {
            console.log(response.data)
            this.video = response.data
            this.videoStatus = VideoStatus.intToStatus(this.video.status)
            this.video.uploaded = this.getElapsedString(this.video.uploaded)
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
            this.error = true
          });
    },
    deleteItemPermanent(videoId) {
      const promise = this.deleteVideoPermanent({videoId: videoId});
      promise.then(() => {
        this.status = this.getStatus();
      })
    },
    getElapsedString(uploadedDate) {
      let elapsed = (Date.now() - Date.parse(uploadedDate)) / 1000;

      elapsed = Math.round(elapsed);

      if (elapsed / 60 < 1) {
        return elapsed + " секунд назад"
      }
      elapsed = Math.round(elapsed / 60);
      if (elapsed / 60 < 1) {
        return elapsed + " минут назад"
      }
      elapsed = Math.round(elapsed / 24);
      if (elapsed / 24 < 1) {
        return elapsed + " часов назад"
      }
      let weeks = Math.round(elapsed / 7);
      if (weeks / 7 < 1) {
        return weeks + " недель назад"
      }
      let months = Math.round(elapsed / 30);
      if (months / 30 < 1) {
        return months + " месяцев назад"
      }
      let years = Math.round(elapsed / 365);
      if (years / 365 < 1) {
        return years + " лет назад"
      }
    },

  }
}
</script>

<style scoped>
.text-align-left {
  text-align: left;
}
</style>