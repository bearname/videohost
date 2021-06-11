<template>
  <div>
    <div>
      <label for="name">Name: <input id="name" type="text" v-model="name"></label>
    </div>
    <div>
      <label for="description">Description: <input id="description" type="text" v-model="description"></label>
    </div>
    <v-btn type="submit" v-on:click="updateTitleAndDescription">save</v-btn>
    <div>
      <p>{{ name }}</p>
      <p>{{ description }} </p>
    </div>
  </div>
</template>

<script>
import {VBtn} from "vuetify/lib";
import {mapActions, mapGetters} from "vuex";

export default {
  name: "VideoEdit",
  components: {
    VBtn
  },
  props: [
    "videoId",
    "videoName",
    "videoDescription"
  ],
  data() {
    return {
      id: this.videoId,
      name: this.videoName,
      description: this.videoDescription
    }
  },
  methods: {
    ...mapActions({
      updateVideoTitleAndDescription: "videoMod/updateTitleAndDescription"
    }),
    ...mapGetters({
      getVideoStatus: "videoMod/getStatus"
    }),
    async updateTitleAndDescription() {
      const video = {
        videoId: this.id,
        name: this.name,
        description: this.description,
      };
      await this.updateVideoTitleAndDescription(video);
      const status = this.getVideoStatus();
      console.log(status)
    }
  }
}
</script>

<style scoped>

</style>