<template>
  <div>
    <form enctype="multipart/form-data"
          action="http://localhost:8000/api/v1/video"
          method="post">
      <label for="name">Name <input id="name" type="text" name="name"/></label>
      <label for="description">Description <input id="description" type="text" name="description"/></label>
      <label for="file">Video file <input id="file" type="file" name="file"/></label>
      <button type="submit" value="upload">Upload</button>
    </form>
  </div>
</template>

<script>
export default {
  name: "UploadVideoPage"
}
</script>

<style scoped>

</style>