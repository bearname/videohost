<template>
  <div>
    <span class="search-wrapper">{{ searchString }}
      <input type="text" v-model="searchString" placeholder="Search title.."/>
    </span>
    <router-link :to="{ name: 'search',  params:{ searchQuery: searchString }}">search</router-link>
  </div>
</template>

<script>
export default {
  name: "SearchRow",
  data() {
    return {
      searchString: ''
    }
  }
}
</script>

<style scoped>

</style>