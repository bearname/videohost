<template>
  <div>
    <div v-if="qualities !== null">
      <div id="videoWrapper" class="video-wrapper">
        <video id="video" width="720px" autoplay="autoplay" :poster="poster" class="player-medium"></video>
        <!--        <div id="videoControls" class="hide video-controls">-->
        <div id="videoControls">
          <div id="defaultBar">
            <div id="progressBar"></div>
          </div>
          <div style="clear:both"></div>
          <!--          <canvas-->
          <!--              id="bufferedCanvas"-->
          <!--              width="720"-->
          <!--              height="15"-->
          <!--              class="videoCentered"-->
          <!--              v-on:click="onClickBufferedRange($event);"-->
          <!--              style="height: fit-content"-->
          <!--          ></canvas>-->
          <div>
            <span>
              <button
                  id="playButton"
                  type="button"
                  title="video.play()"
                  v-on:click="playOrPause()">Play</button>
            </span>
            <span>
              <span id="currentTime"></span>
              <span></span>
            </span>
            <span>
              <span>Playback speed</span>
              <select name="playSpeed" id="playSpeed" v-on:input="setPlaybackSpeed()">
                <option disabled value="">Playback Speed</option>
                <option value="0.25">0.25x</option>
                <option value="0.50">0.5x</option>
                <option value="0.75">0.75x</option>
                <option value="1" selected="selected">Normal</option>
                <option value="1.25">1.25x</option>
                <option value="1.50">1.50x</option>
                <option value="1.75">1.75x</option>
              </select>
            </span>
            <span>
            <button
                type="button"
                class="btn btn-sm btn-info"
                title="video.currentTime -= 10"
                v-on:click="shiftCurrentTime(-10)"
            >
              - 10 s
            </button>
            <button
                type="button"
                class="btn btn-sm btn-info"
                title="video.currentTime += 10"
                v-on:click="shiftCurrentTime(10)"
            >
              + 10 s
            </button>
          </span>
            <span>
            <button
                type="button"
                class="btn btn-xs btn-warning"
                title="hls.startLoad()"
                v-on:click="startLoadHls()"
            >
              Start loading
            </button>
            <button
                type="button"
                class="btn btn-xs btn-warning"
                title="hls.stopLoad()"
                v-on:click="stopLoadHls()"
            >
              Stop loading
            </button>
          </span>
            <span>Volume
            <span id="volume" v-if="this.videoElement !== null">{{ this.videoElement.volume * 100 }}</span>
            <span v-else>100%</span>
          </span>
            <button v-on:click="toggleFullScreen">full screen</button>
          </div>
          <span>
            <span>Quality: </span>
            <select id="selectQuality" v-on:input="changeQuality()">
              <option disabled value="">Please select one</option>
              <option value="-1">auto</option>
              <option v-for="(quality, index) in qualities" :key="quality" :value="index">{{ index }} ! {{
                  quality
                }}p</option>
            </select>
          </span>
          <span>
            <button v-on:click="togglePlayerSize" id="changePlayerSizeButton" class="vi">Увеличить</button>
          </span>
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
      videoDuration: this.duration,
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
      videos: null,
      updateBar: null,
      defaultBar: null,
      progressBar: null,
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
      if (this.bufferedCanvas === null) {
        this.bufferedCanvas = document.querySelector('#bufferedCanvas');
      }
      console.log(this.bufferedCanvas.offsetLeft)

      this.videoElement.currentTime = ((event.clientX - this.bufferedCanvas.offsetLeft) / this.bufferedCanvas.width) * this.getSeekableEnd();
    },
    getSeekableEnd() {
      if (isFinite(this.videoElement.duration)) {
        return this.videoElement.duration;
      }
      if (this.videoElement.seekable.length) {
        return this.videoElement.seekable.end(this.videoElement.seekable.length - 1);
      }
      return 0;
    },
    checkBuffer() {
      // if (this.isPause) {
      //   return
      // }
      // const buffered = this.videoElement.buffered;
      //
      // const seekableEnd = this.getSeekableEnd();

      // console.log('seekableEnd')
      // console.log(seekableEnd)
      // let bufferingDuration;
      // const ctx = this.bufferedCanvas.getContext('2d');
      //
      // if (buffered) {
      // console.log('buffered')
      // console.log(buffered)

      // const pos = this.video.currentTime;
      // let bufferLen = 0;
      // ctx.fillStyle = '#03bb2b';
      //
      // for (let i = 0; i < buffered.length; i++) {
      //   const start = (buffered.start(i) / seekableEnd) * this.canvas.width;
      //   const end = (buffered.end(i) / seekableEnd) * this.canvas.width;
      //   ctx.fillRect(start, 2, Math.max(2, end - start), 11);
      //   if (pos >= buffered.start(i) && pos < buffered.end(i)) {
      //     // play position is inside this buffer TimeRange, retrieve end of buffer position and buffer length
      //     bufferLen = buffered.end(i) - pos;
      //   }
      // }
      // console.log(bufferLen);
      // } else if (ctx.fillStyle !== '#000') {
      //   ctx.fillStyle = '#000';
      //   ctx.fillRect(0, 0, this.bufferedCanvas.width, this.bufferedCanvas.height);
      // }
      // let date = new Date(0);
      // let sec = this.videoElement.currentTime;
      // date.setSeconds(sec);
      // let timeString = date.toISOString().substr(11, 8);
      //
      // if (sec < 3600) {
      //   timeString = timeString.substring(2, timeString.length)
      // }
      // console.log(timeString)
      // this.currentTimeElement.innerText = timeString
      //
      // ctx.fillStyle = '#0511b6';
      // const x = this.videoElement.currentTime / seekableEnd * this.bufferedCanvas.width;
      // ctx.fillRect(x, 0, 2, 15);
    },
    getCurrentTime(seekableEnd) {
      return (this.videoElement.currentTime / seekableEnd) * this.bufferedCanvas.width
    },
    initPlayer() {
      this.videoWrapperElement = document.getElementById('videoWrapper');
      this.videoElement = document.getElementById('video');
      this.videoElement.addEventListener('click', () => {
        // this.togglePause()
        this.playOrPause()
      })
      this.volume = document.getElementById('volume');
      this.bufferedCanvas = document.getElementById('bufferedCanvas')
      this.videoControlElement = document.getElementById('videoControls')
      this.playButton = document.getElementById('playButton')
      this.pauseButton = document.getElementById('pauseButton')
      this.changePlayerSizeButton = document.getElementById('changePlayerSize')
      this.selectQualityElement = document.getElementById('selectQuality')
      this.currentTimeElement = document.getElementById('currentTime')
      // this.video.addEventListener('mouseover', (e) => {
      //   e.preventDefault()
      //   this.videoControlElement.classList.remove('hide');
      // })
      // this.video.addEventListener('mouseout', () => {
      //   if (!this.video.isPause) {
      //     this.videoControlElement.classList.add('hide');
      //   }
      // })
      console.log(this.videoElement.getVideoPlaybackQuality())
      this.initHls();
    },

    playOrPause() {
      if (!this.videoElement.paused && !this.videoElement.ended) {
        this.videoElement.pause();
        this.playButton.innerHTML = 'Play';
        window.clearInterval(this.updateBar);
      } else {
        this.videoElement.play();
        this.playButton.innerHTML = 'Pause';
        this.updateBar = setInterval(this.update, 500);
      }
    },
    update() {
      if (!this.videoElement.ended) {
        let size = parseInt(this.videoElement.currentTime * this.videoElement.width / this.videoElement.duration);
        this.progressBar.style.width = size + 'px';
      } else {
        this.progressBar.style.width = '0px';
        this.playButton.innerHTML = 'Play';
        window.clearInterval(this.updateBar);
      }
    },
    clickedBar(e) {
      if (!this.videoElement.paused && !this.videoElement.ended) {
        let mouseX = e.pageX - this.defaultBar.offsetLeft;
        this.videoElement.currentTime = mouseX * this.videoElement.duration / this.videoElement.width;
        this.progressBar.style.width = mouseX + 'px';
      }
    },
    initProgressBar() {
      this.defaultBar = document.getElementById('defaultBar');
      this.progressBar = document.getElementById('progressBar');

      this.playButton.addEventListener('click', this.playOrPause, false);
      this.defaultBar.addEventListener('click', this.clickedBar, false);
    },
    initHls() {
      if (Hls.isSupported()) {
        this.prepareHls(this.id, this.selectedQuality);
      } else if (this.videoElement.canPlayType('application/vnd.apple.mpegurl')) {
        this.videoElement.src = process.env.VUE_APP_VIDEO_SERVER_ADDRESS + '/media/' + this.id + '/stream/';
        this.videoElement.addEventListener('loadedmetadata', function () {
          this.video.play();
          this.initProgressBar()
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
        clearInterval(this.hls.bufferTimer);
        this.hls = null;
      }
      console.log(quality)
      this.hls = new Hls();
      this.hls.loadSource(process.env.VUE_APP_VIDEO_SERVER_ADDRESS + '/media/' + id + '/stream/');
      this.hls.attachMedia(this.videoElement);
      this.hls.on(Hls.Events.MANIFEST_PARSED, () => {
        this.videoElement.play();
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
            // this.togglePause()
            this.playOrPause()
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
              this.setVolumeText()
            }
            break
          }
          case 'ArrowDown': {
            if (this.videoElement.volume > 0) {
              e.preventDefault()
              this.videoElement.volume -= 0.05;
              this.setVolumeText();
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
              // this.togglePause()
              this.playOrPause()
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
    togglePause() {
      this.playOrPause()
      // if (this.isPause) {
      //   this.play()
      // } else {
      //   this.pause()
      // }
    },
    toggleFullScreen() {
      // if (!this.isFullScreen) {
      //   this.isFullScreen = true;
      //   console.log( this.screen)
      //   this.videoElement.width = this.screen.width;
      //   this.videoElement.height = this.screen.height;
      // } else {
      //   this.isFullScreen = false;
      // }
      if (!this.isFullScreen) {
        if (this.video.mozRequestFullScreen) {
          this.video.mozRequestFullScreen();
        } else {
          this.video.webkitRequestFullScreen(Element.ALLOW_KEYBOARD_INPUT);
        }
        this.isFullScreen = true;
      } else {
        if (document.mozCancelFullScreen) {
          document.mozCancelFullScreen();
        } else {
          document.webkitCancelFullScreen();
        }
        this.isFullScreen = false;
      }
    },
    togglePlayerSize() {
      if (this.isMedium) {
        this.isMedium = false
        this.videoElement.classList.remove('player-medium')
        this.videoElement.classList.add('player-big')
        this.changePlayerSizeButton.innerText = "Уменьшить"
        this.bufferedCanvas.width = '1327px'
      } else {
        this.isMedium = true
        this.videoElement.classList.add('player-medium')
        this.videoElement.classList.remove('player-big')
        this.changePlayerSizeButton.innerText = "Увеличить"
        this.bufferedCanvas.width = '992px'
      }
    },
    changeQuality() {
      console.log(this.hls)

      let loadLevelNumber = this.selectQualityElement.value;
      console.log(loadLevelNumber);
      this.hls.loadLevel = parseInt(loadLevelNumber)

      // this.hls.loadSource(process.env.VUE_APP_VIDEO_SERVER_ADDRESS + '/media/' + this.id + '/' + quality + '/stream/');
      // let tmpVideoElement = document.createElement("video");
      // tmpVideoElement.setAttribute('id', 'video');
      // let width;
      // let cssClass;
      //
      // if (this.isMedium) {
      //   width="720px"
      //   cssClass = "player-medium"
      // } else {
      //   width="1327px"
      //   cssClass = "player-big"
      // }
      //
      // tmpVideoElement.setAttribute('autoplay', 'autoplay');
      // tmpVideoElement.setAttribute('class', cssClass);
      // tmpVideoElement.setAttribute('width', width);
      //
      // this.videoWrapperElement.insertBefore(tmpVideoElement, this.videoControlElement)
      // // this.videoWrapperElement.appendChild(tmpVideoElement)
      // // this.hls.attachMedia(tmpVideoElement);
      // tmpVideoElement.setAttribute('id', 'video');
      // let currentTime = this.videoElement.currentTime;
      // this.videoElement.parentNode.removeChild(this.videoElement)
      // this.videoElement = tmpVideoElement;
      // this.videoElement.currentTime = currentTime
      // // this.prepareHls(this.id, value)
    },
    play() {
      this.videoElement.play()
      this.isPause = false
      this.playButton.innerText = "Pause"
    },
    pause() {
      this.videoElement.pause()
      this.isPause = true
      this.playButton.innerText = "Play"
    },
    startLoadHls() {
      this.hls.startLoad()
    },
    stopLoadHls() {
      this.hls.stopLoad()
    },
    setPlaybackSpeed() {
      this.playbackRate = document.getElementById("playSpeed");
      const playbackRate = this.playbackRate.value;
      if (playbackRate >= 0.25 && playbackRate <= 1.75) {
        this.videoElement.defaultPlaybackRate = playbackRate;
        this.videoElement.playbackRate = playbackRate;
      }
    },
    shiftPlaybackSpeed(shift) {
      if (this.videoElement.defaultPlaybackRate + shift >= 0.25 && this.videoElement.defaultPlaybackRate + shift <= 1.75) {
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

        console.log('this.video.defaultPlaybackRate')
        console.log(this.videoElement.defaultPlaybackRate)
        const defaultPlaybackRate = this.videoElement.playbackRate;
        for (let i = 0; i < children.length; i++) {
          const optionElement = children[i];
          const value = optionElement.getAttribute('value');
          if (value !== null) {
            const number = parseFloat(value);
            console.log(number)
            if (number === defaultPlaybackRate) {
              // console.log(value)
              // console.log('value')
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
    setVolumeText() {
      this.volume.innerText = Math.ceil(this.videoElement.volume * 100) + '%'
    },
    // updateVolume() {
    //   if (this.video !== null) {
    //     const videoVolumeOutput = document.getElementById('videoVolumeOutput');
    //     const newVolume = parseInt(videoVolumeOutput.innerText) / 100;
    //     console.log(newVolume)
    //     while (this.video.volume > 0 || this.video.volume < 1) {
    //       if ((newVolume > this.video.volume && newVolume < this.video.volume + 5) || (newVolume < this.video.volume && newVolume < this.video.volume - 5)) {
    //         break
    //       }
    //       if (this.video.volume > newVolume) {
    //         this.video.volume -= 0.05
    //       } else if (this.video.volume < newVolume) {
    //         this.video.volume += 0.05
    //       }
    //     }
    //
    //     this.setVolumeText()
    //   }
    // },
  },
}
</script>

<style scoped>
.player-big {
  width: 1327px;
}

.player-medium {
  width: 992px;
}

.hide {
  display: none;
}

.video-wrapper {
  position: relative;
  overflow: hidden;
}

.video-controls {
  position: absolute;
  top: 320px;
  z-index: 1;
  background-color: transparent;
}

#defaultBar {
  position: relative;
  float: left;
  width: 720px;
  height: 10px;
  padding: 4px;
  background: yellow;
}

#progressBar {
  position: absolute;
  width: 0;
  height: 5px;
  background: blue;
}
</style>