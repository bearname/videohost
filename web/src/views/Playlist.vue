<template>
  <v-container>
    <v-row v-if="playlists !== null">
      <div v-for="playlist in playlists">
        <div>
          <span class="subtitle-2" v-on:click="fetchPlaylist(playlist.id)">{{
              playlist.name
            }} privacy{{ playlist.privacy }}</span>
        </div>
      </div>
    </v-row>
    <v-row v-if="playlist!== null">
      {{ playlist.name }}
      {{ playlist.privacy }}
      <VideoList :show-status="true" :user-page="true" :videos="playlist.videos" :key="page"/>
    </v-row>
    <v-row v-if="videos!== null">
      {{ currentPlaylist }}
      <VideoList :show-status="true" :user-page="true" :videos="videos" :key="page"/>
    </v-row>
  </v-container>
</template>

<script>
import {mapActions, mapGetters} from "vuex";
import VideoList from "@/components/VideoList";
import {playlistModification} from "@/store/playlistStore/actions";
import {publishEvent} from "@/events/event-bus";

export default {
  name: "Playlist",
  components: {
    VideoList,
  },
  data() {
    return {
      videos: null,
      page: 0,
      countVideo: 30,
      currentPlaylist: null,
      playlists: null,
      playlist: null,
    }
  },
  async created() {
    this.currentPlaylist = this.$route.query.PL;
    console.log(`this.currentPlaylist`);
    console.log(this.currentPlaylist);
    console.log(Number.isInteger(parseInt(this.currentPlaylist)));
    if (this.currentPlaylist === "LL") {
      await this.fetchUserVideos(this.page, this.countVideo);
    } else if (this.isPlaylistId(this.currentPlaylist)) {
      await this.fetchPlaylist(this.currentPlaylist);
    } else {
      console.error("Unsupported list type " + this.currentPlaylist);
    }
    await this.fetchUserPlaylists();
    this.userPlaylists = this.getUserPlaylists();
  },
  methods: {
    ...mapActions({
      findUserLikedVideos: "userMod/getUserLikedVideos",
      getPlaylists: "playlistMod/getUserPlaylists",
      findPlaylist: "playlistMod/getPlaylist",
      fetchUserPlaylists: "playlistMod/getUserPlaylists",
      removeVideo: "playlistMod/getUserPlaylists",
      doPlaylistModification: "playlistMod/modifyPlaylist"
    }),
    ...mapGetters({
      getVideoResult: "userMod/getUserVideos",
      getPlaylistsRequestResult: "playlistMod/getPlaylists",
      getPlaylistResult: "playlistMod/getPlaylist",
      getUserPlaylists: "playlistMod/getPlaylists",
    }),
    async fetchUserVideos(page, countVideo) {
      await this.findUserLikedVideos({page: page, countVideoOnPage: countVideo});
      const result = this.getVideoResult();

      this.videos = result.videos;
    },
    async fetchUserPlaylists() {
      await this.getPlaylists();
      this.playlists = this.getPlaylistsRequestResult();
    },
    async fetchPlaylist(playlistId) {
      await this.findPlaylist({playlistId: playlistId});
      this.playlist = this.getPlaylistResult();
      console.log(`      this.playlist = this.getPlaylist();`);
      console.log(this.playlist);
    },
    async removeVideoFromPlaylist(playlistId, videoId) {
      const videos = [];
      videos.push(videoId);

      await this.doPlaylistModification({action: playlistModification.RemoveVideo, playlistId: playlistId, videos})
      this.errors = this.getErrors();

      publishEvent(this.getStatus(), this.getMessage());
    },
    isPlaylistId(currentPlaylist) {
      return Number.isInteger(parseInt(currentPlaylist));
    },
  },
}
</script>

<style scoped>

</style>