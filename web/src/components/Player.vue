<template>
  <div>
    <video id="video" controls></video>
  </div>
</template>

<script>
let Hls = require('hls.js');

export default {
  name: "Player",
  props: [
    'videoId'
  ],
  mounted() {
    this.initPlayer()
  },
  updated() {
    this.initPlayer()
  },

  methods: {
    initPlayer() {
      const mediaStreamServerAddress = 'http://localhost:8000';
      let video = document.getElementById('video');
      if (Hls.isSupported()) {
        let hls = new Hls();
        hls.loadSource(mediaStreamServerAddress + '/api/v1/media/' + this.id + '/stream/');
        hls.attachMedia(video);
        hls.on(Hls.Events.MANIFEST_PARSED, function () {
          video.play();
        });
      } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
        video.src = mediaStreamServerAddress + '/api/v1/media/1/stream/';
        video.addEventListener('loadedmetadata', function () {
          video.play();
        });
      }
    }
  },
  data() {
    return {
      id: this.videoId
    }
  },
}
</script>

<style scoped>

</style>