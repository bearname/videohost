<template>
  <div>
    <h2>Home</h2>
    <Pagination :show-status="false" :user-page="false"/>
  </div>
</template>

<script>
import Pagination from '../components/Pagination.vue'

export default {
  name: "Home",
  components: {
    Pagination
  }
}
</script>

<style scoped>

</style>