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
      <!--      <v-alert v-if="message" border="left" color="blue-grey" dark>-->
      <!--        {{ message }}-->
      <!--      </v-alert>-->
      <v-form>

        <v-container>
          <v-row>
            <v-col
                cols="12"
                md="4"
            >
              <v-file-input
                  show-size
                  label="File input"
                  v-on:change="onChangeFile"
              ></v-file-input>
            </v-col>
            <v-col
                cols="12"
                md="4"
            >
              <v-text-field
                  :rules="nameRules"
                  :counter="255"
                  label="Title"
                  required
                  @change="onChangeTitle"
              ></v-text-field>
            </v-col>

            <v-col
                cols="12"
                md="4"
            >
              <v-text-field
                  label="Description"
                  required
                  v-on:change="onChangeDescription"
              ></v-text-field>
            </v-col>
            <v-col>
              <v-btn v-on:click="submitFile" type="submit" value="upload">Upload</v-btn>
            </v-col>
          </v-row>
        </v-container>
      </v-form>

      <!--      <div>-->
      <!--        <label for="title">Name-->
      <!--          <input id="title" type="text" name="title" ref="title" v-on:change="onChangeTitle()"/>-->
      <!--        </label>-->
      <!--      </div>-->
      <!--      <div>-->
      <!--        <label for="description">Description-->
      <!--          <input id="description" type="text" name="description" ref="description"-->
      <!--                 v-on:change="onChangeDescription()"/>-->
      <!--        </label>-->
      <!--      </div>-->
      <!--      <div>-->
      <!--        <label>Video file&ndash;&gt;-->
      <!--          <v-file-input-->
      <!--              counter-->
      <!--              multiple-->
      <!--              small-chips-->
      <!--              truncate-length="50"-->
      <!--              name="file"-->
      <!--              ref="file"-->
      <!--              v-on:change="onChangeFile()"-->
      <!--          ></v-file-input>-->
      <!--        </label>-->
      <!--        <label for="file">Video file-->
      <!--          <input id="file" type="file" name="file" ref="file" v-on:change="onChangeFile()"/>-->
      <!--        </label>-->
      <!--      </div>-->
    </div>
  </div>
</template>

<script>
import {mapActions, mapGetters} from "vuex";

export default {
  name: "UploadVideoPage",
  data() {
    return {
      file: undefined,
      title: '',
      description: '',
      processing: null,
      videoId: null,
      valid: false,
      firstname: '',
      lastname: '',
      nameRules: [
        v => !!v || 'Name is required',
        v => v.length <= 255 || 'Name must be less than 10 characters',
      ],
      email: '',
      emailRules: [
        v => !!v || 'E-mail is required',
        v => /.+@.+/.test(v) || 'E-mail must be valid',
      ],
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
    onChangeFile(file) {
      this.file = file;
    },
    onChangeDescription(description) {
      this.description = description;
    },
    onChangeTitle(title) {
      this.title = title;
    },
    async submitFile(e) {
      e.preventDefault();
      // if (!this.file) {
      //   this.message = "Please select a file!";
      //   return;
      // }
      // this.message = "";

      await this.uploadVideo({file: this.file, title: this.title, description: this.description});
      console.log('uploaded');
      this.processing = this.getIsProcessing();
      this.videoId = this.getVideoId();
    },
  }
}

</script>

<style scoped>

</style>