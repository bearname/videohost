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
      <v-btn v-on:click="submitFile()" type="submit" value="upload">Upload</v-btn>
    </div>
  </div>
</template>

<script>
import {mapActions} from "vuex";
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
    ...mapActions({
      uploadVideo: "video/uploadVideo"
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
    submitFile() {
      console.log("submit file")

      this.uploadVideo({file: this.file, title: this.title, description: this.description})
      // const formData = new FormData();
      // formData.append("file", this.file)
      // formData.append("title", this.title)
      // formData.append("description", this.description)
      // const config = {
      //   headers: {
      //     'Content-Type': 'application/x-www-form-urlencoded',
      //     // 'Content-Type': 'video/mp4',
      //     Authorization: context.rootGetters["auth/getTokenHeader"]
      //     // 'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjEzMjg4MzgsInVzZXJuYW1lIjoibWlraGEifQ.xbHVE_z8l_MLFXogE9jkIfEnFOd_FmjtFoYfu3r-xss'
      //   }
      // }
      // const onSuccess = response => {
      //   console.log(response);
      //   const status = response.status
      //   if (status !== 200) {
      //     this.processing = false
      //   } else {
      //     console.log('SUCCESS!!')
      //     this.processing = true
      //     this.videoId = response.data
      //   }
      // };
      // const onFail = error => {
      //   this.processing = false
      //   console.error(error)
      //   console.log('FAILURE!!')
      // };
      //
      // return axios.post("http://localhost:8000/api/v1/video", formData, config).then(onSuccess).catch(onFail)

      // this.uploadVideo({file: this.file, title: this.title, description: this.description})
    },
  }
}
</script>

<style scoped>

</style>