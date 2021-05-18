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
      <label for="title">Name
        <input id="title" type="text" name="title" ref="title" v-on:change="onChangeTitle()"/>
      </label>
      <label for="description">Description
        <input id="description" type="text" name="description" ref="description"
               v-on:change="onChangeDescription()"/>
      </label>
      <label for="file">Video file
        <input id="file" type="file" name="file" ref="file" v-on:change="onChangeFile()"/>
      </label>
      <button v-on:click="submitFile()" type="submit" value="upload">Upload</button>
    </div>
  </div>
</template>

<script>
import axios from "axios";

// const apiUrl = "http://localhost:8000/api/v1";
// let apiUploadVideo = apiUrl + "/video/";
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
    submitFile() {
      const formData = new FormData();
      formData.append("file", this.file)
      formData.append("title", this.title)
      formData.append("description", this.description)
      const config = {
        headers: {
          'Content-Type': 'video/mp4'
        }
      }

      const onSuccess = response => {
        console.log(response);
        const status = response.status
        if (status !== 200) {
          this.processing = false
        } else {
          console.log('SUCCESS!!')
          this.processing = true
          this.videoId = response.data
        }
      };
      const onFail = error => {
        this.processing = false
        console.error(error)
        console.log('FAILURE!!')
      };
      axios.post("http://localhost:8000/api/v1/video", formData, config).then(onSuccess).catch(onFail)
    },
  }
}
</script>

<style scoped>

</style>