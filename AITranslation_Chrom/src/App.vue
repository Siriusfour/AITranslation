<script  setup>
import { ref } from 'vue'



const value4 = ref('');
const visible = ref(true);
const status = ref("翻译页面")
function changeStatus() {
 status.value == "翻译页面" ? status.value = "显示原文" : status.value = "翻译页面"
}


// 点击按钮触发
async function greet() {
  changeStatus()

  let [tab] = await chrome.tabs.query({ active: true, currentWindow: true });

  // 向目标页面里注入js方法
  chrome.scripting.executeScript({
    target: { tabId: tab.id },
    files: ['/assets/AITranslation.js']
  });
}


</script>

<template>
  <div>
    <a href="https://vite.dev" target="_blank">
      <img src="/vite.svg" class="logo" alt="Vite logo" />
    </a>
  </div>
    <a-space>
      <a-input-password
        v-model:value="value4"
        v-model:visible="visible"
        placeholder="input password"
      />
      <a-button @click="visible = !visible">{{ visible ? 'Hide' : 'Show' }}</a-button>
    </a-space>

  <button @click="greet" style="margin-top: 30px;">{{ status }}</button>

</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
.Translation_btn {
  width:50% ;
  height: 50px;
  color:crimson;
}
</style>
