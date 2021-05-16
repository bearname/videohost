<template>
  <ul>
    <li>
      <router-link :to="{ name: 'home'}">Каталог</router-link>
    </li>
    <li>
      <router-link :to="{ name: 'uploadVideo'}">Загрузить видео</router-link>
    </li>
  </ul>
</template>

<script>
export default {
  name: "Navigation"
}
</script>

<style scoped>

</style>