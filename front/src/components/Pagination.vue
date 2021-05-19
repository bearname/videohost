<template>
  <div
      class="mx-auto"
      max-width="1000"
  >
    <div v-if="videos !== null">
      <VideoList :videos="videos" :key="pageNumber"/>
      <v-btn class="v-btn"
             :disabled="pageNumber === 0"
             @click="previousPage">
        Previous
      </v-btn>
      <v-btn>{{ pageNumber + 1 }}</v-btn>
      <v-btn class="v-btn"
             :disabled="(pageNumber >= countPage && countPage < 2) || (pageNumber >= countPage - 1 && countPage >= 2) "
             @click="nextPage">
        Next
      </v-btn>
    </div>
    <p v-else>Not found any video</p>
  </div>

</template>

<script>
import axios from 'axios'
import VideoList from "./VideoList";

export default {
  name: "Pagination",
  components: {
    VideoList
  },

  data() {
    return {
      error: false,
      pageNumber: 1,
      videos: null,
      countVideoOnPage: 12,
      countPage: 0,
      url: {
        type: String,
        required: false,
        default: "list"
      }
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
              this.videos.forEach(function (part, index) {
                this[index].thumbnail = "http://localhost:8000/" + this[index].thumbnail
              }, this.videos)
              console.log("vidoes", this.videos)
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