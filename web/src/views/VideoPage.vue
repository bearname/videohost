<template>
  <v-container class="grey lighten-5">
    <v-row no-gutters>
      <v-col
          cols="12"
      >
        <div v-if="video !== null" class="text-align-left">
          <div v-if="video.status === 3">
            <Player :videoId="videoId" :duration="video.duration" :thumbnail="video.thumbnail"
                    :availableQualities="video.quality" :chapters="video.chapters" :key="key"/>
          </div>
          <div v-else> status {{ videoStatus }}</div>
        </div>
        <h4 v-else>
          <v-text-field
              color="success"
              loading
              disabled
          ></v-text-field>
        </h4>
      </v-col>
    </v-row>
    <v-row no-gutters v-if="video !== null">
      <v-col
          cols="12"
          sm="8"
      >
        <h3>{{ video.name }}</h3>
        <p class="subtitle-1">Watch video {{ video.description }}</p>
        <p class="subtitle-2">Добавлено {{ video.uploaded }}</p>
        <p class="subtitle-2">{{ video.views }} views</p>
        <v-btn v-on:click="likeVideo(true)">like {{ video.countLikes }}</v-btn>
        <v-btn v-on:click="likeVideo(false)">dislike {{ video.countDisLikes }}</v-btn>
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
            <div v-if="error !== null"><span v-if="!error">Success</span><span v-else>Failed</span> updated video
              title
              and description
            </div>
          </div>
          <v-btn v-on:click="deleteItemPermanent(video.id)">delete</v-btn>
        </div>
      </v-col>
    </v-row>
    <!--    <v-row>-->
    <!--      <h4>Also see</h4>-->
    <!--      <Pagination :show-status="false" :user-page="false"/>-->
    <!--    </v-row>-->
  </v-container>


</template>

<script>
// import Pagination from '../components/Pagination.vue'
import Player from '../components/Player.vue'
import {mapActions, mapGetters} from "vuex";
import Cookie from "../util/cookie";
import VideoStatus from "../store/videoStore/videoStatus";
import videosUtil from "../store/videoStore/videoUtil"
import logError from "../util/logger";
import RESPONSE_CODES from "@/store/videoStore/responseCode";
import {VBtn} from "vuetify/lib";

export default {
  name: "StreamPage",
  components: {
    Player,
    VBtn,
    // Pagination,
  },
  data() {
    return {
      videoId: null,
      key: 0,
      video: null,
      currentUserId: null,
      error: null,
      code: null,
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
      findVideoById: "videoMod/getVideoById",
      deleteVideoPermanent: "videoMod/deleteVideoPermanent",
      likeVideoRequest: "videoMod/likeVideo",
    }),
    ...mapGetters({
      getVideoResult: "videoMod/getVideo",
    }),
    async setVideoId() {
      this.videoId = this.$route.params.videoId;
      await this.fetchVideo(this.$route.params.videoId);
    },
    async fetchVideo(videoId) {
      try {
        await this.findVideoById({videoId: videoId});
        this.video = this.getVideoResult();
        if (!('chapters' in this.video)) {
          this.video.chapters = [];
        }
        console.log(this.video);

        this.videoStatus = VideoStatus.intToStatus(this.video.status);
        this.video.uploaded = videosUtil.getElapsedString(this.video.uploaded);
      } catch (error) {
        logError(error);
      }
    },
    async updateTitleAndDescription() {
      const video = {
        videoId: this.video.id,
        name: this.video.name,
        description: this.video.description,
      };

      await this.$store.dispatch("videoMod/updateTitleAndDescription", video);
      console.log("update status");
      this.error = this.getStatus();
    },
    async deleteItemPermanent(videoId) {
      await this.deleteVideoPermanent({videoId: videoId});
      this.error = this.getStatus();
    },
    toggleEdit() {
      this.showEdit = !this.showEdit;
    },
    async likeVideo(isLike) {
      await this.likeVideoRequest({videoId: this.videoId, isLike: isLike, ownerId: this.video.ownerId});
      this.code = this.getCode();
      switch (this.code) {
        case RESPONSE_CODES.SuccessAddLike: {
          this.video.countLikes++;
          break;
        }
        case RESPONSE_CODES.SuccessAddDislike: {
          this.video.countDisLikes++;
          break;
        }
        case RESPONSE_CODES.SuccessDeleteLike: {
          this.video.countLikes--;
          break;
        }
        case RESPONSE_CODES.SuccessDeleteDisLike: {
          this.video.countDisLikes--;
          break;
        }
        default: {
          break
        }
      }
      this.error = this.getStatus();
    }
  },
}
</script>

<style scoped>
.text-align-left {
  text-align: left;
}
</style>