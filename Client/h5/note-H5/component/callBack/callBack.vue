<script setup>
import { ref, onMounted } from "vue";
import api from "../../src/Utils/Api/Api.js";

const code = ref("");
const state = ref("");

onMounted(async () => {

  const params = new URLSearchParams(window.location.search);
  code.value = params.get("code") || "";
  state.value = params.get("state") || "";

  api.post("/Login").then(data => {})

});
</script>

<template>
  <div class="container">
    <p>code: {{ code }}</p>
    <p>state: {{ state }}</p>
  </div>
</template>

<style scoped>
.container {
  padding: 16px;
}
</style>
