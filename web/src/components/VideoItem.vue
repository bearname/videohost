<template>
  <v-card style="overflow-y: auto; height:340px">
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
    <div v-if="status !== null">
      <span v-if="status">Success </span>
      <span v-else>Failed</span>
      <span>delete video
      </span>
    </div>
  </v-card>
</template>

<script>

import {mapActions, mapGetters} from "vuex";
import VideoStatus from "../store/videoStore/videoStatus";
import {VBtn, VCard, VCardActions, VCardText, VCardTitle, VDialog, VImg, VToolbar} from "vuetify/lib";

export default {
  name: "VideoItem",
  components: {
    VBtn,
    VCard,
    VCardTitle,
    VCardText,
    VToolbar,
    VCardActions,
    VImg,
    VDialog,
  },
  props: [
    'video',
    'showStatus',
    'userPage',
    'currentUserId',
  ],
  created() {
    this.videoStatus = VideoStatus.intToStatus(this.videoItem.status)
  },
  data() {
    return {
      videoItem: this.video,
      videoStatus: null,
      isUserPage: this.userPage,
      status: null,
      userId: this.currentUserId,
      saveToPlaylist: false,
      userPlaylists: null,
      dialog: false,
      select: [
        { text: 'State 1' },
        { text: 'State 2' },
        { text: 'State 3' },
        { text: 'State 4' },
        { text: 'State 5' },
        { text: 'State 6' },
        { text: 'State 7' },
      ],

    }
  },
  methods: {
    ...mapActions({
      deleteVideoPermanent: "videoMod/deleteVideoPermanent"
    }),
    ...mapGetters({
      getStatus: "videoMod/getStatus"
    }),
    isCurrentUserOwner() {
      return this.userId !== null && this.videoItem.ownerId === this.userId
    },
    async deleteItemPermanent(videoId) {
      await this.deleteVideoPermanent({videoId: videoId});
      this.status = this.getStatus();
    },
    toggleDialog() {
      console.log('dialog')
      this.dialog = true;
    }
  }
}
</script>

<style scoped>

</style>
