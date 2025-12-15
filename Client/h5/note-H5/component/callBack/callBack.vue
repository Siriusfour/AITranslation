<script setup>
import { ref, onMounted } from "vue";
import api from "../../src/Utils/Api/Api.js";

const code = ref("");
const state = ref("");
const userName = ref("")
const avatar = ref("")


onMounted(async () => {

  const params = new URLSearchParams(window.location.search);
  const data  ={
    code:params.get("code") || "",
    state: params.get("state") || ""
  }

  api.post("/NotAuth/LoginByGithub",data).then(res => {
    console.log(res);
    userName.value=res.nickname;
    avatar.value=res.avatar;
  })

});
</script>

<template>
  <div class="container">
    <p>code: {{ code }}</p>
    <p>state: {{ state }}</p>

    <img :src="avatar" />
    <p>userName: {{ userName }}</p>
  </div>
</template>

<style scoped>
.container {
  padding: 16px;
}
</style>
