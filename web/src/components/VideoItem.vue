<template>
  <div style="overflow-y: auto; height:340px">
    <v-img
        :src="videoItem.thumbnail"
        class="white--text align-end"
        gradient="to bottom, rgba(0,0,0,.1), rgba(0,0,0,.5)"
        height="200px"
    ></v-img>
    <v-card-title>
      <router-link :to="{ name: 'videoStream', params: { videoId: videoItem.id }}" class="subtitle-2"> <span
          class="subtitle-2">{{
          videoItem.name
        }} </span>
      </router-link>
    </v-card-title>
    <v-card-text class="caption text-lg-left">
      <span v-if="videoItem.status === 3">{{ videoItem.duration }}</span> s. {{ videoItem.uploaded }} {{
        videoItem.views
      }} views
      <span v-if="showStatus"> {{
          videoStatus
        }}</span>
    </v-card-text>
    <div v-if="isCurrentUserOwner">
      <v-btn v-on:click="deleteItemPermanent(videoItem.id)" :data-id="videoItem.id">delete</v-btn>
    </div>
    <div v-if="status !== null">
      <span v-if="status">Success </span>
      <span v-else>Failed</span>
      <span>delete video
      </span>
    </div>
  </div>
</template>

<script>

import {mapActions, mapGetters} from "vuex";
import VideoStatus from "../store/videoStore/videoStatus";
import Cookie from "@/util/cookie";

export default {
  name: "VideoItem",
  props: [
    'video',
    'showStatus',
    'userPage',
  ],
  created() {
    this.currentUserId = Cookie.getCookie("userId");
    this.videoStatus = VideoStatus.intToStatus(this.videoItem.status)
  },
  data() {
    return {
      videoItem: this.video,
      videoStatus: null,
      isUserPage: this.userPage,
      status: null,
      currentUserId: null,
    }
  },
  methods: {
    ...mapActions({
      deleteVideoPermanent: "video/deleteVideoPermanent"
    }),
    ...mapGetters({
      getStatus: "video/getStatus"
    }),
    isCurrentUserOwner() {
      return this.currentUserId !== null && this.videoItem.ownerId === this.currentUserId
    },
    deleteItemPermanent(videoId) {
      const promise = this.deleteVideoPermanent({videoId: videoId});
      promise.then(() => {
        this.status = this.getStatus();
      })
    },
  }
}
</script>

<style scoped>

</style>