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
      {{ videoItem.duration }} s. {{ videoItem.uploaded }}  {{videoItem.views}} views
    </v-card-text>
  </div>
</template>

<script>
export default {
  name: "VideoItem",
  props: [
    'video'
  ],
  data() {
    return {
      videoItem: this.video
    }
  },
}
</script>

<style scoped>

</style>