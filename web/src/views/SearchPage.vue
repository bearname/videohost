<template>
  <div>
    {{ search }}
    <!--    <Search :search-query="search" :key="search"/>-->
  </div>
</template>

<script>
// import Search from "../components/Search";

export default {
  name: "SearchPage",
  props: [
    'searchQuery'
  ],
  data() {
    return {
      search: this.$route.params.searchQuery
    }
  },
  components: {
    // Search,
  }
}
</script>

<style scoped>

</style>