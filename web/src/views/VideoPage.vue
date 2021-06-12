<template>
  <v-container class="grey lighten-5">
    <v-row no-gutters>
      <v-col
          cols="12"
      >
        asdfadsf
        <div v-if="video !== null" class="text-align-left">
          {{ video.status }}
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
        <PlaylistAddPopup :video-id="video.id"/>
        <v-btn v-if="!isCurrentUserOwner" v-on:click="followToAuthor">follow to author</v-btn>
        <h3>{{ video.name }}</h3>
        <p class="subtitle-1">Watch video {{ video.description }}</p>
        <p class="subtitle-2">Добавлено {{ video.uploaded }}</p>
        <p class="subtitle-2">{{ video.views }} views</p>
        <v-btn id="countLikes" v-on:click="likeVideo(true)">Likes {{ countLikes }}</v-btn>
        <v-btn id="countDislikes" v-on:click="likeVideo(false)">Dislikes {{ countDisLikes }}</v-btn>
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
import PlaylistAddPopup from "@/components/PlaylistAddPopup";
import {publishEvent} from "@/events/event-bus";

export default {
  name: "StreamPage",
  components: {
    PlaylistAddPopup,
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
      countLikes: 0,
      countDisLikes: 0,
      likeElement: null,
      disLikeElement: null,
    }
  },
  created() {
    this.getVideo()
    this.key = Date.now()
    this.currentUserId = Cookie.getCookie("userId");
  },
  watch: {
    '$route'() {
      this.getVideo()
      this.key = Date.now()
    },
    likes: function () {
      return "likes " + this.countLikes
    },
    dislikes: function () {
      return "dislikes " + this.countDisLikes
    },
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
      follow: "userMod/subscribe",
    }),
    ...mapGetters({
      isSuccessResult: "userMod/isSuccess",
      getVideoResult: "videoMod/getVideo",
      getResponseCode: "videoMod/getCode",
    }),
    async followToAuthor() {
      await this.follow({followingToUserId: this.video.ownerId});
      const isOk = this.isSuccessResult();

      publishEvent((isOk) ? "true" : "false", (isOk) ? "Success" : "Failed" + "subscribe");
    },
    async getVideo() {
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
        this.countLikes = this.video.countLikes;
        this.countDisLikes = this.video.countDisLikes;
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
      const data = {videoId: this.videoId, isLike: isLike, ownerId: this.video.ownerId};
      console.log(data);
      await this.likeVideoRequest(data);
      this.code = this.getResponseCode();
      console.log('this.code');
      console.log(this.code);
      const code = parseInt(this.code);
      switch (code) {
        case RESPONSE_CODES.ErrInternalServer: {
          publishEvent('false', "Internal server error");
          break;
        }
        case RESPONSE_CODES.SuccessAddLike: {
          this.countLikes++;
          publishEvent('success', "Success Add Like");
          document.getElementById("countLikes").innerText = "likes " + this.countLikes;
          break;
        }
        case RESPONSE_CODES.SuccessAddDislike: {
          publishEvent('success', "Success Add Dislike");
          this.countDisLikes++;
          document.getElementById("countDisLikes").innerText = "dislikes " + this.countLikes;
          break;
        }
        case RESPONSE_CODES.SuccessDeleteLike: {
          publishEvent('success', "Success Remove Like");
          this.countLikes--;
          document.getElementById("countLikes").innerText = "likes " + this.countLikes;
          break;
        }
        case RESPONSE_CODES.SuccessDeleteDisLike: {
          publishEvent('success', "Success Remove Dislike");
          this.countDisLikes--;
          document.getElementById("countDisLikes").innerText = "dislikes " + this.countLikes;
          break;
        }
        case RESPONSE_CODES.ErrAlreadyLike: {
          publishEvent('false', "Already Liked");
          break;
        }
        case RESPONSE_CODES.ErrAlreadyDisLike: {
          publishEvent('false', "Already DisLiked");
          break;
        }
        case RESPONSE_CODES.ErrFailedAddLike: {
          publishEvent('false', "Already Failed Add Like");
          break;
        }
        case RESPONSE_CODES.ErrFailedDeleteLike: {
          publishEvent('false', "Already Failed Remove Like");
          break;
        }
        default: {
          break
        }
      }
    },
  },
}

</script>

<style scoped>
.text-align-left {
  text-align: left;
}
</style>