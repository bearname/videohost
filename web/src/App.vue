<template>
  <layout>
    <div>
      <Navigation/>
      <h1>Pretty player</h1>
      <Search/>
      <Pagination :show-status="false" :user-page="false"/>
    </div>
  </layout>
</template>

<script>
import Layout from './layouts/LayoutDefault.vue';
import Navigation from './components/Navigation.vue';
import Search from "./components/Search";
import Pagination from './components/Pagination.vue';

export default {
  name: 'App',
  components: {
    Layout,
    Navigation,
    Search,
    Pagination
  },
};
</script>

<style>
a {
  text-decoration: none;
}

</style>