<template>
  <div>
    Video page
    <div class="text-align-left" v-if="video !== null">
      <div v-if="video.status === '3'">
        <Player :videoId="videoId" :duration="video.duration" :thumbnail="video.thumbnail" :availableQualities="video.quality" :key="key"/>
      </div>
      <div v-else> status {{ video.status }}{{ videoStatus }}</div>
      <h3>{{ video.name }}</h3>
      <p class="subtitle-1">Watch video {{ video.description }}</p>
      <p class="subtitle-2">Добавлено {{ video.uploaded }}</p>
      <p class="subtitle-2">{{ video.views }} views</p>
      <div v-if="isCurrentUserOwner">
        <v-btn v-on:click="toggleEdit" :data-id="video.id">edit</v-btn>
        <div v-if="showEdit">
          <div>
            <label for="name">Name: <input id="name" type="text" v-model="video.name"></label>
          </div>
          <div>
            <label for="description">Description: <input id="description" type="text"
                                                         v-model="video.description"></label>
          </div>
          <v-btn type="submit" v-on:click="updateTitleAndDescription">update</v-btn>
          <div v-if="error !== null"><span v-if="!error">Success</span><span v-else>Failed</span> updated video title
            and description
          </div>
        </div>
        <v-btn v-on:click="deleteItemPermanent(video.id)">delete</v-btn>
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
import getElapsedString from "../store/videoStore/video"

export default {
  name: "StreamPage",
  components: {
    Player,
    Pagination,
  },
  data() {
    return {
      videoId: null,
      key: 0,
      video: null,
      currentUserId: null,
      error: null,
      videoStatus: null,
      userVideos: null,
      showEdit: false,
    }
  },
  created() {
    this.setVideoId()
    this.key = Date.now()
    this.currentUserId = Cookie.getCookie("userId");
  },
  watch: {
    '$route'() {
      this.setVideoId()
      this.key = Date.now()
    }
  },
  computed: {
    isCurrentUserOwner() {
      return this.currentUserId !== null && this.video.ownerId === this.currentUserId
    },
  },
  methods: {
    ...mapActions({
      deleteVideoPermanent: "video/deleteVideoPermanent",
    }),
    ...mapGetters({
      getStatus: "video/getStatus",
    }),
    setVideoId() {
      this.videoId = this.$route.params.videoId
      this.fetchVideo(this.videoId)
    },
    fetchVideo(videoId) {
      let url = process.env.VUE_APP_VIDEO_SERVER_ADDRESS + '/api/v1/videos/' + videoId;
      axios.get(url)
          .then(response => {
            console.log(response.data)
            this.video = response.data
            this.videoStatus = VideoStatus.intToStatus(this.video.status)
            this.video.uploaded = getElapsedString(this.video.uploaded)
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
    async updateTitleAndDescription() {
      const video = {
        videoId: this.video.id,
        name: this.video.name,
        description: this.video.description,
      };
      await this.$store.dispatch("video/updateTitleAndDescription", video);
      console.log("update status")
      this.error = this.getStatus();
    },
    async deleteItemPermanent(videoId) {
      const promise = this.deleteVideoPermanent({videoId: videoId});
      const newVar = await promise;
      console.log(newVar)
      this.error = this.getStatus();
    },
    toggleEdit() {
      this.showEdit = !this.showEdit
    },
  },
}
</script>

<style scoped>
.text-align-left {
  text-align: left;
}
</style>