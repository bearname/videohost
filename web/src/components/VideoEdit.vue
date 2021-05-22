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
export default {
  name: "VideoEdit",
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
    updateTitleAndDescription() {
      const video = {
        videoId: this.id,
        name: this.name,
        description: this.description,
      };
      this.$store.dispatch("video/updateTitleAndDescription", video)
      const getter = this.$store.getters["video/getStatus"];
      console.log(getter)
    }
  }
}
</script>

<style scoped>

</style>