<template>
  <div>
    <div v-if="qualities !== null">
      <div id="videoWrapper" class="player-wrapper player-medium">
        <video id="video" width="720px" autoplay="autoplay" :poster="poster" class="player-video player-medium"></video>
        <div id="videoControls" class="player-controls">
          <div id="progressBarWrapper" class="wrapper-bar player-medium" v-on:click="onClickBufferedRange($event);">
            <div id="progressBar" class="progress-bar player-medium"></div>
          </div>
          <div id="progressBarCircle" class="progress-bar-circle"></div>
          <div>
            <span id="showHoverTime" class="hide"></span>
            <span class="player-controls-item">
              <button
                  id="playButton"
                  type="button"
                  title="video.play()"
                  v-on:click="togglePlay()">&#x23F5;</button>
            </span>
            <span class="player-controls-item">
              <button id="volumeMute">mute</button>
            </span>
            <span class="player-controls-item">
              <input id="volumeChanger" type="range" name="volume" class="player-volume" min="0" max="1" step="0.05"
                     value="1">
              <span id="volume" class="player-controls-item"
                    v-if="this.videoElement !== null">{{ this.videoElement.volume * 100 }}       </span>
              <span v-else>100%       </span>
            </span>
            <span>
              <span id="currentTime"></span>
              <span> / {{ videoDuration }}</span>
            </span>
            <div class="float-right ">
              <span id="settingButton" class="gear player-controls-item "></span>
              <button v-on:click="togglePlayerSize" id="changePlayerSize" class="player-controls-item">&#9645;</button>
              <button v-on:click="toggleFullScreen" class="float-right player-controls-item">&#x26F6;</button>
            </div>
          </div>
          <div id="settings" class="settings-popup hide">
            <div>
              <span id="playbackSpeedSetting">Playback speed</span>
              <select name="playSpeed" id="playSpeed" v-on:input="setPlaybackSpeed()">
                <option disabled value="">Playback Speed</option>
                <option value="0.25">0.25x</option>
                <option value="0.50">0.5x</option>
                <option value="0.75">0.75x</option>
                <option value="1" selected="selected">Normal</option>
                <option value="1.25">1.25x</option>
                <option value="1.50">1.50x</option>
                <option value="1.75">1.75x</option>
                <option value="2">2x</option>
              </select>
            </div>
            <div id="qualitySetting">
              <span>Quality: </span>
              <select id="selectQuality" v-on:input="changeQuality()">
                <option disabled value="">Please select one</option>
                <option value="-1">auto</option>
                <option v-for="(quality, index) in qualities" :key="quality" :value="index">{{ index }} ! {{
                    quality
                  }}p
                </option>
              </select>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-else>
      Video not available
    </div>
    <ul id="videoList" v-if="videos !== null">
      <li v-for="video in videos" :key="video.id" :data-videoId="video.src">{{ video.name }}</li>
    </ul>
  </div>
</template>

<script>
import {mapActions, mapGetters} from "vuex";

let Hls = require('hls.js');

export default {
  name: "Player",
  props: [
    'videoId',
    'availableQualities',
    'duration',
    'thumbnail',
  ],
  created() {
    this.qualities = this.availableQualities.split(",");
    console.log(this.qualities)
    if (this.qualities[0] === "") {
      this.qualities = null
    }
    this.qualities.reverse()
    console.log("this.qualities");
  },
  data() {
    return {
      id: this.videoId,
      selectedQuality: '1080',
      qualities: null,
      videoDuration: this.formatTimeString(this.duration),
      isMedium: true,
      previousVolume: 1,
      playbackRate: null,
      shiftTime: 10,
      volume: null,
      hls: null,
      isPause: false,
      videoElement: null,
      playbackSpeed: null,
      isFullScreen: false,
      bufferingIdx: null,
      lastStartPosition: 0,
      bufferedCanvas: null,
      poster: this.thumbnail,
      videoControlElement: null,
      playButton: null,
      pauseButton: null,
      changePlayerSizeButton: null,
      selectQualityElement: null,
      videoWrapperElement: null,
      currentTimeElement: null,
      volumeMuteElement: null,
      showHoverTimeElement: null,
      progressBarWrapperElement: null,
      progressBarElement: null,
      progressBarCircleElement: null,
      volumeChangerElement: null,
      settingButtonElement: null,
      settingsPopupElement: null,
      playbackSpeedSettingElement: null,
      qualitySettingElement: null,
      videos: null,
      volumePrevious: 1,
      showSetting: false,
    }
  },
  mounted() {
    this.initPlayer()
    this.initKeyHandler()
    this.fetchVideosByPage(0, 10)
  },
  updated() {
    this.initPlayer()
  },
  methods: {
    ...mapActions({
      getVideoOnPage: "video/getVideoOnPage"
    }),
    ...mapGetters({
      getVideos: "video/getVideos",
      getPageCount: "video/getPageCount"
    }),
    async fetchVideosByPage(page, countVideoOnPage) {
      await this.getVideoOnPage(page, countVideoOnPage)
          .then(() => {
            this.videos = this.getVideos()
            console.log(this.videos)
            this.countPage = this.getPageCount()
          })
    },
    onClickBufferedRange(event) {
      this.videoElement.currentTime = this.getEventMouseTime(event);
    },

    initPlayer() {
      this.videoWrapperElement = document.getElementById('videoWrapper');
      this.videoElement = document.getElementById('video');
      this.videoElement.addEventListener('click', () => {
        this.togglePlay()
      })
      this.volume = document.getElementById('volume');
      this.bufferedCanvas = document.getElementById('bufferedCanvas')
      this.videoControlElement = document.getElementById('videoControls')
      this.playButton = document.getElementById('playButton')
      this.pauseButton = document.getElementById('pauseButton')
      this.changePlayerSizeButton = document.getElementById('changePlayerSize')
      this.selectQualityElement = document.getElementById('selectQuality')
      this.currentTimeElement = document.getElementById('currentTime')
      this.volumeMuteElement = document.getElementById('volumeMute')
      this.progressBarWrapperElement = document.getElementById('progressBarWrapper')
      this.progressBarElement = document.getElementById('progressBar')
      this.showHoverTimeElement = document.getElementById('showHoverTime')
      this.progressBarCircleElement = document.getElementById('progressBarCircle')
      this.volumeChangerElement = document.getElementById('volumeChanger')
      this.settingButtonElement = document.getElementById('settingButton')
      this.qualitySettingElement = document.getElementById('qualitySetting')
      this.playbackSpeedSettingElement = document.getElementById('playbackSpeedSetting')
      this.settingsPopupElement = document.getElementById('settings')

      this.videoElement.addEventListener('timeupdate', this.handleTimeUpdate, false)
      this.videoElement.addEventListener('stop', this.onVideoStop, false)
      this.volumeMuteElement.addEventListener('click', this.toggleMute, false)

      this.progressBarWrapperElement.addEventListener('mousemove', this.onHoverProgressBar, false)
      this.progressBarWrapperElement.addEventListener('mouseout', this.hideShowTime, false)
      this.volumeChangerElement.addEventListener('change', this.handleVolumeChange, false)
      this.volumeChangerElement.addEventListener('mousemove', this.handleVolumeChange, false)

      this.settingButtonElement.addEventListener('click', this.toggleSettingPopup, false)
      console.log(this.videoElement.getVideoPlaybackQuality())
      this.initHls();
    },
    togglePlay() {
      if (!this.videoElement.paused && !this.videoElement.ended) {
        this.pause();
      } else {
        this.play();
      }
    },
    initHls() {
      if (Hls.isSupported()) {
        this.prepareHls(this.id, this.selectedQuality);
      } else if (this.videoElement.canPlayType('application/vnd.apple.mpegurl')) {
        this.videoElement.src = process.env.VUE_APP_VIDEO_SERVER_ADDRESS + '/media/' + this.id + '/stream/';
        this.videoElement.addEventListener('loadedmetadata', function () {
          this.video.play();
          console.log('getVideoPlaybackQuality()')
        });
      } else {
        let videoWrapper = document.getElementById('videoWrapper');
        videoWrapper.innerText = "Not support streaming video"
      }
    },
    prepareHls(id, quality) {
      if (this.hls !== null) {
        this.hls.destroy();
        this.hls = null;
      }
      console.log(quality)
      this.hls = new Hls();
      this.hls.loadSource(process.env.VUE_APP_VIDEO_SERVER_ADDRESS + '/media/' + id + '/stream/');
      this.hls.attachMedia(this.videoElement);
      this.hls.on(Hls.Events.MANIFEST_PARSED, () => {
        this.play();
        console.log('getVideoPlaybackQuality()')
      });
      console.log('frag_buffered')
      this.hls.on(Hls.Events.FRAG_BUFFERED, () => {
        console.log('frag,et buffered')
        // this.hls.bufferTimer = setInterval(this.checkBuffer, 1000);
      })
      // this.hls.on(Hls.Events.MEDIA_DETACHED, function () {
      //   clearInterval(this.hls.bufferTimer);
      // });
      // this.hls.on(Hls.Events.DESTROYING, function () {
      //   clearInterval(this.hls.bufferTimer);
      // });
      // this.hls.on(Hls.Events.BUFFER_RESET, function () {
      //   clearInterval(this.hls.bufferTimer);
      // });
      console.log(this.hls)
    },
    initKeyHandler() {
      console.log('ADDING EVENT HANDLER');
      window.addEventListener("keydown", e => {
        const key = e.key;
        console.log('key')
        console.log(key)
        switch (key) {
          case 'j' || 'J': {
            this.shiftCurrentTime(-this.shiftTime);
            break
          }
          case 'f' || 'F': {
            this.toggleFullScreen()
            break
          }
          case '<': {
            this.shiftPlaybackSpeed(-0.25);
            break
          }
          case '>': {
            this.shiftPlaybackSpeed(0.25);
            break
          }
          case  'l' || 'L' : {
            this.shiftCurrentTime(this.shiftTime);
            break
          }
          case 'k' || 'K': {
            this.togglePlay()
            break;
          }
          case 'ArrowLeft': {
            this.shiftCurrentTime(-5)
            break
          }
          case 'ArrowRight': {
            this.shiftCurrentTime(5)
            break
          }
          case 'ArrowUp': {
            if (this.videoElement.volume + 0.05 <= 1) {
              e.preventDefault()
              this.videoElement.volume += 0.05
              this.updateVolumeText()
            }
            break
          }
          case 'ArrowDown': {
            if (this.videoElement.volume > 0) {
              e.preventDefault()
              this.videoElement.volume -= 0.05;
              this.updateVolumeText();
            }
            break
          }
          case 'm' || 'M': {
            if (this.videoElement.volume > 0) {
              this.previousVolume = this.videoElement.volume
              this.videoElement.volume = 0;
            } else if (this.videoElement.volume === 0) {
              this.videoElement.volume = this.previousVolume
              this.previousVolume = 0
            }
            break;
          }
          default: {
            const eventCode = e.code;

            console.log(eventCode)
            if (eventCode === 'Space') {
              e.preventDefault()
              this.togglePlay()
            } else {
              const s = eventCode.substring(0, eventCode.length - 1);
              const keyElement = eventCode[eventCode.length - 1];
              console.log(s)
              console.log(keyElement)

              if ((s === 'Numpad' || s === 'Digit') && keyElement >= '0' && keyElement <= '9') {
                const number = parseInt(keyElement);
                console.log(number / 10)
                const shift = this.videoDuration * (number / 10);
                console.log(shift)
                this.setCurrentTime(shift)
              }
            }

            break
          }
        }
      }, false);
    },
    toggleFullScreen() {
      if (!this.isFullScreen) {
        this.openFullscreen()
      } else {
        this.closeFullscreen()
      }
    },
    openFullscreen() {
      let elem = this.videoWrapperElement;
      if (elem.requestFullscreen) {
        elem.requestFullscreen();
        this.toggleFullScreenOnVideoElement()
      } else if (elem.webkitRequestFullscreen) { /* Safari */
        elem.webkitRequestFullscreen();
        this.toggleFullScreenOnVideoElement()
      } else if (elem.msRequestFullscreen) { /* IE11 */
        elem.msRequestFullscreen();
        this.toggleFullScreenOnVideoElement()
      }
    },
    closeFullscreen() {
      if (document.exitFullscreen) {
        document.exitFullscreen();
        this.toggleFullScreenOnVideoElement()
      } else if (document.webkitExitFullscreen) { /* Safari */
        document.webkitExitFullscreen();
        this.toggleFullScreenOnVideoElement()
      } else if (document.msExitFullscreen) { /* IE11 */
        document.msExitFullscreen();
        this.toggleFullScreenOnVideoElement()
      }
    },
    toggleFullScreenOnVideoElement() {
      if (!this.isFullScreen) {
        this.changePlayerSizeButton.classList.add('hide')
        this.isFullScreen = true;
      } else {
        this.isFullScreen = false;
        this.changePlayerSizeButton.classList.remove('hide')
      }
    },
    togglePlayerSize() {
      if (this.isMedium) {
        this.isMedium = false
        this.videoElement.classList.remove('player-medium')
        this.videoElement.classList.add('player-big')
        this.videoWrapperElement.classList.remove('player-medium')
        this.videoWrapperElement.classList.add('player-big')
        this.progressBarElement.classList.remove('player-medium')
        this.progressBarElement.classList.add('player-big')

        this.bufferedCanvas.width = '1327px'
      } else {
        this.isMedium = true
        this.videoElement.classList.add('player-medium')
        this.videoElement.classList.add('player-medium')
        this.videoWrapperElement.classList.add('player-medium')
        this.videoWrapperElement.classList.remove('player-big')
        this.progressBarElement.classList.add('player-medium')
        this.progressBarElement.classList.remove('player-big')
        this.bufferedCanvas.width = '992px'
      }
    },
    changeQuality() {
      console.log(this.hls)

      let loadLevelNumber = this.selectQualityElement.value;
      console.log(loadLevelNumber);
      this.hls.loadLevel = parseInt(loadLevelNumber)
    },
    play() {
      this.videoElement.play()
      this.isPause = false
      this.playButton.innerText = '❚ ❚'
    },
    pause() {
      this.videoElement.pause()
      this.isPause = true
      this.playButton.textContent = '►'
      this.videoControlElement.style.transform = 'translateY(-10px)'
    },
    formatTimeString(seconds) {
      let date = new Date(0);

      date.setSeconds(seconds);
      let timeString = date.toISOString().substr(11, 8);

      if (seconds < 3600) {
        timeString = timeString.substring(3, timeString.length)
      }

      return timeString
    },
    setPlaybackSpeed() {
      this.playbackRate = document.getElementById("playSpeed");
      const playbackRate = this.playbackRate.value;
      if (playbackRate >= 0.25 && playbackRate <= 2) {
        this.videoElement.defaultPlaybackRate = playbackRate;
        this.videoElement.playbackRate = playbackRate;
      }
    },
    shiftPlaybackSpeed(shift) {
      if (this.videoElement.defaultPlaybackRate + shift >= 0.25 && this.videoElement.defaultPlaybackRate + shift <= 2) {
        this.videoElement.defaultPlaybackRate += shift;
        this.videoElement.playbackRate += shift;

        if (this.playbackSpeed === null) {
          this.playbackSpeed = document.getElementById('playSpeed');
        }

        const children = this.playbackSpeed.children;
        children.forEach(function (part, index) {
          const optionElement = this[index];
          const attribute = optionElement.getAttribute('selected');
          if (attribute !== null) {
            optionElement.removeAttribute('selected')
          }
        }, children)

        const defaultPlaybackRate = this.videoElement.playbackRate;
        for (let i = 0; i < children.length; i++) {
          const optionElement = children[i];
          const value = optionElement.getAttribute('value');
          if (value !== null) {
            const number = parseFloat(value);
            console.log(number)
            if (number === defaultPlaybackRate) {
              children[i].setAttribute('selected', 'selected')
              break;
            }
          }
        }
      }
    },
    shiftCurrentTime(shift) {
      this.videoElement.currentTime += shift
    },
    setCurrentTime(time) {
      this.videoElement.currentTime = time
    },
    updateVolumeText() {
      this.volume.innerText = Math.ceil(this.videoElement.volume * 100) + '%'
    },
    handleTimeUpdate() {
      this.updateTime()
      this.updateProgressBar()
    },
    updateTime() {
      this.currentTimeElement.innerText = this.formatTimeString(this.videoElement.currentTime)
    },
    updateProgressBar() {
      let boundingClientRect = this.videoElement.getBoundingClientRect();
      let size = parseInt(this.videoElement.currentTime * boundingClientRect.width / this.duration);
      this.progressBarElement.style.width = size + 'px';
      this.progressBarCircleElement.style.left = size + 'px'
    },
    toggleMute() {
      if (this.videoElement.muted) {
        this.videoElement.volume = this.volumePrevious
        this.videoElement.muted = false
        this.changeButtonType(this.volumeMuteElement, 'mute')
        this.volumeChangerElement.value = this.videoElement.volume
      } else {
        this.volumePrevious = this.videoElement.volume
        this.videoElement.muted = true
        this.changeButtonType(this.volumeMuteElement, 'unmute')
        this.volumeChangerElement.value = 0
      }

      this.handleVolumeChange()
    },
    changeButtonType(btn, value) {
      btn.title = value;
      btn.innerHTML = value;
      btn.className = value;
    },
    onHoverProgressBar(event) {
      let eventMouseTime = this.getEventMouseTime(event);
      let boundingClientRect = this.progressBarWrapperElement.getBoundingClientRect();
      let x = event.x - boundingClientRect.x;
      this.showHoverTimeElement.style.left = x - 15 + 'px'
      this.showHoverTimeElement.style.top = -20 + 'px'
      this.showHoverTimeElement.classList.remove('hide')
      this.showHoverTimeElement.innerText = this.formatTimeString(eventMouseTime);
    },
    getEventMouseTime(event) {
      let boundingClientRect = this.progressBarWrapperElement.getBoundingClientRect();

      return (event.pageX - boundingClientRect.left) * this.duration / boundingClientRect.width;
    },
    hideShowTime() {
      this.showHoverTimeElement.classList.add('hide')
    },
    onVideoStop() {
      this.videoControlElement.style.transform = 'translateY(-20px)'
    },
    handleVolumeChange() {
      this.videoElement.volume = this.volumeChangerElement.value
      if (this.videoElement.volume === 0) {
        this.volumeMuteElement.innerText = "unmute"
        this.volumeChangerElement.value = 0
        this.videoElement.muted = true
      } else {
        this.videoElement.muted = false
        this.volumeMuteElement.innerText = "mute"
        this.volumeChangerElement.value = this.volumePrevious

      }
      this.updateVolumeText();
    },
    toggleSettingPopup() {
      let boundingClientRect = this.settingButtonElement.getBoundingClientRect();
      this.settingsPopupElement.style.top = boundingClientRect.top - 300
      this.settingsPopupElement.style.left = boundingClientRect.left + 100
      this.settingsPopupElement.classList.toggle('hide')
    },
  },
}
</script>

<style scoped>

select {
  color: #e8dbdb;
  background: rgba(28, 28, 28, 0.9);
}

.float-right {
  float: right;
}

.player-video {
  width: 100%;
}

.player-big {
  width: 1327px;
}

.player-medium {
  width: 992px;
}

.player-wrapper:fullscreen, .player-video {
  max-width: none;
  width: 100%;
  margin: auto 0;
}

.player-wrapper:-webkit-full-screen {
  max-width: none;
  width: 100%;
}

.hide {
  opacity: 0;
}

.player-wrapper {
  position: relative;
  overflow: hidden;
}

.player-controls {
  /*display: flex;*/
  position: absolute;
  bottom: 0;
  width: 100%;
  transform: translateY(100px);
  transition-delay: 0.5s;
  transition: all .3s;
  flex-wrap: wrap;
  background: rgba(0, 0, 0, 0.5);
  color: #e8dbdb;
  opacity: 0;
  padding: 0 4px 10px;
}

.player-controls-item {
  padding: 10px;
}

.player-wrapper:hover .player-controls {
  transform: translateY(-10px);
  opacity: 1;
}

.player-wrapper:hover .progress-bar {
  height: 4px;
}

.wrapper-bar {
  position: relative;
  float: left;
  width: 100%;
  height: 3px;
  background: rgba(255, 255, 255, 0.4);
  margin-bottom: 14px;
}

.progress-bar,
.wrapper-bar {
  cursor: pointer;
  position: relative;
  height: 3px;
}

.progress-bar {
  background-color: #f00;
  height: 3px;
}

.wrapper-bar:hover {
  height: 4px;
}

.wrapper-bar:hover .progress-bar {
  height: 5px;
}

.wrapper-bar:hover, .progress-bar-circle {
  opacity: 1;
}

.progress-bar-circle {
  border-radius: 50%;
  width: 14px;
  height: 14px;
  background-color: #f00;
  position: absolute;
  top: -5px;
  left: 0;
}

#showHoverTime {
  position: absolute;
  z-index: 10;
  color: #e8dbdb;
  text-shadow: #9c9b9b;
}

.player-volume,
.player-volume[value] {
  appearance: none;
  -moz-appearance: none;
  -webkit-appearance: none;
  width: 52px;
  height: 4px;
  color: white;
  background-color: white;
  border-radius: 2px;
  background-size: 35px 20px, 100% 100%, 100% 100%;
}

.gear {
  display: inline-block;
  position: relative;
  margin: 0.25em;
  width: 1em;
  height: 1em;
  background: white;
  border-radius: 50%;
  border: 0.3em solid gray;
  box-sizing: border-box;
}

.gear:before,
.gear:after {
  content: "×";
  position: absolute;
  z-index: -1;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-weight: bold;
  font-size: 2.5em;
  color: white;
}

.gear:after {
  transform: translate(-50%, -50%) rotate(45deg);
}

.gear:hover {
  cursor: pointer;
}

.settings-popup {
  position: absolute;
  top: -100px;
  left: 600px;
  z-index: 100;
  padding: 30px 30px;
  background: rgba(28, 28, 28, 0.9);
  border-radius: 4px;
}

</style>