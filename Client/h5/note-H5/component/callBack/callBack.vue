<script setup>
import { ref, onMounted } from "vue";
import api from "../../src/Utils/Api/Api.js";

const code = ref("");
const state = ref("");

onMounted(async () => {

  const params = new URLSearchParams(window.location.search);
  const data  ={
    code:params.get("code") || "",
    state: params.get("state") || ""
  }

  api.post("/Auth/LoginByGithub",data).then(res => {




  })

});
</script>

<template>
  <div class="container">
    <p>code: {{ code }}</p>
    <p>state: {{ state }}</p>

    <p>avatar: {{ avatar }}</p>
    <p>userName: {{ userName }}</p>
  </div>
</template>

<style scoped>
.container {
  padding: 16px;
}
</style>
