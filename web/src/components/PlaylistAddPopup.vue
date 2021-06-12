<template>
  <div>
    <button v-on:click="togglePopup">+ save</button>
    <button v-on:click="publish">publish</button>
    <div id="popup" class="popup hide">
      <div v-for="playlist in userPlaylists" :key="playlist.id">
        <div v-on:click="saveToPlaylist(playlist.id)">{{ playlist.name }}</div>
      </div>
      <button id="createPlaylist" v-on:click="showCreatePlaylistPopup">+ create new playlist</button>
      <div id="createPlaylistPopup" class="hide">
        <form @submit.prevent>
          <select name="privacy" id="privacy" v-model="privacy">
            <option disabled value="">Privacy type</option>
            <option value="2">Private</option>
            <option value="1">Unlisted</option>
            <option value="0">Public</option>
          </select>
          <div>
            <label for="playlistName"><span>Name</span>
              <input id="playlistName" type="text" v-model="name" placeholder="name" required/></label>
          </div>

          <v-btn @click="createPlaylist">create playlist</v-btn>
        </form>
      </div>
      <div></div>
    </div>
  </div>
</template>

<script>
import {mapActions, mapGetters} from "vuex";
import {toggleHideElement, unHideElement} from "@/util/dom";
import {publishEvent} from "@/events/event-bus";
import {playlistModification} from "@/store/playlistStore/actions";

export default {
  name: "PlaylistAddPopup",
  props: [
    "videoId"
  ],
  async mounted() {
    this.idVideo = this.videoId;
    this.createPlaylistPopupElement = document.getElementById("createPlaylistPopup");
    this.generalPopupElement = document.getElementById("popup");
    console.log(`this.generalPopupElement`);
    console.log(this.generalPopupElement);
    await this.fetchUserPlaylists();
    this.userPlaylists = this.getUserPlaylists();
  },
  data() {
    return {
      idVideo: null,
      show: false,
      userPlaylists: [],
      name: "",
      privacy: 0,
      errors: null,
      showNewPlaylistPopup: false,
      createPlaylistPopupElement: null,
      generalPopupElement: null,
    }
  },
  methods: {
    ...mapActions({
      createPlayList: "playlistMod/createPlaylist",
      doPlaylistModification: "playlistMod/modifyPlaylist",
      fetchUserPlaylists: "playlistMod/getUserPlaylists"
    }),
    ...mapGetters({
      getErrors: "playlistMod/getError",
      getStatus: "playlistMod/getSuccess",
      getMessage: "playlistMod/getResponseMessage",
      getUserPlaylists: "playlistMod/getPlaylists",
    }),
    publish() {
      publishEvent( 'success', "data" + Date.now())
    },
    async createPlaylist() {
      const videos = this.getCurrentVideos();
      await this.createPlayList({name: this.name, privacy: this.privacy, videos: videos});
      this.errors = this.getErrors();

      publishEvent(this.getStatus(), this.getMessage());
    },
    async saveToPlaylist(playlistId) {
      const videos = this.getCurrentVideos();
      await this.doPlaylistModification({action: playlistModification.AddVideo, playlistId: playlistId, videos})
      this.errors = this.getErrors();

      publishEvent(this.getStatus(), this.getMessage());
    },
    showCreatePlaylistPopup() {
      unHideElement(this.createPlaylistPopupElement);
    },
    togglePopup() {
      console.log('toggle popup')
      toggleHideElement(this.generalPopupElement);
    },
    getCurrentVideos() {
      const videos = [];
      videos.push(this.videoId);
      return videos;
    }
  }
}
</script>

<style scoped>
.hide {
  opacity: 0;
  display: none;
}

.popup {
  position: absolute;
  left: 200px;
  top: 200px;
  width: 200px;
  height: 400px;
  z-index: 1;
  -webkit-box-shadow: 5px 6px 15px -6px #E5E5E5;
  box-shadow: 5px 6px 15px -6px #E5E5E5;
  background-color: #e8dbdb;
}
</style>