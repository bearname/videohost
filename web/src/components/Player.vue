<template>
  <div>
    <div>
      <video id="video" controls></video>
    </div>
    <div v-if="qualities !== null">
      <select v-model="selectedQuality">
        <option disabled value="">Please select one</option>
        <option v-for="quality in qualities" :key="quality" :value="quality">{{ quality }}p</option>
        <!--        <option value="720">720p</option>-->
        <!--        <option value="480">480p</option>-->
        <!--        <option value="360">360p</option>-->
      </select>
    </div>
  </div>
</template>

<script>
let Hls = require('hls.js');

export default {
  name: "Player",
  props: [
    'videoId',
    'availableQualities'
  ],
  created() {
    this.qualities = this.availableQualities.split(",");
    console.log(this.qualities)
  },
  data() {
    return {
      id: this.videoId,
      selectedQuality: '1080',
      qualities: null
    }
  },
  mounted() {
    this.initPlayer()
  },
  updated() {
    this.initPlayer()
  },
  methods: {
    initPlayer() {
      let video = document.getElementById('video');
      if (Hls.isSupported()) {
        let hls = new Hls();
        hls.loadSource(process.env.VUE_APP_VIDEO_SERVER_ADDRESS + '/media/' + this.id + '/' + this.selectedQuality + '/stream/');
        hls.attachMedia(video);
        hls.on(Hls.Events.MANIFEST_PARSED, function () {
          video.play();
        });
      } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
        video.src = process.env.VUE_APP_VIDEO_SERVER_ADDRESS + '/media/' + this.id + '/' + this.selectedQuality + '/stream/';
        video.addEventListener('loadedmetadata', function () {
          video.play();
        });
      }
    }
  },
}
</script>

<style scoped>

</style>