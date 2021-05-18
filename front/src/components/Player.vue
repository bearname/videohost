<template>
  <div>
    <video id="video" controls></video>
  </div>
</template>

<script>
let Hls = require('hls.js');

export default {
  name: "Player",
  props: ['key',
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
      let video = document.getElementById('video');
      if (Hls.isSupported()) {
        let hls = new Hls();
        hls.loadSource('http://localhost:8000/api/v1/media/' + this.id + '/stream/');
        hls.attachMedia(video);
        hls.on(Hls.Events.MANIFEST_PARSED, function () {
          video.play();
        });
      } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
        video.src = 'http://localhost:8000/api/v1/media/1/stream/';
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