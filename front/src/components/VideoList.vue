<template>
  <div
      class="mx-auto"
      max-width="1000"
  >
    <v-container fluid>
      <v-row dense>
        <v-col
            v-for="video in videos" :key="video.id"
            :cols="3"
        >
          <div style="overflow-y: auto; height:340px">
            <v-img
                :src="video.thumbnail"
                class="white--text align-end"
                gradient="to bottom, rgba(0,0,0,.1), rgba(0,0,0,.5)"
                height="200px"
            ></v-img>
            <v-card-title>
              <router-link :to="{ name: 'videoStream', params: { videoId: video.id }}" class="subtitle-2"><span
                  class="subtitle-2">{{
                  video.name
                }}</span>
              </router-link>
              {{ video.duration }} s.
            </v-card-title>
          </div>
        </v-col>
      </v-row>
    </v-container>
    <v-btn class="v-btn"
           :disabled="pageNumber === 0"
           @click="previousPage">
      Previous
    </v-btn>
    <v-btn class="v-btn"
           :disabled="(pageNumber >= countPage && countPage < 2) || (pageNumber >= countPage - 1 && countPage >= 2) "
           @click="nextPage">
      Next
    </v-btn>
  </div>

</template>

<script>
import axios from 'axios'

export default {
  name: "VideoList",
  props: {
    countVideoOnPage: {
      type: Number,
      required: false,
      default: 12
    },
    countPage: {
      type: Number,
      required: false,
      default: 0
    },
  },
  data() {
    return {
      error: false,
      pageNumber: 1,
      videos: null
    }
  },
  methods: {
    previousPage() {
      if (this.pageNumber > 0) {
        this.pageNumber--;
        this.fetchVideos(this.pageNumber, this.countVideoOnPage)
      }
    },
    nextPage() {
      if (this.pageNumber < this.countPage) {
        this.pageNumber++;
        this.fetchVideos(this.pageNumber, this.countVideoOnPage)
      }
    },
    fetchVideos(page = '1', countVideoOnPage = '10') {
      let url = 'http://localhost:8000/api/v1/list?page=' + page + '&countVideoOnPage=' + countVideoOnPage;
      console.log(url)
      axios.get(url)
          .then(response => {
            console.log(response.data)
            if (Object.keys(response.data).includes("pageCount")) {
              this.countPage = response.data.pageCount
            }
            if (Object.keys(response.data).includes("videos")) {
              this.videos = response.data.videos
            }
          })
          .catch(function (error) {
            if (error.response) {
              console.log(error.response.data);
              console.log(error.response.status);
              console.log(error.response.headers);
            } else if (error.request) {
              console.log(error.request);
            } else {
              console.log('Error', error.message);
            }
            this.error = true
          });
    }
  },
  mounted() {
    console.log("axios")
    this.fetchVideos()
  }
}
</script>

<style scoped>

</style>