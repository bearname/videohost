<template>
  <div>
    <div v-if="qualities !== null">
      <div id="video-wrapper" class="video-wrapper">
        <video id="video" width="720px" autoplay="autoplay" :poster="poster"></video>
        <div>
          <span>Quality: </span>
          <select v-model="selectedQuality">
            <option disabled value="">Please select one</option>
            <option v-for="quality in qualities" :key="quality" :value="quality">{{ quality }}p</option>
          </select>
          <button v-on:click="changePlayerSize"><span v-if="isMedium">Увеличить</span><span v-else>Уменьшить</span>
          </button>
        </div>
        <div id="videoControls" class="hide video-controls">
          <canvas
              id="bufferedCanvas"
              width="720"
              height="15"
              class="videoCentered"
              v-on:click="onClickBufferedRange($event);"
              style="height: fit-content"
          ></canvas>
          <div>
            <button
                type="button"
                class="btn btn-sm btn-info"
                title="video.play()"
                v-on:click="play()"
            >
              Play
            </button>
          </div>
          <div>
            <button
                type="button"
                class="btn btn-sm btn-info"
                title="video.pause()"
                v-on:click="pause()"
            >
              Pause
            </button>
          </div>
          <div>
            <button type="button"
                    class="btn btn-sm btn-info"
                    title="video.playbackRate = text input"
                    v-on:click="setPlaybackSpeed()"
            >
              Playback speed
            </button>
            <select name="playSpeed" id="playSpeed">
              <option disabled value="">Playback Speed</option>
              <option value="0.25">0.25x</option>
              <option value="0.50">0.5x</option>
              <option value="0.75">0.75x</option>
              <option value="1" selected="selected">Normal</option>
              <option value="1.25">1.25x</option>
              <option value="1.50">1.50x</option>
              <option value="1.75">1.75x</option>
            </select>
          </div>
          <div>
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
          </div>
          <div>
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
          </div>
          <div>Volume
            <span id="volume" v-if="this.video !== null">{{ this.video.volume * 100 }}</span>
            <span v-else>100%</span>
          </div>
          <span v-on:click="toggleFullScreen">full screen</span>
        </div>
      </div>

    </div>
    <div v-else>
      Video not available
    </div>

  </div>
</template>

<script>
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
    console.log("this.qualities");
  },
  data() {
    return {
      id: this.videoId,
      selectedQuality: '1080',
      qualities: null,
      videoDuration: this.duration,
      isMedium: true,
      video: null,
      previousVolume: 1,
      playbackRate: null,
      shiftTime: 10,
      volume: null,
      hls: null,
      isPause: true,
      poster: this.thumbnail,
      playbackSpeed: null,
      isFullScreen: false,
      bufferingIdx: null,
      lastStartPosition: 0,
      canvas: null,
      videoControlElement: null
    }
  },
  mounted() {
    this.initPlayer()
    this.initKeyHandler()
  },
  updated() {
    this.initPlayer()
  },
  methods: {
    onClickBufferedRange(event) {
      if (this.canvas === null) {
        this.canvas = document.querySelector('#bufferedCanvas');
      }
      this.video.currentTime = ((event.clientX - this.canvas.offsetLeft) / this.canvas.width) * this.getSeekableEnd();
    },
    getSeekableEnd() {
      if (isFinite(this.video.duration)) {
        return this.video.duration;
      }
      if (this.video.seekable.length) {
        return this.video.seekable.end(this.video.seekable.length - 1);
      }
      return 0;
    },
    checkBuffer() {
      // const colors = ['#000000', '#d21c1c', '#52fd16', '#0409fd',]
      // const rndInt = Math.floor(Math.random() * colors.length)
      // ctx.fillStyle = colors[rndInt]
      // ctx.fillRect(0, 0, this.canvas.clientWidth, this.canvas.clientHeight)
      const buffered = this.video.buffered;
      const seekableEnd = this.getSeekableEnd();
      console.log('buffered')
      console.log(buffered)
      console.log('seekableEnd')
      console.log(seekableEnd)
      // let bufferingDuration;
      const ctx = this.canvas.getContext('2d');
      if (buffered) {
        ctx.fillStyle = '#000000';
        if (!this.canvas.width || this.canvas.width !== this.video.clientWidth) {
          this.canvas.width = this.video.clientWidth;
        }
        ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
        const pos = this.video.currentTime;
        let bufferLen = 0;
        ctx.fillStyle = '#03bb2b';
        for (let i = 0; i < buffered.length; i++) {
          const start = (buffered.start(i) / seekableEnd) * this.canvas.width;
          const end = (buffered.end(i) / seekableEnd) * this.canvas.width;
          ctx.fillRect(start, 2, Math.max(2, end - start), 11);
          if (pos >= buffered.start(i) && pos < buffered.end(i)) {
            // play position is inside this buffer TimeRange, retrieve end of buffer position and buffer length
            bufferLen = buffered.end(i) - pos;
          }
        }
        console.log(bufferLen)
        // check if we are in buffering / or playback ended state
        // if (
        //     bufferLen <= 0.1 &&
        //     this.video.paused === false && pos - this.lastStartPosition > 0.5
        // ) {
        //   if (lastDuration - pos <= 0.5 && events.isLive === false) {
        //     // don't create buffering event if we are at the end of the playlist, don't report ended for live playlist
        //   } else {
        //     // we are not at the end of the playlist ... real buffering
        //     if (bufferingIdx !== -1) {
        //       bufferingDuration =
        //           self.performance.now() -
        //           events.t0 -
        //           events.video[bufferingIdx].time;
        //       this.video[bufferingIdx].duration = bufferingDuration;
        //       this.video[bufferingIdx].name = bufferingDuration;
        //     } else {
        //       this.video.push({
        //         type: 'buffering',
        //         // time: self.performance.now() - events.t0,
        //       });
        //       // trimEventHistory();
        //       // we are in buffering state
        //       // bufferingIdx = events.video.length - 1;
        //     }
        //   }
        // }
        //
        // if (bufferLen > 0.1 && bufferingIdx !== -1) {
        //   bufferingDuration =
        //       self.performance.now() - events.t0 - events.video[bufferingIdx].time;
        //   events.video[bufferingIdx].duration = bufferingDuration;
        //   events.video[bufferingIdx].name = bufferingDuration;
        //   // we are out of buffering state
        //   bufferingIdx = -1;
        // }

        // update buffer/position for current Time
        // const event = {
        //   time: self.performance.now() - events.t0,
        //   buffer: Math.round(bufferLen * 1000),
        //   pos: Math.round(pos * 1000),
        // };
        // const bufEvents = events.buffer;
        // const bufEventLen = bufEvents.length;
        // if (bufEventLen > 1) {
        //   const event0 = bufEvents[bufEventLen - 2];
        //   const event1 = bufEvents[bufEventLen - 1];
        //   const slopeBuf0 =
        //       (event0.buffer - event1.buffer) / (event0.time - event1.time);
        //   const slopeBuf1 =
        //       (event1.buffer - event.buffer) / (event1.time - event.time);
        //
        //   const slopePos0 = (event0.pos - event1.pos) / (event0.time - event1.time);
        //   const slopePos1 = (event1.pos - event.pos) / (event1.time - event.time);
        //   // compute slopes. if less than 30% difference, remove event1
        //   if (
        //       (slopeBuf0 === slopeBuf1 ||
        //           Math.abs(slopeBuf0 / slopeBuf1 - 1) <= 0.3) &&
        //       (slopePos0 === slopePos1 || Math.abs(slopePos0 / slopePos1 - 1) <= 0.3)
        //   ) {
        //     bufEvents.pop();
        //   }
        // }
        // events.buffer.push(event);
        // trimEventHistory();
        // canvas.refreshCanvas();

        ctx.fillStyle = '#0511b6';
        const x = (this.video.currentTime / seekableEnd) * this.canvas.width;
        ctx.fillRect(x, 0, 2, 15);
      } else if (ctx.fillStyle !== '#000') {
        ctx.fillStyle = '#000';
        ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
      }
    },
    initPlayer() {
      this.video = document.getElementById('video');
      this.volume = document.getElementById('volume');
      this.canvas = document.getElementById('bufferedCanvas')
      this.videoControlElement = document.getElementById('videoControls')
      this.video.addEventListener('mouseover', () => {
        this.videoControlElement.classList.remove('hide');
      })
      this.video.addEventListener('mouseout', () => {
        this.videoControlElement.classList.add('hide');
      })
      console.log(this.video.getVideoPlaybackQuality())
      if (Hls.isSupported()) {
        this.hls = new Hls();
        this.hls.loadSource(process.env.VUE_APP_VIDEO_SERVER_ADDRESS + '/media/' + this.id + '/' + this.selectedQuality + '/stream/');
        this.hls.attachMedia(this.video);
        this.hls.on(Hls.Events.MANIFEST_PARSED, () => {
          this.video.play();
          console.log('getVideoPlaybackQuality()')
        });
        console.log('frag_buffered')
        this.hls.on(Hls.Events.FRAG_BUFFERED, () => {
          console.log('frag,et buffered')
          // this.hls.bufferTimer = setInterval(this.checkBuffer, 300);
        })
      } else if (this.video.canPlayType('application/vnd.apple.mpegurl')) {
        this.video.src = process.env.VUE_APP_VIDEO_SERVER_ADDRESS + '/media/' + this.id + '/' + this.selectedQuality + '/stream/';
        this.video.addEventListener('loadedmetadata', function () {
          this.video.play();
          console.log('getVideoPlaybackQuality()')
        });
      } else {
        let videoWrapper = document.getElementById('video-wrapper');
        videoWrapper.innerText = "Not support streaming video"
      }
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
            // if (!this.isFullScreen) {
            //   this.video.requestFullscreen();
            //   this.isFullScreen = true;
            // } else {
            //   this.video.exitFullscreen();
            //   this.isFullScreen = false;
            // }
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
            this.togglePause()
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
            if (this.video.volume + 0.05 <= 1) {
              this.video.volume += 0.05
              this.setVolumeText()
            } else if (this.video.volume + 0.05 > 1) {
              this.video.volume += 1 - this.video.volume
            }
            break
          }
          case 'ArrowDown': {
            if (this.video.volume > 0) {
              this.video.volume -= 0.05
              this.setVolumeText()
            }
            break
          }
          case 'm' || 'M': {
            if (this.video.volume > 0) {
              this.previousVolume = this.video.volume
              this.video.volume = 0;
            } else if (this.video.volume === 0) {
              this.video.volume = this.previousVolume
              this.previousVolume = 0
            }
            break;
          }
          default: {
            const eventCode = e.code;

            console.log(eventCode)
            if (eventCode ==='Space') {
              e.preventDefault()
              this.togglePause()
            }
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
            break
          }
        }
      }, false);
    },
    togglePause() {
      if (this.isPause) {
        this.play()
      } else {
        this.pause()
      }
    },
    toggleFullScreen() {
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
    changePlayerSize() {
      let video = document.getElementById('video');
      if (this.isMedium) {
        this.isMedium = false
        video.setAttribute('width', '1420px');
        this.shiftCurrentTime(0)
      } else {
        this.isMedium = true
        video.setAttribute('width', '720px')
        this.shiftCurrentTime(0)
      }
    },
    play() {
      this.video.play()
      this.isPause = false
    },
    pause() {
      this.video.pause()
      this.isPause = true
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
        this.video.defaultPlaybackRate = playbackRate;
        this.video.playbackRate = playbackRate;
      }
    },
    shiftPlaybackSpeed(shift) {
      if (this.video.defaultPlaybackRate + shift >= 0.25 && this.video.defaultPlaybackRate + shift <= 1.75) {
        this.video.defaultPlaybackRate += shift;
        this.video.playbackRate += shift;

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
        console.log(this.video.defaultPlaybackRate)
        const defaultPlaybackRate = this.video.playbackRate;
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
      this.video.currentTime += shift
    },
    setCurrentTime(time) {
      this.video.currentTime = time
    },
    setVolumeText() {
      this.volume.innerText = Math.ceil(this.video.volume * 100) + '%'
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

.hide {
  display: none;
}

.player-medium {
  width: 992px;
}

.video-wrapper {
  position: relative;
  overflow: hidden;
}

.video-controls {
  position: absolute;
  top: 360px;
  z-index: 1;
  background-color: transparent;
}
</style>