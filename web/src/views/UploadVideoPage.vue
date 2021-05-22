<template>
  <div>
    <div v-if="processing !== null">
      <span v-if="processing">
       <span v-if="videoId !== null">
          In Processing. Coming soon available at
         <router-link :to="{ name: 'videoStream', params: { videoId: videoId }}" class="subtitle-2"><span
             class="subtitle-2">video</span>
          </router-link></span>
      </span>
      <span v-else>
        Failed upload video
      </span>
    </div>
    <div v-else>
      <div>
        <label for="title">Name
          <input id="title" type="text" name="title" ref="title" v-on:change="onChangeTitle()"/>
        </label>
      </div>
      <div>
        <label for="description">Description
          <input id="description" type="text" name="description" ref="description"
                 v-on:change="onChangeDescription()"/>
        </label>
      </div>
      <div>
        <label for="file">Video file
          <input id="file" type="file" name="file" ref="file" v-on:change="onChangeFile()"/>
        </label>
      </div>
      <v-btn v-on:click="submitFile" type="submit" value="upload">Upload</v-btn>
    </div>
  </div>
</template>

<script>
import {mapActions, mapGetters} from "vuex";

export default {
  name: "UploadVideoPage",
  data() {
    return {
      file: '',
      title: '',
      description: '',
      processing: null,
      videoId: null
    }
  },
  methods: {
    ...mapActions({
      uploadVideo: "video/uploadVideo"
    }),
    ...mapGetters({
      getVideoId: "video/getVideoId",
      getIsProcessing: "video/getIsProcessing"
    }),
    onChangeFile() {
      this.file = this.$refs.file.files[0]
    },
    onChangeDescription() {
      this.description = this.$refs.description.value
    },
    onChangeTitle() {
      console.log(this.$refs.title)
      this.title = this.$refs.title.value
    },
    async submitFile() {
      await this.uploadVideo({file: this.file, title: this.title, description: this.description})
          .then(() => {
            console.log('uploaded')
            this.processing = this.getIsProcessing();
            this.videoId = this.getVideoId()
          });

    },
  }
}
</script>

<style scoped>

</style>