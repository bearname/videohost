<template>
  <v-container>
     <span class="group pa-2">
        <v-icon>home</v-icon>
        <v-icon>event</v-icon>
        <v-icon>info</v-icon>
      </span>
    <v-row>
      {{ listType }}
      <VideoList v-if="videos!== null" :show-status="true" :user-page="true" :videos="videos" :key="page"/>
    </v-row>
  </v-container>
</template>

<script>
import {mapActions, mapGetters} from "vuex";
import VideoList from "@/components/VideoList";
import {VAvatar} from "vuetify/lib";

export default {
  name: "Playlist",
  components: {
    VideoList,
    VAvatar,

  },
  async created() {
    this.listType = this.$route.query.PL;
    await this.fetchUserVideos(this.page, this.countVideo)
  },
  data() {
    return {
      videos: null,
      page: 0,
      countVideo: 30,
      listType: null,
      messages: [
        {
          avatar: 'https://avatars0.githubusercontent.com/u/9064066?v=4&s=460',
          name: 'John Leider',
          title: 'Welcome to Vuetify!',
          excerpt: 'Thank you for joining our community...',
        },
        {
          color: 'red',
          icon: 'mdi-account-multiple',
          name: 'Social',
          new: 1,
          total: 3,
          title: 'Twitter',
        },
        {
          color: 'teal',
          icon: 'mdi-tag',
          name: 'Promos',
          new: 2,
          total: 4,
          title: 'Shop your way',
          exceprt: 'New deals available, Join Today',
        },
      ],
      lorem: 'Lorem ipsum dolor sit amet, at aliquam vivendum vel, everti delicatissimi cu eos. Dico iuvaret debitis mel an, et cum zril menandri. Eum in consul legimus accusam. Ea dico abhorreant duo, quo illum minimum incorrupte no, nostro voluptaria sea eu. Suas eligendi ius at, at nemore equidem est. Sed in error hendrerit, in consul constituam cum.',
    }
  },
  methods: {
    ...mapActions({
      findUserLikedVideos: "userMod/getUserLikedVideos",
    }),
    ...mapGetters({
      getVideoResult: "userMod/getUserVideos",
    }),
    async fetchUserVideos(page, countVideo) {
      await this.findUserLikedVideos({page: page, countVideoOnPage: countVideo});
      const result = this.getVideoResult();
      this.videos = result.videos;
    }
  },
}
</script>

<style scoped>

</style>