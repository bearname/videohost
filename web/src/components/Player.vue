<template>
  <div>
    <div v-if="qualities !== null">
      <div id="videoWrapper" class="player-wrapper player-medium">
        <video id="video" autoplay="autoplay" :poster="poster" class="player-video">
          <!--          <track label="English" kind="subtitles" srclang="en"-->
          <!--                 src="http://localhost:8000/content/afe67012-bc76-11eb-afc7-e4e74940035b/subtitle.srt" default>-->
        </video>
        <div id="videoControls" class="player-controls">
          <div id="progressBarWrapper" class="wrapper-bar" v-on:click="onClickBufferedRange($event);">
            <div id="progressBar" class="progress-bar"></div>
          </div>
          <div id="progressBarCircle" class="progress-bar-circle"></div>
          <div>
            <div id="thumbnail-wrapper" class="thumbnail-wrapper hide"><img id="thumbnail" src="" alt="preview"/></div>
            <span id="showHoverTime" class="hide"></span>
            <span class="player-controls-item">
              <button
                  id="playButton"
                  type="button"
                  title="video.play()"
                  v-on:click="togglePlayPause()">&#x23F5;</button>
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
            <div class="float-right right-control">
              <span id="settingButton" class="gear player-controls-item" title="Setting"></span>
              <span v-on:click="togglePlayerSize" id="changePlayerSize" class="player-controls-item"
                    title="Theater mode">&#9645;
              </span>
              <span v-on:click="toggleFullScreen" id="toggleFullScreen " class="float-right player-controls-item"
                    title="Full screen">&#x26F6;
              </span>
            </div>
          </div>
          <div id="settings" class="settings-popup hide">
            <div>
              <span id="playbackSpeedSetting">Playback speed</span>
              <select name="playSpeed" id="playSpeed" v-on:input="updatePlaybackSpeed()">
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
          <div id="contextMenu" class="settings-popup hide">
            <div>
              <button id="loopButton">loop</button>
            </div>
            <div>
              <button id="copyVideo">copy video url</button>
            </div>
            <div>
              <button id="copyVideoURLWithCurrentTime">Copy video URL at current time</button>
            </div>
            <div>
              <button id="statsOfNerds">Stats of nerds</button>
            </div>
            <div>
              <button id="keyboardHelp">Keyboard help</button>
            </div>
          </div>
        </div>
        <div id="submenu" class="stats-show hide"></div>
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

const USER_HELP = [
  {key: 'j', message: 'Seek backward 10 seconds in player.'},
  {key: 'l', message: 'Seek forward 10 seconds in player.'},
  {key: 'k', message: 'Pause/Play in player.'},
  {key: 'Spacebar', message: 'Pause/Play in player.'},
  {key: 'm', message: 'Mute/unmute the video.'},
  {key: '<', message: 'Slow down the video playback rate.'},
  {key: '>', message: 'Speed up the video playback rate.'},
  {key: 'Numbers 0 to 9 on the seek bar (not on the numeric pad)', message: 'Seek to the 0% to 90% of the video.'},
  {
    key: 'f',
    message: 'Activate full screen. If full screen mode is enabled, activate F again or press escape to exit full screen mode.'
  },
]
const PREVIEW_WIDTH = 256;
const PREVIEW_HEIGHT = 144;
const TILE_SIZE = {x: 5, y: 5}
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
    if (this.qualities.length === 0 && this.qualities[0] === "") {
      this.qualities = null
    }
    this.qualities.reverse()
    this.qualities = this.qualities.filter(quality => {
      return quality.length !== 0
    })
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
      previousPlaybackSpeed: 1,
      shiftTime: 10,
      volume: null,
      hls: null,
      isPause: false,
      videoElement: null,
      previewImageElement: null,
      previewImageWrapperElement: null,
      playbackSpeed: null,
      isFullScreen: false,
      bufferingIdx: null,
      lastStartPosition: 0,
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
      toggleFullScreenElement: null,
      settingsPopupElement: null,
      playbackSpeedSettingElement: null,
      qualitySettingElement: null,
      toggleLoopButtonElement: null,
      contextMenuElement: null,
      copyVideoElement: null,
      copyVideoURLWithCurrentTimeElement: null,
      statsOfNerdsElement: null,
      keyboardHelpElement: null,
      submenuShowElement: null,
      videos: null,
      volumePrevious: 1,
      showSetting: false,
      isLoopedVideo: false,
      prevSeekTime: 0,
      previousImageIndex: 0,
      previewImages: [],
    }
  },
  async mounted() {
    const thumbnails = [];
    for (let i = 0; i < this.duration / (TILE_SIZE.x * TILE_SIZE.y); i++) {
      thumbnails.push(process.env.VUE_APP_VIDEO_API + '/content/' + this.videoId + "/output-" + i + ".png");
    }
    await this.preloadPreviewImages(thumbnails);

    this.initPlayer()
    this.initKeyHandler()
    // this.fetchVideosByPage(0, 10)
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
    async preloadPreviewImages(array) {
      if (!this.previewImages) {
        this.previewImages = [];
      }

      for (let i = 0; i < array.length; i++) {
        let img = new Image();
        img.onload = function () {
          let index = this.previewImages.indexOf(i);
          if (index !== -1) {
            this.previewImages.splice(index, 1);
          }
        }
        this.previewImages.push(img);
        img.src = array[i];
      }
    },
    async fetchVideosByPage(page, countVideoOnPage) {
      await this.getVideoOnPage(page, countVideoOnPage);
      this.videos = this.getVideos();
      console.log(this.videos);
      this.countPage = this.getPageCount();
    },
    onClickBufferedRange(event) {
      this.videoElement.currentTime = this.getEventMouseTime(event);
    },
    initPlayer() {
      this.initElements()
      this.addToElementEventListeners()
      this.initHls();
    },
    initElements() {
      this.initVideoElement();
      this.initVideoControls();
      this.initContextMenu();
    },
    initVideoElement() {
      this.videoWrapperElement = document.getElementById('videoWrapper');
      this.videoElement = document.getElementById('video');
      this.previewImageElement = document.getElementById('thumbnail');
      this.previewImageWrapperElement = document.getElementById('thumbnail-wrapper');
      this.seekRequestTime()
    },
    seekRequestTime() {
      let time = this.$route.query.time;
      if (time !== undefined) {
        this.videoElement.currentTime = time;
      }
    },
    initVideoControls() {
      this.initVideoControlsItem()
      this.initProgressBar();
      this.initSettingElement();
    },
    initContextMenu() {
      this.settingsPopupElement = document.getElementById('settings');
      this.contextMenuElement = document.getElementById('contextMenu');
      this.toggleLoopButtonElement = document.getElementById('loopButton');
      this.copyVideoElement = document.getElementById('copyVideo');
      this.copyVideoURLWithCurrentTimeElement = document.getElementById('copyVideoURLWithCurrentTime');
      this.statsOfNerdsElement = document.getElementById('statsOfNerds');
      this.keyboardHelpElement = document.getElementById('keyboardHelp');
      this.submenuShowElement = document.getElementById('submenu');
    },
    initVideoControlsItem() {
      this.videoControlElement = document.getElementById('videoControls');
      this.currentTimeElement = document.getElementById('currentTime');
      this.volume = document.getElementById('volume');
      this.playButton = document.getElementById('playButton');
      this.pauseButton = document.getElementById('pauseButton');
      this.volumeMuteElement = document.getElementById('volumeMute');
      this.volumeChangerElement = document.getElementById('volumeChanger');
      this.changePlayerSizeButton = document.getElementById('changePlayerSize')
      this.toggleFullScreenElement = document.getElementById('toggleFullScreen');
    },
    initProgressBar() {
      this.progressBarWrapperElement = document.getElementById('progressBarWrapper');
      this.progressBarElement = document.getElementById('progressBar');
      this.showHoverTimeElement = document.getElementById('showHoverTime');
      this.progressBarCircleElement = document.getElementById('progressBarCircle');
    },
    initSettingElement() {
      this.settingButtonElement = document.getElementById('settingButton');
      this.selectQualityElement = document.getElementById('selectQuality');
      this.qualitySettingElement = document.getElementById('qualitySetting');
      this.playbackSpeedSettingElement = document.getElementById('playbackSpeedSetting');
    },
    addToElementEventListeners() {
      this.addEventListenerOnVideoElement();
      this.addEventListenerOnVideoControlItems();
      this.addEventListenerOnProgressBar();
      this.addEventListenerOnSettingMenu();
      this.addEventListenerOnContextMenu();
    },
    addEventListenerOnVideoElement() {
      this.videoElement.addEventListener('click', () => {
        this.hideElement(this.contextMenuElement);
        this.hideElement(this.submenuShowElement);
        this.togglePlayPause()
      })

      this.videoElement.addEventListener('dblclick', this.toggleFullScreen, false)
      this.videoElement.addEventListener('timeupdate', this.handleTimeUpdate, false)
      this.videoElement.addEventListener('stop', this.onVideoStop, false)
      // document.getElementById('video').addEventListener('loadeddata', this.generateThumbnails, false);
    },
    addEventListenerOnVideoControlItems() {
      this.volumeMuteElement.addEventListener('click', this.toggleMute, false)
      this.volumeChangerElement.addEventListener('change', this.handleVolumeChange, false)
      this.volumeChangerElement.addEventListener('mousemove', this.handleVolumeChange, false)
    },
    addEventListenerOnProgressBar() {
      this.progressBarWrapperElement.addEventListener('mousemove', this.onHoverProgressBar, false)
      this.progressBarWrapperElement.addEventListener('mouseout', this.onMouseOutProgressBar, false)
    },
    addEventListenerOnSettingMenu() {
      this.settingButtonElement.addEventListener('click', this.toggleSettingPopup, false)
      this.toggleLoopButtonElement.addEventListener('click', this.toggleLoop, false)
      this.copyVideoElement.addEventListener('click', this.copyVideoPath, false)
      this.copyVideoURLWithCurrentTimeElement.addEventListener('click', this.copyVideoPathWithTime, false)
    },
    addEventListenerOnContextMenu() {
      if (this.videoWrapperElement.addEventListener) {
        this.videoWrapperElement.addEventListener('contextmenu', this.toggleContextMenu, false);
      } else {
        this.videoWrapperElement.attachEvent('oncontextmenu', function () {
          window.event.returnValue = false;
        });
      }
      this.videoWrapperElement.addEventListener('click', this.togglePlayPause, false);
      this.statsOfNerdsElement.addEventListener('click', this.showStatsOfNerds, false)
      this.keyboardHelpElement.addEventListener('click', this.showKeyboardHelp, false)
    },
    initHls() {
      if (Hls.isSupported()) {
        this.prepareHls(this.id, this.selectedQuality);
      } else if (this.videoElement.canPlayType('application/vnd.apple.mpegurl')) {
        this.videoElement.src = process.env.VUE_APP_VIDEO_API + '/media/' + this.id + '/stream/';
        this.videoElement.addEventListener('loadedmetadata', function () {
          this.video.play();
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
      this.hls.loadSource(process.env.VUE_APP_VIDEO_API + '/media/' + id + '/stream/');
      this.hls.attachMedia(this.videoElement);
      this.hls.on(Hls.Events.MANIFEST_PARSED, () => {
        this.play();
      });
    },
    initKeyHandler() {
      window.addEventListener("keydown", this.handleKeyDownEvent, false);
    },
    togglePlayPause() {
      if (!this.videoElement.paused && !this.videoElement.ended) {
        this.pause();
      } else if (this.videoElement.ended) {
        this.videoReplay()
      } else {
        this.play();
      }
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
        this.hideElement(this.changePlayerSizeButton);
        this.isFullScreen = true;
        this.toggleFullScreenElement.setAttribute('title', "Full screen");
      } else {
        this.isFullScreen = false;
        this.unHideElement(this.changePlayerSizeButton);
        this.toggleFullScreenElement.setAttribute('title', "Exit full screen")
      }
    },
    togglePlayerSize() {
      if (this.isMedium) {
        this.isMedium = false;
        this.removeCss(this.videoWrapperElement, 'player-medium');
        this.addCss(this.videoWrapperElement, 'player-big');
        this.changePlayerSizeButton.setAttribute('title', 'Theater mode');
      } else {
        this.isMedium = true;
        this.addCss(this.videoWrapperElement, 'player-medium');
        this.removeCss(this.videoWrapperElement, 'player-big');
        this.changePlayerSizeButton.setAttribute('title', 'Default view');
      }
    },
    changeQuality() {
      console.log(this.hls)

      let loadLevelNumber = this.selectQualityElement.value;
      console.log(loadLevelNumber);
      this.hls.loadLevel = parseInt(loadLevelNumber)
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
    copyTextToClipboard(text) {
      if (!navigator.clipboard) {
        console.log("failed copy")
        return;
      }
      navigator.clipboard.writeText(text).then(function () {
        console.log('Async: Copying to clipboard was successful!');
      }, function (err) {
        console.error('Async: Could not copy text: ', err);
      });
    },
    updatePlaybackSpeed() {
      this.playbackRate = document.getElementById("playSpeed");
      const playbackRate = this.playbackRate.value;
      this.setPlaybackSpeed(playbackRate);
    },
    toggleSlowMotion() {
      if (this.playbackRate === 0.25) {
        this.playbackRate = this.previousPlaybackSpeed
      } else {
        this.previousPlaybackSpeed = this.playbackRate;
        this.playbackRate = 0.25
      }
    },
    setPlaybackSpeed(playbackRate) {
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
        children.forEach((part, index) => {
          const optionElement = this[index];
          const attribute = optionElement.getAttribute('selected');
          if (attribute !== null) {
            optionElement.removeAttribute('selected');
          }
        }, children);

        const defaultPlaybackRate = this.videoElement.playbackRate;
        for (let i = 0; i < children.length; i++) {
          const optionElement = children[i];
          const value = optionElement.getAttribute('value');
          if (value !== null) {
            const number = parseFloat(value);
            console.log(number);
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
      this.volumeChangerElement.value = this.videoElement.volume
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
        this.videoElement.volume = 1;
        this.videoElement.muted = false;
        this.changeButtonType(this.volumeMuteElement, 'mute');
        this.volumeChangerElement.value = this.videoElement.volume;
      } else {
        this.volumePrevious = this.videoElement.volume;
        this.videoElement.muted = true;
        this.changeButtonType(this.volumeMuteElement, 'unmute');
        this.volumeChangerElement.value = 0;
      }

      this.handleVolumeChange();
    },
    changeButtonType(btn, value) {
      btn.title = value;
      btn.innerHTML = value;
      btn.className = value;
    },
    onHoverProgressBar(event) {
      const boundingClientRect = this.progressBarWrapperElement.getBoundingClientRect();

      const seekTime = this.getEventMouseTime(event);
      if (this.prevSeekTime === seekTime) {
        return;
      }
      this.setNewPreviewImageIfNeeded(seekTime)

      this.prevSeekTime = seekTime;
      const offsetX = event.x - boundingClientRect.x;

      const point = this.getPreviewImageOffset(seekTime);
      this.changePreviewImage(point.x, point.y, offsetX, boundingClientRect)

      this.changeHoverTimeElement(seekTime, offsetX)
    },
    getEventMouseTime(event) {
      const boundingClientRect = this.progressBarWrapperElement.getBoundingClientRect();

      return parseInt((event.pageX - boundingClientRect.left) * this.duration / boundingClientRect.width);
    },
    onMouseOutProgressBar() {
      this.hideElement(this.previewImageWrapperElement);
      this.hideElement(this.showHoverTimeElement)
    },
    onVideoStop() {
      this.videoControlElement.style.transform = 'translateY(-20px)';
    },
    handleVolumeChange() {
      this.volumePrevious = this.volumeChangerElement.value;
      this.videoElement.volume = this.volumeChangerElement.value;
      if (this.videoElement.volume === 0) {
        this.volumeMuteElement.innerText = "unmute";
        this.volumeChangerElement.value = 0;
        this.videoElement.muted = true;
      } else {
        this.videoElement.muted = false;
        this.volumeMuteElement.innerText = "mute";
        this.volumeChangerElement.value = this.volumePrevious;
      }
      this.updateVolumeText();
    },
    toggleSettingPopup() {
      let boundingClientRect = this.settingButtonElement.getBoundingClientRect();
      this.setElemPosition(this.settingsPopupElement, {left:  boundingClientRect.left + 100, top: boundingClientRect.top - 300})
      this.toggleHideElement(this.settingsPopupElement);
    },
    toggleLoop() {
      if (!this.isLoopedVideo) {
        this.isLoopedVideo = true;
        this.toggleLoopButtonElement.textContent = "looped";
        this.videoElement.addEventListener('ended', this.loopVideo, false);
      } else {
        this.toggleLoopButtonElement.textContent = "loop";
        this.isLoopedVideo = false;
      }
    },
    loopVideo() {
      if (this.isLoopedVideo) {
        this.videoReplay();
      } else {
        this.playButton.textContent = 'Replay';
      }
    },
    toggleContextMenu(e) {
      e.preventDefault();
      this.toggleHideElement(this.contextMenuElement);
    },
    videoReplay() {
      this.currentTime = 0;
      this.play();
    },
    copyVideoPath() {
      this.copyTextToClipboard(document.URL);
    },
    copyVideoPathWithTime() {
      const string = document.URL.split('?')[0];
      this.copyTextToClipboard(string + "?time=" + this.videoElement.currentTime);
    },
    showStatsOfNerds() {
      const innerHtml = this.getDebugInfo();

      this.toggleShowPopup(this.submenuShowElement, innerHtml);
    },
    showKeyboardHelp() {
      let innerHtml = '';
      USER_HELP.forEach(element => {
        innerHtml += '<p>' + element.key + ' - ' + element.message + '</p>';
      });
      this.toggleShowPopup(this.submenuShowElement, innerHtml);
    },
    toggleShowPopup(element, innerHtml) {
      this.toggleHideElement(element);
      this.submenuShowElement.innerHTML = innerHtml;
    },
    volumeUp(event) {
      if (this.videoElement.volume + 0.05 <= 1) {
        event.preventDefault()
        this.videoElement.volume += 0.05;
        this.updateVolumeText();
      }
    },
    volumeDown(event) {
      if (this.videoElement.volume >= 0) {
        event.preventDefault();
        this.videoElement.volume -= 0.05;
        this.updateVolumeText();
      }
    },
    handleKeyDownEvent(e) {
      const key = e.key;
      console.log('key');
      console.log(key);
      switch (key) {
        case 'j': {
          this.shiftCurrentTime(-this.shiftTime);
          break;
        }
        case 's': {
          this.toggleSlowMotion();
          break;
        }
        case 'f': {
          this.toggleFullScreen();
          break;
        }
        case 't': {
          this.togglePlayerSize();
          break;
        }
        case '<': {
          this.shiftPlaybackSpeed(-0.25);
          break;
        }
        case '>': {
          this.shiftPlaybackSpeed(0.25);
          break;
        }
        case  'l': {
          this.shiftCurrentTime(this.shiftTime);
          break;
        }
        case 'k': {
          this.togglePlayPause();
          break;
        }
        case 'ArrowLeft': {
          this.shiftCurrentTime(-5);
          break;
        }
        case 'ArrowRight': {
          this.shiftCurrentTime(5);
          break;
        }
        case 'ArrowUp': {
          this.volumeUp(e);
          break;
        }
        case 'ArrowDown': {
          this.volumeDown(e);
          break;
        }
        case 'm': {
          this.toggleMute();
          break;
        }
        default: {
          const eventCode = e.code;

          console.log(eventCode)
          if (eventCode === 'Space') {
            e.preventDefault();
            this.togglePlayPause();
          } else {
            const s = eventCode.substring(0, eventCode.length - 1);
            const keyElement = eventCode[eventCode.length - 1];
            console.log(s);
            console.log(keyElement);

            if ((s === 'Numpad' || s === 'Digit') && keyElement >= '0' && keyElement <= '9') {
              const percent = parseInt(keyElement);
              console.log(percent / 10);
              const shift = this.duration * (percent / 10);
              console.log('shift');
              console.log(shift);
              this.setCurrentTime(shift);
            }
          }
          break;
        }
      }
    },
    getDebugInfo() {
      const hls = this.hls;
      const boundingClientRect = this.videoElement.getBoundingClientRect();

      return '<p>Video ID' + this.videoId + '</p>' +
          '<p>Viewport ' + boundingClientRect.width + 'x' + boundingClientRect.height + '</p>' +
          '<p>Volume ' + this.videoElement.volume * 100 + '%</p>' +
          '<p>Codecs ' + hls.levels[hls.currentLevel].videoCodec + '</p>' +
          '<p>Latency / Bandwidth Estimated ' + hls.latency.toFixed(3) + '  ' + hls.bandwidthEstimate + '</p>' +
          '<p>Max Latency: ' + this.hls.maxLatency + '</p>' +
          '<p>Target Latency: ' + this.hls.targetLatency.toFixed(3) + '</p>' +
          '<p>Bitrate: ' + Math.round(hls.levels[hls.currentLevel].bitrate / 1000) + '</p>' +
          '<p>Latency: ' + this.hls.latency.toFixed(3) + '</p>' +
          '<p>Drift: ' + this.hls.drift.toFixed(3) + '  (edge advance rate)</p>';
    },
    setNewPreviewImageIfNeeded(seekTime) {
      const currentImageIndex = parseInt(seekTime / (TILE_SIZE.x + TILE_SIZE.y));
      if (this.previousImageIndex !== currentImageIndex) {
        this.previousImageIndex = currentImageIndex;
        this.previewImageElement.src = this.previewImages[this.previousImageIndex].src;
      }
    },
    changePreviewImage(thumbX, thumbY, offsetX, boundingRect) {
      this.previewImageElement.style.transform = "translate(-" + PREVIEW_WIDTH * thumbX + "px, -" + PREVIEW_HEIGHT * thumbY + "px)";
      const offset = this.calculatePreviewImageOffset(offsetX, boundingRect)

      this.setElemPosition(this.previewImageWrapperElement, {left:   offset.left, top:offset.top })

      this.unHideElement(this.previewImageWrapperElement);
    },
    getPreviewImageOffset(seekTime) {
      const seconds = seekTime % (TILE_SIZE.x + TILE_SIZE.y);
      const thumbY = parseInt(seconds / TILE_SIZE.y);
      const thumbX = parseInt(seconds % TILE_SIZE.x);
      return {x: thumbX, y: thumbY}
    },
    changeHoverTimeElement(seekTime, x) {
      this.showHoverTimeElement.style.left = x - 15 + 'px';
      this.showHoverTimeElement.style.top = -20 + 'px';
      this.setElemPosition(this.showHoverTimeElement, {left:  x - 15 + 'px', top:  -20  })

      this.unHideElement(this.showHoverTimeElement);
      this.showHoverTimeElement.innerText = this.formatTimeString(seekTime);
    },

    calculatePreviewImageOffset(offsetX, boundingRect) {
      let x = offsetX;
      const centerX = PREVIEW_WIDTH/ 2;
      const leftX = boundingRect.x + centerX - 50;
      const rightX = boundingRect.x + boundingRect.width - PREVIEW_WIDTH - 50;

      if (x > leftX && x < rightX) {
        x -= centerX;
      } else if (x <= leftX) {
        x += 0;
      } else if (x >= rightX) {
        x -= PREVIEW_WIDTH;
      }

      return {left: x, top: -PREVIEW_HEIGHT - 20}
    },
    hideElement(element) {
      this.addCss(element, 'hide')
    },
    unHideElement(element) {
      this.removeCss(element, 'hide')
    },
    toggleHideElement(element) {
      element.classList.toggle('hide')
    },
    addCss(element, cssClassName) {
      element.classList.add(cssClassName);
    },
    removeCss(element, cssClassName) {
      element.classList.remove(cssClassName);
    },
    setElemPosition(element, {left, top}) {
      element.style.top = top + "px";
      element.style.left = left + "px";
    }
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
  max-width: 95%;
}

.player-medium {
  max-width: 992px;
}

.player-wrapper:fullscreen .player-video {
  max-width: none;
  width: 100%;
  margin: auto 0;
}

.player-wrapper:-webkit-full-screen .player-video {
  max-width: none;
  width: 100%;
}

.player-wrapper:fullscreen .player-controls {
  max-width: none;
  width: 100%;
  margin: auto 0;
}

.player-wrapper:-webkit-full-screen .player-controls {
  max-width: none;
  width: 100%;
}

.hide {
  opacity: 0;
  display: none;
}

.player-wrapper {
  position: relative;
  overflow: hidden;
}

.player-controls {
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
  vertical-align: middle;
  cursor: pointer;
}

.right-control {
  margin: 0 5px;
  line-height: 26px;
  font-size: 20px;
}

.player-wrapper:hover .player-controls {
  transform: translateY(0px);
  opacity: 1;
}

.player-wrapper:hover .progress-bar {
  height: 4px;
}

.wrapper-bar {
  position: relative;
  float: left;
  width: 100%;
  height: 8px;
  background: rgba(255, 255, 255, 0.4);
  margin-bottom: 14px;
}

.progress-bar,
.wrapper-bar {
  cursor: pointer;
  position: relative;
  box-sizing: inherit;
}

.progress-bar {
  background-color: #f00;
  height: 3px;
}

.wrapper-bar:hover .progress-bar {
  height: 5px;
}

.wrapper-bar:hover, .progress-bar-circle {
  opacity: 1;
}

.progress-bar-circle {
  position: absolute;
  top: -5px;
  left: 0;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background-color: #f00;
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

.stats-show {
  position: absolute;
  top: 10px;
  left: 10px;
  padding: 10px 10px;
  background: rgba(28, 28, 28, 0.9);
  border-radius: 4px;
  color: #e8dbdb;
}

.thumbnail-wrapper {
  position: absolute;
  overflow: hidden;
  width: 256px;
  height: 144px;
  border: 2px white solid;
  border-radius: 3%;
}

.thumbnail-wrapper img {
  width: 1280px;
  height: 720px;
  position: relative;
  /*transform: translate(0, 0);*/
}
</style>